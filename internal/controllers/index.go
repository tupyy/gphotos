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

		ownAlbums, err := albumRepo.GetAlbumsByOwnerID(reqCtx, *session.User.ID)
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
			sharedAlbums, err = albumRepo.GetAlbumsByUserID(reqCtx, *session.User.ID)
			if err != nil {
				if errors.Is(repo.ErrAlbumNotFound, err) {
					logger.Info("user has no shared albums")
				} else {
					panic(err) // TODO 500 page
				}
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"hasOwnAlbums":    len(ownAlbums) != 0,
			"ownAlbums":       ownAlbums,
			"hasSharedAlbums": len(sharedAlbums) != 0,
			"sharedAlbums":    sharedAlbums,
		})
	})
}
