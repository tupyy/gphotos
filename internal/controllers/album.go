package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/form"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// GET /album/:id
func GetAlbum(r *gin.RouterGroup, repos Repositories) {
	//albumRepo := repos[AlbumRepoName].(AlbumRepo)

	r.GET("/album/:id", func(c *gin.Context) {

	})
}

// GET /album
func GetCreateAlbumForm(r *gin.RouterGroup, repos Repositories) {
	keycloakRepo := repos[KeycloakRepoName].(KeycloakRepo)

	r.GET("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		reqCtx := c.Request.Context()

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))

			return
		}

		users, err := keycloakRepo.GetUsers(reqCtx)
		if err != nil && errors.Is(err, repo.ErrInternalError) {
			AbortInternalError(c, err, "cannot fetch users")

			return
		}

		// filter out admins and can_share is false
		userFilter := entity.NewUserFilter(users)
		filteredUsers := userFilter.Filter(func(u entity.User) bool {
			return u.CanShare == true && u.Role != entity.RoleAdmin && u.Username != session.User.Username
		})

		userMap, err := mapNames(filteredUsers)
		if err != nil {
			AbortInternalError(c, err, "cannot encrypt usernames")

			return
		}

		groups, err := keycloakRepo.GetGroups(reqCtx)
		if err != nil && errors.Is(err, repo.ErrInternalError) {
			AbortInternalError(c, err, "cannot fetch groups")

			return
		}

		c.HTML(http.StatusOK, "album_form.html", gin.H{
			"users":          userMap,
			"groups":         groups,
			"canShare":       session.User.CanShare,
			"isOwner":        true,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	})
}

// POST /album
func CreateAlbum(r *gin.RouterGroup, repos Repositories) {
	albumRepo := repos[AlbumRepoName].(AlbumRepo)
	keycloakRepo := repos[KeycloakRepoName].(KeycloakRepo)

	r.POST("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))
		}

		var albumForm form.Album
		if err := c.ShouldBind(&albumForm); err != nil {
			AbortBadRequest(c, err, "fail to bind to form")

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
			album.UserPermissions = make(map[string][]entity.Permission)

			// get all the users
			users, err := keycloakRepo.GetUsers(reqCtx)
			if err != nil {
				AbortInternalError(c, err, "error fetching users")

				return
			}

			// put users into a map
			usersID := make(map[string]string)
			for _, u := range users {
				// remove the current user
				if u.ID != album.OwnerID {
					usersID[u.Username] = u.ID
				}
			}

			perms := parsePermissions(cleanForm.UserPermissions)

			if len(perms) == 0 {
				logger.WithField("permissions_string", cleanForm.UserPermissions).Warn("cannot user parse permissions")
			} else {
				for k, v := range perms {
					if userID, found := usersID[k]; found {
						logger.WithFields(logrus.Fields{
							"userID":      userID,
							"permissions": v,
						}).Trace("permissions added")

						album.UserPermissions[userID] = v
					} else {
						logger.WithField("username", k).Warn("username not found in db")
					}
				}
			}
		}

		if len(cleanForm.GroupPermissions) > 0 {
			album.GroupPermissions = make(map[string][]entity.Permission)

			// get all the users
			groups, err := keycloakRepo.GetGroups(reqCtx)
			if err != nil {
				AbortInternalError(c, err, "error fetching groups")

				return
			}

			// put groups into a map
			groupsID := make(map[string]string)
			for _, g := range groups {
				groupsID[g.Name] = g.Name
			}

			perms := parsePermissions(cleanForm.GroupPermissions)

			if len(perms) == 0 {
				logger.WithField("permissions_string", cleanForm.GroupPermissions).Warn("cannot group parse permissions")
			} else {
				for k, v := range perms {
					if groupID, found := groupsID[k]; found {
						album.GroupPermissions[groupID] = v
					} else {
						logger.WithField("group name", k).Warn("group not found in db")
					}
				}
			}
		}

		albumID, err := albumRepo.Create(reqCtx, album)
		if err != nil {
			AbortInternalError(c, err, fmt.Sprintf("album: %+v", album))

			return
		}

		logger.WithFields(logrus.Fields{
			"album": fmt.Sprintf("%+v", album),
			"id":    albumID,
		}).Info("album entity created")

		c.Redirect(http.StatusFound, "/")
	})
}

