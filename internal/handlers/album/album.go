package album

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/dto"
	"github.com/tupyy/gophoto/internal/form"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
)

const (
	rootURL = "/"
)

// TODO fix the error management. it totally crap.
// GET /album/:id
func GetAlbum(r *gin.RouterGroup, albumService *album.Service, usersService *users.Service) {

	r.GET("/album/:id", parseAlbumIDHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
		if err != nil {
			logger.WithError(err).WithField("id", c.GetInt("id")).Error("album not found")
			common.AbortNotFound(c, err, "failed to album")

			return
		}

		// check permissions to this album
		ats := permissions.NewAlbumPermissionService()
		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
			Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
			Policy(permissions.AnyUserPermissionPolicty{}).
			Policy(permissions.AnyGroupPermissionPolicy{}).
			Strategy(permissions.AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("user has no permissions to access this album")

			common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionReadAlbum, album, session.User), "")

			return
		}

		users, err := usersService.Query().
			Where(users.NotUsername(session.User.Username)).
			Where(users.CanShare(true)).
			Where(users.Roles([]entity.Role{entity.RoleEditor, entity.RoleUser})).
			AllUsers(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get users")
			common.AbortInternalError(c)

			return
		}

		// if not owner get the owner from keycloak
		owner, err := usersService.Query().FirstUser(ctx, album.OwnerID)
		if err != nil {
			logger.WithError(err).WithField("album id", album.ID).Error("failed to fetch owner from keycloak")
			common.AbortInternalError(c)

			return
		}

		// replace ids with names in user permissions maps and OwnerID with owner's name
		userPermissions := make(map[string][]entity.Permission)
		for _, u := range users {
			if perms, found := album.UserPermissions[u.ID]; found {
				name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)

				if name == "" {
					logger.WithField("username", u.Username).Warn("user has not first or last name set")

					continue
				}

				userPermissions[name] = perms
			}
		}

		// check individual permissions for this album
		permissions := make(map[entity.Permission]bool)
		permissions[entity.PermissionReadAlbum] = entity.HasUserPermission(album, session.User.ID, entity.PermissionReadAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin
		permissions[entity.PermissionWriteAlbum] = entity.HasUserPermission(album, session.User.ID, entity.PermissionWriteAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin
		permissions[entity.PermissionEditAlbum] = entity.HasUserPermission(album, session.User.ID, entity.PermissionEditAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin
		permissions[entity.PermissionDeleteAlbum] = entity.HasUserPermission(album, session.User.ID, entity.PermissionDeleteAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin

		for _, g := range session.User.Groups {
			if perms, found := album.GroupPermissions[g.Name]; found {
				for _, p := range perms {
					permissions[p] = true
				}
			}
		}

		albumDTO, err := dto.NewAlbumDTO(album, owner)
		if err != nil {
			logger.WithError(err).WithField("album", album.String()).Error("failed to serialize album")

			common.AbortInternalError(c)

			return
		}

		c.HTML(http.StatusOK, "album_view.html", gin.H{
			"album":             albumDTO,
			"is_owner":          session.User.ID == album.OwnerID,
			"owner":             fmt.Sprintf("%s %s", owner.FirstName, owner.LastName),
			"user_permissions":  userPermissions,
			"group_permissions": album.GroupPermissions,
			"delete_link":       fmt.Sprintf("/album/%s/delete", albumDTO.ID),
			"edit_link":         fmt.Sprintf("/album/%s/edit", albumDTO.ID),
			"read_permission":   permissions[entity.PermissionReadAlbum],
			"write_permission":  permissions[entity.PermissionWriteAlbum],
			"edit_permission":   permissions[entity.PermissionEditAlbum],
			"delete_permission": permissions[entity.PermissionDeleteAlbum],
			"is_admin":          session.User.Role == entity.RoleAdmin,
		})
	})
}

// GET /album
func GetCreateAlbumForm(r *gin.RouterGroup, usersService *users.Service) {
	r.GET("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			logger.Error("user with user role cannot create albums")
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))

			return
		}

		users, err := usersService.Query().
			Where(users.NotUsername(session.User.Username)).
			Where(users.CanShare(true)).
			Where(users.Roles([]entity.Role{entity.RoleEditor, entity.RoleUser})).
			AllUsers(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get users")
			common.AbortInternalError(c)

			return
		}

		groups, err := usersService.Query().AllGroups(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get groups")
			common.AbortInternalError(c)

			return
		}

		usersDTO := dto.NewUserDTOs(users)

		c.HTML(http.StatusOK, "album_form.html", gin.H{
			"users":          usersDTO,
			"groups":         groups,
			"canShare":       session.User.CanShare,
			"isOwner":        true,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	})
}

