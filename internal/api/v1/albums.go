package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/filter"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/utils/encryption"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

// (GET /api/gphotos/v1/albums)
func (server *Server) GetAlbums(c *gin.Context, params apiv1.GetAlbumsParams) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(c)

	albumService := server.GetAlbumService()

	q := albumService.Query()
	if params.Personal != nil {
		q.OwnAlbums(*params.Personal)
	}
	if params.Shared != nil {
		q.SharedAlbums(*params.Shared)
	}
	if params.Filter != nil {
		filter, err := filter.New(*params.Filter)
		if err != nil {
			logger.WithError(err).Error("failed to create filter engine")
			common.AbortBadRequestWithJson(c, err, err.Error())

			return
		}

		q.Filter(filter)
	}
	// setup sort
	if params.Sort != nil {
		switch *params.Sort {
		case "name":
			q.Sort(album.SortByName, album.NormalOrder)
		case "location":
			q.Sort(album.SortByLocation, album.NormalOrder)
		case "date-normal":
			q.Sort(album.SortByDate, album.NormalOrder)
		default:
			q.Sort(album.SortByDate, album.ReverseOrder)
		}
	}

	// paginate
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
		q.Page(page)
	}

	if params.Size != nil {
		q.Size(int(*params.Size))
	}

	albums, total, err := q.All(ctx, session.User)
	if err != nil {
		logger.WithError(err).Error("failed to get albums")
		common.AbortInternalErrorWithJson(c)
		return
	}

	albumModels := make([]apiv1.Album, 0, len(albums))
	for _, album := range albums {
		albumModels = append(albumModels, mapAlbumToModel(album))
	}

	c.JSON(http.StatusOK, &apiv1.AlbumList{
		Kind:  "AlbumList",
		Page:  page,
		Size:  len(albumModels),
		Total: total,
		Items: albumModels,
	})

	return
}

// (GET /api/gphotos/v1/albums/groups/{group_id})
func (server *Server) GetAlbumsByGroup(c *gin.Context, groupId string) {
}

// (GET /api/gphotos/v1/albums/users/{user_id})
func (server *Server) GetAlbumsByUser(c *gin.Context, userId string) {
}

// (GET /api/gphotos/v1/albums/{album_id}/permissions)
func (server *Server) GetAlbumPermissions(c *gin.Context, albumId string) {
}

func (server *Server) GetAlbumByID(c *gin.Context, albumID apiv1.AlbumId) {}

func (server *Server) CreateAlbum(c *gin.Context) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	// only editors and admins have the right to create albums
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.RolePolicy{Role: entity.RoleEditor}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(entity.Album{}, session.User)

	if !hasPermission {
		common.AbortForbidden(c, errors.New("user has no editor or admin role"), "user role forbids the creation of albums")
		return
	}

	var payload apiv1.AlbumRequestPayload
	if err := c.BindJSON(&payload); err != nil {
		common.AbortBadRequest(c, err, "fail to bind to form")
		return
	}

	logger.WithField("form", fmt.Sprintf("%+v", payload)).Info("create album request submitted")
	album := entity.Album{
		Name:        html.EscapeString(payload.Name),
		Description: escapeFieldPtr(payload.Description),
		CreatedAt:   time.Now(),
		Location:    escapeFieldPtr(payload.Location),
		Owner:       session.User.Username,
	}
	if payload.UserPermissions != nil {
		permForm := make(map[string][]string)

		err := json.Unmarshal(bytes.NewBufferString(*payload.UserPermissions).Bytes(), &permForm)
		if err != nil {
			logger.WithField("permissions_string", *payload.UserPermissions).WithError(err).Warn("unmarshal error")
		} else {
			var pp = make(entity.Permissions)
			pp.Parse(permForm, true)
			album.UserPermissions = pp
		}
	}

	if payload.GroupPermissions != nil {
		permForm := make(map[string][]string)

		err := json.Unmarshal(bytes.NewBufferString(*payload.GroupPermissions).Bytes(), &permForm)
		if err != nil {
			logger.WithField("permissions_string", *&payload.GroupPermissions).WithError(err).Warn("unmarshal error")
		} else {
			var pp = make(entity.Permissions)
			pp.Parse(permForm, false)
			album.GroupPermissions = pp
		}
	}

	albumID, err := server.GetAlbumService().Create(ctx, album)
	if err != nil {
		common.AbortInternalError(c)

		return
	}

	logger.WithFields(logrus.Fields{
		"album": fmt.Sprintf("%+v", album),
		"id":    albumID,
	}).Info("album entity created")

	ss := sessions.Default(c)
	ss.Set(session.SessionID, session)
	ss.Save()

	c.JSON(http.StatusOK, mapAlbumToModel(album))
}