// GET /album/:id/edit
func GetUpdateAlbumForm(r *gin.RouterGroup, repos Repositories) {
	albumRepo := repos[AlbumRepoName].(AlbumRepo)
	keycloakRepo := repos[KeycloakRepoName].(KeycloakRepo)

	r.GET("/album/:id/edit", func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		param := c.Param("id")

		id, err := strconv.Atoi(param)
		if err != nil {
			logger.WithError(err).WithField("id", param).Error("cannot parse album id")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		album, err := albumRepo.GetByID(reqCtx, int32(id))
		if err != nil {
			AbortNotFound(c, err, "update album")

			return
		}

		// check if user is the owner or it has the edit permission set
		if album.OwnerID == session.User.ID {
			logger.Info("edit permission granted. user is the owner")

			users, err := keycloakRepo.GetUsers(reqCtx)
			if err != nil && errors.Is(err, repo.ErrInternalError) {
				AbortInternalError(c, err, "cannot fetch users")

				return
			}

			// filter out admins and can_share is false
			userFilter := entity.NewUserFilter(users)
			filteredUsers := userFilter.Filter(func(u entity.User) bool {
				return u.CanShare == true && u.Role != entity.RoleAdmin && u.Username != session.User.Username
			})

			userMap, err := mapNames(filteredUsers)
			if err != nil {
				AbortInternalError(c, err, "cannot encrypt usernames")

				return
			}

			groups, err := keycloakRepo.GetGroups(reqCtx)
			if err != nil && errors.Is(err, repo.ErrInternalError) {
				AbortInternalError(c, err, "cannot fetch groups")

				return
			}

			c.HTML(http.StatusOK, "album_form.html", gin.H{
				"album":    album,
				"canShare": session.User.CanShare,
				"isOwner":  true,
				"users":    userMap,
				"groups":   groups,
			})

			return
		}

		editPermissionFound := false

		// the user is not the owner
		// check if user has edit permission set
		if album.HasUserPermission(session.User.ID, entity.PermissionEditAlbum) {
			logger.Info("edit permission granted. user has given the edit permission by the owner.")
			editPermissionFound = true
		}

		// check if one of user's groups has edit permission set
		for _, group := range session.User.Groups {
			if album.HasGroupPermission(group.Name, entity.PermissionEditAlbum) {
				logger.Info("edit permission granted. user's group has given the edit permission.")
				editPermissionFound = true
				break
			}
		}

		if !editPermissionFound {
			logger.WithFields(logrus.Fields{
				"request user id": session.User.ID,
				"album owner id":  album.OwnerID,
			}).Error("album cannot be edit either by user with edit permission or the owner")
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("user has no edit permissions for this album"))

			return
		}

		c.HTML(http.StatusOK, "album_form.html", gin.H{
			"album":    album,
			"canShare": session.User.CanShare,
			"isOwner":  false,
		})
	})

}

// PUT /album/:id/
func UpdateAlbum(r *gin.RouterGroup, repos Repositories) {
	//albumRepo := repos[AlbumRepoName].(AlbumRepo)
	//keycloakRepo := repos[KeycloakRepoName].(KeycloakRepo)

	r.PUT("/album/:id/", func(c *gin.Context) {

	})
}

// DELETE /album/:id
func DeleteAlbum(r *gin.RouterGroup, repos Repositories) {
	r.DELETE("/album/:id", func(c *gin.Context) {

	})
}

// parsePermissions will parse the permission string (e.g. (username#r,w)(uername2#e,d))
func parsePermissions(perms string) map[string][]entity.Permission {
	permRe := regexp.MustCompile(`(\((\w+)#(([rwed],?)+)\))`)
	permissions := make(map[string][]entity.Permission)

	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	for matchIdx, match := range permRe.FindAllStringSubmatch(perms, -1) {
		logutil.GetDefaultLogger().WithFields(logrus.Fields{"idx": matchIdx, "match": fmt.Sprintf("%+v", match)}).Debug("permission matched")
		// get 2nd and 3rd groups only
		name, err := gen.DecryptData(match[2])
		if err != nil {
			logutil.GetDefaultLogger().WithError(err).WithField("data", match[2]).Error("decrypt name")

			continue
		}

		permList := strings.Split(match[3], ",")
		entities := make([]entity.Permission, 0, len(permList))

		for _, p := range permList {
			switch p {
			case "r":
				entities = append(entities, entity.PermissionReadAlbum)
			case "w":
				entities = append(entities, entity.PermissionWriteAlbum)
			case "e":
				entities = append(entities, entity.PermissionEditAlbum)
			case "d":
				entities = append(entities, entity.PermissionDeleteAlbum)
			}
		}

		if len(entities) > 0 {
			permissions[name] = entities
		}
	}

	return permissions
}

// encode the permission map into a string like (e.g. (username#r,w)(uername2#e,d))
func encodePermissions(perms map[string][]entity.Permission) string {
	return ""
}

// return a map with encrypted username as key and First + Last name as value
func mapNames(users []entity.User) (map[string]string, error) {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	encryptedUsernames := make(map[string]string)
	for _, u := range users {
		encryptedUsername, err := gen.EncryptData(u.Username)
		if err != nil {
			return nil, err
		}

		encryptedUsernames[encryptedUsername] = fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	}

	return encryptedUsernames, nil
}
