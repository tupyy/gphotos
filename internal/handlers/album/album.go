package album

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/utils"
	"github.com/tupyy/gophoto/internal/form"
	"github.com/tupyy/gophoto/internal/handlers/common"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

const (
	rootURL = "/"
)

// GET /album/:id
func GetAlbum(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	keycloakRepo := repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	r.GET("/album/:id", func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		// decrypt album id
		gen := encryption.NewGenerator(conf.GetEncryptionKey())

		decryptedID, err := gen.DecryptData(c.Param("id"))
		if err != nil {
			logger.WithError(err).Error("cannot decrypt album id")
			c.AbortWithError(http.StatusNotFound, err) // explicit return not found here

			return
		}

		id, err := strconv.Atoi(decryptedID)
		if err != nil {
			logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		album, err := albumRepo.GetByID(reqCtx, int32(id))
		if err != nil {
			common.AbortNotFound(c, err, "update album")

			return
		}

		// check permissions to this album
		atr := NewAlbumPermissionResolver()
		hasPermission := atr.Policy(OwnerPolicy{}).
			Policy(RolePolicy{entity.RoleAdmin}).
			Policy(AnyUserPermissionPolicty{}).
			Policy(AnyGroupPermissionPolicy{}).
			Strategy(AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("user has no permissions to access this album")

			common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionReadAlbum, album, session.User), "")

			return
		}

		users, err := keycloakRepo.GetUsers(reqCtx, nil)
		if err != nil {
			logger.WithError(err).Error("fetch users")
			common.AbortInternalError(c, errors.New("fetch users from keyclosk"), "")

			return
		}

		// if not owner get the owner from keycloak
		owner, err := keycloakRepo.GetUserByID(reqCtx, album.OwnerID)
		if err != nil {
			logger.WithError(err).WithField("id", id).Error("fetch owner")
			common.AbortInternalError(c, errors.New("fetch owner from keyclosk"), "")

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
		permissions[entity.PermissionReadAlbum] = utils.HasUserPermission(album, session.User.ID, entity.PermissionReadAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin
		permissions[entity.PermissionWriteAlbum] = utils.HasUserPermission(album, session.User.ID, entity.PermissionWriteAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin
		permissions[entity.PermissionEditAlbum] = utils.HasUserPermission(album, session.User.ID, entity.PermissionEditAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin
		permissions[entity.PermissionDeleteAlbum] = utils.HasUserPermission(album, session.User.ID, entity.PermissionDeleteAlbum) || session.User.ID == album.OwnerID || session.User.Role == entity.RoleAdmin

		for _, g := range session.User.Groups {
			if perms, found := album.GroupPermissions[g.Name]; found {
				for _, p := range perms {
					permissions[p] = true
				}
			}
		}

		// encrypt album id
		encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", album.ID))
		if err != nil {
			logger.WithError(err).Error("encrypt album id")
			common.AbortInternalError(c, err, fmt.Sprintf("album id: %d", album.ID))

			return
		}

		c.HTML(http.StatusOK, "album_view.html", gin.H{
			"name":              album.Name,
			"description":       album.Description,
			"location":          album.Location,
			"created_at":        album.CreatedAt,
			"is_owner":          session.User.ID == album.OwnerID,
			"owner":             fmt.Sprintf("%s %s", owner.FirstName, owner.LastName),
			"user_permissions":  userPermissions,
			"group_permissions": album.GroupPermissions,
			"delete_link":       fmt.Sprintf("/album/%s", encryptedID),
			"edit_link":         fmt.Sprintf("/album/%s/edit", encryptedID),
			"read_permission":   permissions[entity.PermissionReadAlbum],
			"write_permission":  permissions[entity.PermissionWriteAlbum],
			"edit_permission":   permissions[entity.PermissionEditAlbum],
			"delete_permission": permissions[entity.PermissionDeleteAlbum],
			"is_admin":          session.User.Role == entity.RoleAdmin,
		})
	})
}

// GET /album
func GetCreateAlbumForm(r *gin.RouterGroup, repos domain.Repositories) {
	keycloakRepo := repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	r.GET("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		reqCtx := c.Request.Context()

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))

			return
		}

		// filter out current user, admins and with can_share false
		userFilters, err := generateFilters(session.User)
		if err != nil {
			logutil.GetLogger(c).WithError(err).Error("create user fitlers")
			common.AbortInternalError(c, err, "")

			return
		}

		users, err := keycloakRepo.GetUsers(reqCtx, userFilters...)
		if err != nil && errors.Is(err, domain.ErrInternalError) {
			common.AbortInternalError(c, err, "cannot fetch users")

			return
		}

		groups, err := keycloakRepo.GetGroups(reqCtx)
		if err != nil && errors.Is(err, domain.ErrInternalError) {
			common.AbortInternalError(c, err, "cannot fetch groups")

			return
		}

		c.HTML(http.StatusOK, "album_form.html", gin.H{
			"users":          serialize(users),
			"groups":         groups,
			"canShare":       session.User.CanShare,
			"isOwner":        true,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	})
}