// POST /album
func CreateAlbum(r *gin.RouterGroup, albumService *album.Service) {
	r.POST("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

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

		var albumForm form.Album
		if err := c.ShouldBind(&albumForm); err != nil {
			common.AbortBadRequest(c, err, "fail to bind to form")

			return
		}

		cleanForm := albumForm.Sanitize()

		logger.WithField("form", fmt.Sprintf("%+v", cleanForm)).Info("create album request submitted")

		album := entity.Album{
			Name:        cleanForm.Name,
			Description: cleanForm.Description,
			CreatedAt:   time.Now(),
			Location:    cleanForm.Location,
			OwnerID:     session.User.ID,
		}

		if len(cleanForm.UserPermissions) > 0 {
			permForm := make(map[string][]string)

			err := json.Unmarshal(bytes.NewBufferString(cleanForm.UserPermissions).Bytes(), &permForm)
			if err != nil {
				logger.WithField("permissions_string", cleanForm.UserPermissions).WithError(err).Warn("unmarshal error")
			} else {
				var pp = make(entity.Permissions)
				pp.Parse(permForm, true)
				album.UserPermissions = pp
			}
		}

		if len(cleanForm.GroupPermissions) > 0 {
			permForm := make(map[string][]string)

			err := json.Unmarshal(bytes.NewBufferString(cleanForm.GroupPermissions).Bytes(), &permForm)
			if err != nil {
				logger.WithField("permissions_string", cleanForm.UserPermissions).WithError(err).Warn("unmarshal error")
			} else {
				var pp = make(entity.Permissions)
				pp.Parse(permForm, false)
				album.GroupPermissions = pp
			}
		}

		albumID, err := albumService.Create(ctx, album)
		if err != nil {
			common.AbortInternalError(c)

			return
		}

		logger.WithFields(logrus.Fields{
			"album": fmt.Sprintf("%+v", album),
			"id":    albumID,
		}).Info("album entity created")

		alert := entity.Alert{
			Message: fmt.Sprintf("Album %s created.", album.Name),
			IsError: false,
		}
		session.AddAlert(alert)
		session.AddAlert(entity.Alert{
			Message: "test",
			IsError: true,
		})

		ss := sessions.Default(c)
		ss.Set(session.SessionID, session)
		ss.Save()

		c.Redirect(http.StatusFound, rootURL)
	})
}

// GET /album/:id/edit
func GetUpdateAlbumForm(r *gin.RouterGroup, albumService *album.Service, usersService *users.Service) {
	r.GET("/album/:id/edit", parseAlbumIDHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(c)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
		if err != nil {
			common.AbortNotFound(c, err, "update album")

			return
		}

		albumDTO, err := dto.NewAlbumDTO(album, session.User)
		if err != nil {
			logger.WithError(err).WithField("id", album.ID).Error("failed to serialize the album")
			common.AbortInternalError(c)

			return

		}

		// check if user is the owner or it has the edit permission set
		if album.OwnerID == session.User.ID || session.User.Role == entity.RoleAdmin {
			logger.Info("edit permission granted. user is the owner")

			users, err := usersService.Query().
				Where(users.NotUsername(session.User.Username)).
				Where(users.CanShare(true)).
				Where(users.Roles([]entity.Role{entity.RoleEditor, entity.RoleUser})).
				AllUsers(ctx)
			if err != nil {
				logger.WithError(err).Error("failed to get users")
				common.AbortInternalError(c)

				return
			}

			groups, err := usersService.Query().AllGroups(ctx)
			if err != nil {
				logger.WithError(err).Error("failed to get groups")
				common.AbortInternalError(c)

				return
			}

			c.HTML(http.StatusOK, "album_form.html", gin.H{
				"update_link":    fmt.Sprintf("/album/%s", albumDTO.ID),
				"album":          albumDTO,
				"canShare":       session.User.CanShare,
				"isOwner":        true,
				"users":          dto.NewUserDTOs(users),
				"groups":         groups,
				"is_admin":       session.User.Role == entity.RoleAdmin,
				csrf.TemplateTag: csrf.TemplateField(c.Request),
			})

			return
		}

		// only users with editPermission set for this album or one of user's group with the same permission
		// can edit this album
		apr := permissions.NewAlbumPermissionService()
		hasPermission := apr.Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionEditAlbum}).
			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionEditAlbum}).
			Strategy(permissions.AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"request user id": session.User.ID,
				"album owner id":  album.OwnerID,
			}).Error("album cannot be edit either by user with edit permission or the owner")
			common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "update album")

			return
		}

		c.HTML(http.StatusOK, "album_form.html", gin.H{
			"album":          albumDTO,
			"canShare":       session.User.CanShare,
			"isOwner":        false,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	})

}