func (server *Server) UpdateAlbum(c *gin.Context, albumID apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := descryptAlbumID(albumID)
	if err != nil {
		logger.WithError(err).WithField("album id", albumID).Error("failed to decrypt album id")
		common.AbortInternalError(c)
		return
	}

	album, err := server.GetAlbumService().Query().First(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
		common.AbortNotFound(c, err, "update album")

		return
	}

	// only users with editPermission set for this album or one of user's group with the same permission
	// can edit this album
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
		Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionEditAlbum}).
		Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionEditAlbum}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		logger.WithFields(logrus.Fields{
			"request user id": session.User.ID,
			"album owner id":  album.Owner,
		}).Error("album can be edit either by user with edit permission or the owner")
		common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "update album")
		return
	}

	var payload apiv1.AlbumRequestPayload
	if err := c.BindJSON(&payload); err != nil {
		logger.WithError(err).WithField("payload", fmt.Sprintf("%v", payload)).Error("failed to bind payload")
		common.AbortBadRequest(c, err, "failed to parse payload")
		return
	}

	// update album
	album.Description = escapeFieldPtr(payload.Description)
	album.Location = escapeFieldPtr(payload.Location)
	if payload.Name != "" {
		album.Name = escapeField(payload.Name)
	}

	// only the owner and an admin has the right to edit permissions
	if session.User.Username == album.Owner || session.User.Role == entity.RoleAdmin {
		// add new permissions if any
		if payload.UserPermissions != nil {
			permForm := make(map[string][]string)
			err := json.Unmarshal(bytes.NewBufferString(*payload.UserPermissions).Bytes(), &permForm)
			if err != nil {
				logger.WithField("permissions_string", *payload.UserPermissions).WithError(err).Warn("unmarshal error")
			} else {
				var pp = make(entity.Permissions)
				pp.Parse(permForm, true)
				album.UserPermissions = pp
			}
		}

		if payload.GroupPermissions != nil {
			permForm := make(map[string][]string)
			err := json.Unmarshal(bytes.NewBufferString(*payload.GroupPermissions).Bytes(), &permForm)
			if err != nil {
				logger.WithField("permissions_string", *payload.UserPermissions).WithError(err).Warn("unmarshal error")
			} else {
				var pp = make(entity.Permissions)
				pp.Parse(permForm, false)
				album.GroupPermissions = pp
			}
		}
	}

	if _, err := server.GetAlbumService().Update(ctx, album); err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"album id": id,
			"album":    fmt.Sprintf("%+v", album),
		}).Error("update album")

		common.AbortInternalError(c)

		return
	}

	c.JSON(http.StatusOK, mapAlbumToModel(album))
}

// (DELETE /api/gphotos/v1/albums/{album_id})
func (server *Server) DeleteAlbum(c *gin.Context, albumId apiv1.AlbumId) {
	s, _ := c.Get("sessionData")
	session := s.(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := descryptAlbumID(albumId)
	if err != nil {
		logger.WithError(err).WithField("album id", albumId).Error("failed to decrypt album id")
		common.AbortInternalError(c)
		return
	}

	album, err := server.GetAlbumService().Query().First(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
		common.AbortNotFound(c, err, "update album")

		return
	}

	// only users with editPermission set for this album or one of user's group with the same permission
	// can edit this album
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionDeleteAlbum}).
		Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionDeleteAlbum}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		logger.WithFields(logrus.Fields{
			"request user id": session.User.ID,
			"album owner id":  album.Owner,
		}).Error("album can be edit either by user with delete permission or the owner")
		common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionDeleteAlbum, album, session.User), "delete album")

		return
	}

	if err := server.GetAlbumService().Delete(ctx, album); err != nil {
		logger.WithError(err).WithField("allbum id", album.ID).Error("failed to delete album")
		common.AbortInternalError(c)

		return
	}

	ss := sessions.Default(c)
	ss.Set(session.SessionID, session)
	ss.Save()

	c.JSON(http.StatusOK, gin.H{})
}

func descryptAlbumID(albumID string) (int32, error) {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	id, err := gen.DecryptData(albumID)
	if err != nil {
		return 0, err
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return int32(intID), nil
}

func escapeFieldPtr(fieldValue *string) string {
	if fieldValue != nil {
		return html.EscapeString(*fieldValue)
	}
	return ""
}

func escapeField(fieldValue string) string {
	return html.EscapeString(fieldValue)
}