// POST /album
func CreateAlbum(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)

	r.POST("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		// only editors and admins have the right to create albums
		apr := NewAlbumPermissionResolver()
		hasPermission := apr.Policy(RolePolicy{entity.RoleEditor}).
			Policy(RolePolicy{entity.RoleAdmin}).
			Strategy(AtLeastOneStrategy).
			Resolve(entity.Album{}, session.User)

		if !hasPermission {
			common.AbortForbidden(c, errors.New("user has no editor of admin role"), "user role forbids the creation of albums")

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

		albumID, err := albumRepo.Create(reqCtx, album)
		if err != nil {
			common.AbortInternalError(c, err, fmt.Sprintf("album: %+v", album))

			return
		}

		logger.WithFields(logrus.Fields{
			"album": fmt.Sprintf("%+v", album),
			"id":    albumID,
		}).Info("album entity created")

		c.Redirect(http.StatusFound, rootURL)
	})
}

// GET /album/:id/edit
func GetUpdateAlbumForm(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	keycloakRepo := repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	r.GET("/album/:id/edit", func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		// decrypt album id
		gen := encryption.NewGenerator(conf.GetEncryptionKey())

		decryptedID, err := gen.DecryptData(c.Param("id"))
		if err != nil {
			logger.WithError(err).Error("cannot decrypt album id")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		id, err := strconv.Atoi(decryptedID)
		if err != nil {
			logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		album, err := albumRepo.GetByID(reqCtx, int32(id))
		if err != nil {
			common.AbortNotFound(c, err, "update album")

			return
		}

		// check if user is the owner or it has the edit permission set
		if album.OwnerID == session.User.ID || session.User.Role == entity.RoleAdmin {
			logger.Info("edit permission granted. user is the owner")

			// filter out current user, admins and with can_share false
			userFilters, err := generateFilters(session.User)
			if err != nil {
				logutil.GetLogger(c).WithError(err).Error("create user fitlers")
				common.AbortInternalError(c, err, "")

				return
			}

			users, err := keycloakRepo.GetUsers(reqCtx, userFilters...)
			if err != nil && errors.Is(err, domain.ErrInternalError) {
				common.AbortInternalError(c, err, "cannot fetch users")

				return
			}

			serializedUsers := serialize(users)
			userPermissions := make(entity.Permissions)
			for _, u := range serializedUsers {
				if perm, found := album.UserPermissions[u.ID]; found {
					userPermissions[u.EncryptedID] = perm
				}
			}

			groups, err := keycloakRepo.GetGroups(reqCtx)
			if err != nil && errors.Is(err, domain.ErrInternalError) {
				common.AbortInternalError(c, err, "cannot fetch groups")

				return
			}

			var permUserJson string
			if len(album.UserPermissions) > 0 {
				permUserJson, err = userPermissions.Json()
				if err != nil {
					logger.WithField("user permissions", fmt.Sprintf("%+v", userPermissions)).WithError(err).Error("marshal to json")
				}
			}

			var permGroupJson string
			if len(album.GroupPermissions) > 0 {
				var err error
				permGroupJson, err = album.GroupPermissions.Json()
				if err != nil {
					logger.WithField("group permissions", fmt.Sprintf("%+v", album.GroupPermissions)).WithError(err).Error("marshal to json")
				}
			}

			encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", album.ID))
			if err != nil {
				logger.WithError(err).WithField("id", album.ID).Error("encrypt album id")
				common.AbortInternalError(c, err, "")

				return
			}

			c.HTML(http.StatusOK, "album_form.html", gin.H{
				"update_link":        fmt.Sprintf("/album/%s", encryptedID),
				"album":              album,
				"canShare":           session.User.CanShare,
				"isOwner":            true,
				"users":              serializedUsers,
				"groups":             groups,
				"users_permissions":  permUserJson,
				"groups_permissions": permGroupJson,
				"is_admin":           session.User.Role == entity.RoleAdmin,
				csrf.TemplateTag:     csrf.TemplateField(c.Request),
			})

			return
		}

		// only users with editPermission set for this album or one of user's group with the same permission
		// can edit this album
		apr := NewAlbumPermissionResolver()
		hasPermission := apr.Policy(UserPermissionPolicy{entity.PermissionEditAlbum}).
			Policy(GroupPermissionPolicy{entity.PermissionEditAlbum}).
			Strategy(AtLeastOneStrategy).
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
			"album":          album,
			"canShare":       session.User.CanShare,
			"isOwner":        false,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	})

}

// PUT /album/:id/
func UpdateAlbum(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)

	r.POST("/album/:id/", func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		// decrypt album id
		gen := encryption.NewGenerator(conf.GetEncryptionKey())

		decryptedID, err := gen.DecryptData(c.Param("id"))
		if err != nil {
			logger.WithError(err).Error("cannot decrypt album id")
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		id, err := strconv.Atoi(decryptedID)
		if err != nil {
			logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		album, err := albumRepo.GetByID(reqCtx, int32(id))
		if err != nil {
			common.AbortNotFound(c, err, "update album")

			return
		}

		var albumForm form.Album
		if err := c.ShouldBind(&albumForm); err != nil {
			common.AbortBadRequest(c, err, "fail to bind to form")

			return
		}

		// only users with editPermission set for this album or one of user's group with the same permission
		// can edit this album
		apr := NewAlbumPermissionResolver()
		hasPermission := apr.Policy(OwnerPolicy{}).
			Policy(RolePolicy{entity.RoleAdmin}).
			Policy(UserPermissionPolicy{entity.PermissionEditAlbum}).
			Policy(GroupPermissionPolicy{entity.PermissionEditAlbum}).
			Strategy(AtLeastOneStrategy).
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

		// edit permissions don't allow a user other than the owner to change permissions
		err = albumRepo.Update(reqCtx, album)
		if err != nil {
			common.AbortInternalError(c, err, "update album")

			return
		}

		c.Redirect(http.StatusFound, rootURL)
	})
}

// DELETE /album/:id
func DeleteAlbum(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)

	r.DELETE("/album/:id", func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		// decrypt album id
		gen := encryption.NewGenerator(conf.GetEncryptionKey())

		decryptedID, err := gen.DecryptData(c.Param("id"))
		if err != nil {
			logger.WithError(err).Error("cannot decrypt album id")
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		id, err := strconv.Atoi(decryptedID)
		if err != nil {
			logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		album, err := albumRepo.GetByID(reqCtx, int32(id))
		if err != nil {
			common.AbortNotFound(c, err, "update album")

			return
		}

		// only users with editPermission set for this album or one of user's group with the same permission
		// can edit this album
		apr := NewAlbumPermissionResolver()
		hasPermission := apr.Policy(OwnerPolicy{}).
			Policy(UserPermissionPolicy{entity.PermissionDeleteAlbum}).
			Policy(GroupPermissionPolicy{entity.PermissionDeleteAlbum}).
			Strategy(AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"request user id": session.User.ID,
				"album owner id":  album.OwnerID,
			}).Error("album can be edit either by user with delete permission or the owner")
			common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionDeleteAlbum, album, session.User), "delete album")

			return
		}

		err = albumRepo.Delete(reqCtx, album.ID)
		if err != nil {
			common.AbortInternalError(c, common.ErrDeleteAlbum, fmt.Sprintf("album id: %d", id))

			return
		}

		c.Redirect(http.StatusFound, rootURL)
	})
}
