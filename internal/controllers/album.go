package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/form"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

// POST /album
func CreateAlbum(r *gin.RouterGroup, repos Repositories) {
	albumRepo := repos[AlbumRepoName].(AlbumRepo)
	userRepo := repos[UserRepoName].(UserRepo)
	groupRepo := repos[GroupRepoName].(GroupRepo)

	r.GET("/album/new", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))

			return
		}

		users, err := userRepo.Get(reqCtx)
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

		groups, err := groupRepo.Get(reqCtx)
		if err != nil && errors.Is(err, repo.ErrInternalError) {
			logger.WithError(err).Error("cannot fetch groups")
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot fetch groups"))

			return
		}

		c.HTML(http.StatusOK, "album_create.html", gin.H{
			"users":    filteredUsers,
			"groups":   groups,
			"canShare": session.User.CanShare,
		})
	})

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

			c.HTML(http.StatusBadRequest, "album_create.html", gin.H{"error": err})

			return
		}

		escapedFormAlbum := albumForm.Sanitize()

		logger.WithField("form", fmt.Sprintf("%+v", escapedFormAlbum)).Info("create album request submitted")

		album := entity.Album{
			Name:        escapedFormAlbum.Name,
			Description: &escapedFormAlbum.Description,
			CreatedAt:   time.Now(),
			Location:    &escapedFormAlbum.Location,
			OwnerID:     *session.User.ID,
		}

		if len(escapedFormAlbum.UserPermissions) > 0 {
			album.UserPermissions = make(map[int32][]entity.Permission)

			// get all the users
			users, err := userRepo.Get(reqCtx)
			if err != nil {
				logger.WithError(err).Error("error fetching users")
				c.HTML(http.StatusInternalServerError, "internal_error.html", nil)

				return
			}

			// put users into a map
			usersID := make(map[string]int32)
			for _, u := range users {
				usersID[u.Username] = *u.ID
			}

			perms := parsePermissions(escapedFormAlbum.UserPermissions)

			if len(perms) == 0 {
				logger.WithField("permissions_string", escapedFormAlbum.UserPermissions).Warn("cannot user parse permissions")
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

		if len(escapedFormAlbum.GroupPermissions) > 0 {
			album.GroupPermissions = make(map[int32][]entity.Permission)

			// get all the users
			groups, err := groupRepo.Get(reqCtx)
			if err != nil {
				logger.WithError(err).Error("error fetching groups")
				c.HTML(http.StatusInternalServerError, "internal_error.html", nil)

				return
			}

			// put groups into a map
			groupsID := make(map[string]int32)
			for _, u := range groups {
				groupsID[u.Name] = *u.ID
			}

			perms := parsePermissions(escapedFormAlbum.GroupPermissions)

			if len(perms) == 0 {
				logger.WithField("permissions_string", escapedFormAlbum.GroupPermissions).Warn("cannot group parse permissions")
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

func parsePermissions(perms string) map[string][]entity.Permission {
	permissions := make(map[string][]entity.Permission)

	// for _, p := range perms {
	// 	switch p {
	// 	case "r":
	// 		permissions = append(permissions, entity.PermissionReadAlbum)
	// 	case "w":
	// 		permissions = append(permissions, entity.PermissionWriteAlbum)
	// 	case "e":
	// 		permissions = append(permissions, entity.PermissionEditAlbum)
	// 	case "d":
	// 		permissions = append(permissions, entity.PermissionDeleteAlbum)
	// 	}
	// }

	return permissions
}
