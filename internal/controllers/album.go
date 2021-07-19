package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

// POST /album
func CreateAlbum(r *gin.RouterGroup, repos repo.Repositories) {
	//	albumRepo := repos[repo.AlbumRepoName].(repo.AlbumRepo)
	userRepo := repos[repo.UserRepoName].(repo.UserRepo)
	groupRepo := repos[repo.GroupRepoName].(repo.GroupRepo)

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
			"users":  filteredUsers,
			"groups": groups,
		})
	})

	r.POST("/album", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		//reqCtx := c.Request.Context()
		//logger := logutil.GetLogger(c)

		// only editors and admins have the right to create albums
		if session.User.Role == entity.RoleUser {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("user with user role cannot create albums"))
		}

		return
	})
}
