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
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/form"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

// POST /album
func CreateAlbum(r *gin.RouterGroup, repos Repositories) {
	albumRepo := repos[AlbumRepoName].(AlbumRepo)
	keycloakRepo := repos[KeycloakRepoName].(KeycloakRepo)

	r.GET("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))

			return
		}

		users, err := keycloakRepo.GetUsers(reqCtx)
		if err != nil && errors.Is(err, repo.ErrInternalError) {
			logger.WithError(err).Error("cannot fetch users")
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot fetch users"))

			return
		}

		// filter out admins and can_share is false
		userFilter := entity.NewUserFilter(users)
		filteredUsers := userFilter.Filter(func(u entity.User) bool {
			return u.CanShare == true && u.Role != entity.RoleAdmin && u.Username != session.User.Username
		})

		groups, err := keycloakRepo.GetGroups(reqCtx)
		if err != nil && errors.Is(err, repo.ErrInternalError) {
			logger.WithError(err).Error("cannot fetch groups")
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot fetch groups"))

			return
		}

		c.HTML(http.StatusOK, "album_form.html", gin.H{
			"users":          filteredUsers,
			"groups":         groups,
			"canShare":       session.User.CanShare,
			"isOwner":        true,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		})
	})

	// TODO CRSF protection
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
			logger.WithError(err).Info("fail to bind to json")

			c.HTML(http.StatusBadRequest, "album_form.html", gin.H{"error": err})

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
				logger.WithError(err).Error("error fetching users")
				c.HTML(http.StatusInternalServerError, "internal_error.html", nil)

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
				logger.WithError(err).Error("error fetching groups")
				c.HTML(http.StatusInternalServerError, "internal_error.html", nil)

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
			logger.WithField("album", fmt.Sprintf("%+v", album)).WithError(err).Error("cannot create album")
			c.Redirect(http.StatusInternalServerError, "/error")

			return
		}

		logger.WithFields(logrus.Fields{
			"album": fmt.Sprintf("%+v", album),
			"id":    albumID,
		}).Info("album entity created")

		c.Redirect(http.StatusFound, "/")
	})
}

func UpdateAlbum(r *gin.RouterGroup, repos Repositories) {
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
			logger.WithError(err).WithField("id", id).Error("update album")
			c.AbortWithError(http.StatusNotFound, err)

			return
		}

		// check if user is the owner or it has the edit permission set
		if album.OwnerID == session.User.ID {
			logger.Info("edit permission granted. user is the owner")

			users, err := keycloakRepo.GetUsers(reqCtx)
			if err != nil && errors.Is(err, repo.ErrInternalError) {
				logger.WithError(err).Error("cannot fetch users")
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot fetch users"))

				return
			}

			// filter out admins and can_share is false
			userFilter := entity.NewUserFilter(users)
			filteredUsers := userFilter.Filter(func(u entity.User) bool {
				return u.CanShare == true && u.Role != entity.RoleAdmin && u.Username != session.User.Username
			})

			groups, err := keycloakRepo.GetGroups(reqCtx)
			if err != nil && errors.Is(err, repo.ErrInternalError) {
				logger.WithError(err).Error("cannot fetch groups")
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot fetch groups"))

				return
			}

			c.HTML(http.StatusOK, "album_form.html", gin.H{
				"album":    album,
				"canShare": session.User.CanShare,
				"isOwner":  true,
				"users":    filteredUsers,
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

	r.PUT("/album/:id/", func(c *gin.Context) {

	})

}

// parsePermissions will parse the permission string (e.g. (username#r,w)(uername2#e,d))
func parsePermissions(perms string) map[string][]entity.Permission {
	permRe := regexp.MustCompile(`(\((\w+)#(([rwed],?)+)\))`)
	permissions := make(map[string][]entity.Permission)

	for matchIdx, match := range permRe.FindAllStringSubmatch(perms, -1) {
		logutil.GetDefaultLogger().WithFields(logrus.Fields{"idx": matchIdx, "match": fmt.Sprintf("%+v", match)}).Debug("permission matched")
		// get 2nd and 3rd groups only
		name := match[2]

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
