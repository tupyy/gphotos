package v1

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/filter"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"go.uber.org/zap"
)

// (GET /api/gphotos/v1/albums)
func (server *Server) GetAlbums(c *gin.Context, params apiv1.GetAlbumsParams) {
	session := c.MustGet("session").(entity.Session)

	albumService := server.AlbumService()

	q := albumService.Query().OwnAlbums(true)
	if params.Personal != nil {
		q.OwnAlbums(*params.Personal)
	}
	if params.Shared != nil {
		q.SharedAlbums(*params.Shared)
	}
	if params.Search != nil {
		searchExp := *params.Search
		if se, err := strconv.Unquote(*params.Search); err == nil {
			searchExp = se
		}
		filter, err := filter.New(searchExp)
		if err != nil {
			zap.S().Errorw("failed to create filter engine", "filter", *&params.Search, "user", session.User.Username)
			c.AbortWithStatusJSON(http.StatusBadRequest, mappersv1.MapFromStatus(http.StatusBadRequest, "malformatted search expression"))
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

	albums, total, err := q.All(c, session.User)
	if err != nil {
		zap.S().Errorw("failed to get albums", "error", err, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	albumModels := make([]apiv1.Album, 0, len(albums))
	for _, album := range albums {
		albumModels = append(albumModels, mappersv1.MapAlbumToModel(album))
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
func (server *Server) GetAlbumsByGroup(c *gin.Context, groupId string, params apiv1.GetAlbumsByGroupParams) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(groupId)
	if err != nil {
		zap.S().Errorw("failed to decrypt group id", "error", err, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "group with id '%s' not found", groupId))
		return
	}

	albumService := server.AlbumService()

	filter, _ := filter.New(fmt.Sprintf("permissions.group = '%s'", id))
	q := albumService.Query().
		OwnAlbums(false).
		SharedAlbums(true).
		Filter(filter)

	// paginate
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
		q.Page(page)
	}

	if params.Size != nil {
		q.Size(int(*params.Size))
	}

	albums, total, err := q.All(c, session.User)
	if err != nil {
		zap.S().Errorw("failed to get albums", "error", err, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	albumModels := make([]apiv1.Album, 0, len(albums))
	for _, album := range albums {
		albumModels = append(albumModels, mappersv1.MapAlbumToModel(album))
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

// (GET /api/gphotos/v1/albums/users/{user_id})
func (server *Server) GetAlbumsByUser(c *gin.Context, userId string, params apiv1.GetAlbumsByUserParams) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(userId)
	if err != nil {
		zap.S().Errorw("failed to decrypt user id", "error", err, "user id", userId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "user with id '%s' not found", userId))
		return
	}

	albumService := server.AlbumService()

	filter, _ := filter.New(fmt.Sprintf("permissions.user = '%s'", id))
	q := albumService.Query().
		OwnAlbums(false).
		SharedAlbums(true).
		Filter(filter)

	// paginate
	page := 1
	if params.Page != nil {
		page = int(*params.Page)
		q.Page(page)
	}

	if params.Size != nil {
		q.Size(int(*params.Size))
	}

	albums, total, err := q.All(c, session.User)
	if err != nil {
		zap.S().Errorw("failed to get albums", "error", err, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	albumModels := make([]apiv1.Album, 0, len(albums))
	for _, album := range albums {
		albumModels = append(albumModels, mappersv1.MapAlbumToModel(album))
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

func (server *Server) GetAlbumByID(c *gin.Context, albumID apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumID)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumID))
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album id", id, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	// only users with editPermission set for this album or one of user's group with the same permission
	// can edit this album
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
		Policy(permissions.AnyUserPermissionPolicty{}).
		Policy(permissions.AnyGroupPermissionPolicy{}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		zap.S().Errorw("permission denied to access album", "album", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "access denied"))
		return
	}

	zap.S().Infow("album access ok", "album id", albumID, "user", session.User.Username)
	c.JSON(http.StatusOK, mappersv1.MapAlbumToModel(album))
}

func (server *Server) CreateAlbum(c *gin.Context) {
	session := c.MustGet("session").(entity.Session)

	// only editors and admins have the right to create albums
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.RolePolicy{Role: entity.RoleEditor}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(entity.Album{}, session.User)

	if !hasPermission {
		zap.S().Errorw("permission denied to create album", "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "access denied"))
		return
	}

	var payload apiv1.AlbumRequestPayload
	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, mappersv1.MapFromStatusf(http.StatusBadRequest, "failed to parse payload: %s", err))
		return
	}

	if len(payload.Name) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, mappersv1.MapFromStatus(http.StatusBadRequest, "album's name is missing"))
		return
	}

	zap.S().Debugw("album request payload", "payload", payload, "user", session.User.Username)
	album := entity.Album{
		Name:        html.EscapeString(payload.Name),
		Description: escapeFieldPtr(payload.Description),
		CreatedAt:   time.Now(),
		Location:    escapeFieldPtr(payload.Location),
		Owner:       session.User.Username,
	}
	albumID, err := server.AlbumService().Create(c, album)
	if err != nil {
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	zap.S().Infow("album created", "album id", albumID, "user", session.User.Username)

	ss := sessions.Default(c)
	ss.Set(session.SessionID, session)
	ss.Save()

	c.JSON(http.StatusOK, mappersv1.MapAlbumToModel(album))
}

func (server *Server) UpdateAlbum(c *gin.Context, albumID apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumID)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumID))
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumID))
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
		zap.S().Errorw("failed to update album. user has no edit permission", "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "access denied"))
		return
	}

	var payload apiv1.AlbumRequestPayload
	if err := c.BindJSON(&payload); err != nil {
		zap.S().Errorw("failed to bind payload", "album id", id, "error", err, "payload", payload, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusBadRequest, mappersv1.MapFromStatusf(http.StatusBadRequest, "failed to parse payload: %s", err))
		return
	}

	// update album
	if payload.Description != nil {
		album.Description = escapeFieldPtr(payload.Description)
	}

	if payload.Location != nil {
		album.Location = escapeFieldPtr(payload.Location)
	}

	if payload.Name != "" {
		album.Name = escapeField(payload.Name)
	}

	if _, err := server.AlbumService().Update(c, album); err != nil {
		zap.S().Errorw("failed to update album", "album id", id, "error", err, "payload", payload, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	zap.S().Infow("album updated", "album id", id, "user", session.User.Username)

	c.JSON(http.StatusCreated, mappersv1.MapAlbumToModel(album))
}

// (DELETE /api/gphotos/v1/albums/{album_id})
func (server *Server) DeleteAlbum(c *gin.Context, albumId apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumId, "album", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumId))
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumId))
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
		zap.S().Errorw("failed to delete album. missing permissions", "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "access denied"))
		return
	}

	if err := server.AlbumService().Delete(c, album); err != nil {
		zap.S().Errorw("failed to delete album", "error", err, "album id", id, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func escapeFieldPtr(fieldValue *string) string {
	if fieldValue != nil {
		return html.EscapeString(*fieldValue)
	}
	return ""
}

func escapeFieldPtr2(fieldValue *string) *string {
	if fieldValue != nil {
		value := html.EscapeString(*fieldValue)
		return &value
	}
	return nil
}

func escapeField(fieldValue string) string {
	return html.EscapeString(fieldValue)
}
