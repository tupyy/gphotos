package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, albumRepo AlbumRepo) {
	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		personalAlbums, err := albumRepo.GetByOwnerID(reqCtx, session.User.ID)
		if err != nil {
			if errors.Is(err, repo.ErrAlbumNotFound) {
				logger.Info("user has no personal albums")
			} else {
				panic(err) // TODO 500 page
			}
		}

		// user with canShare true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		var sharedAlbums []entity.Album

		if session.User.CanShare {
			sharedAlbums, err = albumRepo.GetByUserID(reqCtx, session.User.ID)
			if err != nil {
				if errors.Is(err, repo.ErrAlbumNotFound) {
					logger.Info("user has no shared albums")
				} else {
					panic(err) // TODO 500 page
				}
			}
		}

		// TODO add album with permissions given by user's group

		c.HTML(http.StatusOK, "index.html", gin.H{
			"username":       session.User.Username,
			"user_role":      session.User.Role.String(),
			"personalAlbums": personalAlbums,
			"sharedAlbums":   sharedAlbums,
		})
	})
}

//type tAlbum struct {
//	Name        string
//	Owner       string
//	Description string
//	Date        time.Time
//	Location    string
//}
//
//func newTAlbum(a entity.Album) tAlbum {
//	t := tAlbum{
//		Name:  a.Name,
//		Date:  a.CreatedAt,
//		Owner: string(a.OwnerID),
//	}
//
//	if a.Description != nil {
//		t.Description = *a.Description
//	}
//
//	if a.Location != nil {
//		t.Location = *a.Location
//	}
//}
