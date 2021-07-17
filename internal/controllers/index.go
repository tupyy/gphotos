package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, albumRepo repo.AlbumRepo) {
	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		ownAlbums, err := albumRepo.GetByOwnerID(reqCtx, *session.User.ID)
		if err != nil {
			if errors.Is(repo.ErrAlbumNotFound, err) {
				logger.Info("albums not found")
			} else {
				panic(err) // TODO 500 page
			}
		}

		// user with canShare true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		var sharedAlbums []entity.Album

		if session.User.CanShare {
			sharedAlbums, err = albumRepo.GetByUserID(reqCtx, *session.User.ID)
			if err != nil {
				if errors.Is(repo.ErrAlbumNotFound, err) {
					logger.Info("user has no shared albums")
				} else {
					panic(err) // TODO 500 page
				}
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"username":     session.User.Username,
			"user_role":    session.User.Role.String(),
			"ownAlbums":    ownAlbums,
			"sharedAlbums": sharedAlbums,
		})
	})
}