// PUT /album/:id/
func UpdateAlbum(r *gin.RouterGroup, albumService *album.Service) {
	r.POST("/album/:id/", parseAlbumIDHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
		if err != nil {
			logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
			common.AbortNotFound(c, err, "update album")

			return
		}

		var albumForm form.Album
		if err := c.ShouldBind(&albumForm); err != nil {
			logger.WithError(err).WithField("query parameters", fmt.Sprintf("%v", albumForm)).Error("failed to bind query parameters to form")
			common.AbortBadRequest(c, err, "fail to bind to form")

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
				"album owner id":  album.OwnerID,
			}).Error("album can be edit either by user with edit permission or the owner")
			common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "update album")

			return
		}

		// update album

		cleanForm := albumForm.Sanitize()
		album.Description = cleanForm.Description
		album.Location = cleanForm.Location

		if cleanForm.Name == "" {
			common.AbortBadRequest(c, errors.New("name is missing"), "update album")

			return
		}

		album.Name = cleanForm.Name

		if session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin {
			// add new permissions if any
			if len(cleanForm.UserPermissions) > 0 {
				permForm := make(map[string][]string)

				err := json.Unmarshal(bytes.NewBufferString(cleanForm.UserPermissions).Bytes(), &permForm)
				if err != nil {
					logger.WithField("permissions_string", cleanForm.UserPermissions).WithError(err).Warn("unmarshal error")
				} else {
					var pp = make(entity.Permissions)
					pp.Parse(permForm, true)
					album.UserPermissions = pp
				}
			}

			if len(cleanForm.GroupPermissions) > 0 {
				permForm := make(map[string][]string)

				err := json.Unmarshal(bytes.NewBufferString(cleanForm.GroupPermissions).Bytes(), &permForm)
				if err != nil {
					logger.WithField("permissions_string", cleanForm.UserPermissions).WithError(err).Warn("unmarshal error")
				} else {
					var pp = make(entity.Permissions)
					pp.Parse(permForm, false)
					album.GroupPermissions = pp
				}
			}
		}

		c.Redirect(http.StatusFound, rootURL)
	})
}

// DELETE /album/:id
func DeleteAlbum(r *gin.RouterGroup, albumService *album.Service) {
	r.GET("/album/:id/delete", parseAlbumIDHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
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
				"album owner id":  album.OwnerID,
			}).Error("album can be edit either by user with delete permission or the owner")
			common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionDeleteAlbum, album, session.User), "delete album")

			return
		}

		if err := albumService.Delete(ctx, album); err != nil {
			logger.WithError(err).WithField("allbum id", album.ID).Error("failed to delete album")
			common.AbortInternalError(c)

			return
		}

		alert := entity.Alert{
			Message: fmt.Sprintf("Album %s deleted.", album.Name),
			IsError: false,
		}
		session.AddAlert(alert)

		ss := sessions.Default(c)
		ss.Set(session.SessionID, session)
		ss.Save()

		c.Redirect(http.StatusFound, rootURL)
	})
}
