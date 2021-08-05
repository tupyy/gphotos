package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, repos repo.Repositories) {
	albumRepo := repos[repo.AlbumRepoName].(repo.Album)
	keycloakRepo := repos[repo.KeycloakRepoName].(repo.KeycloakRepo)

	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		// get users to fetch name of the user
		users, err := getUsers(reqCtx, keycloakRepo)
		if err != nil {
			logger.WithError(err).Error("index fetch users")
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		personalAlbums, err := albumRepo.GetByOwnerID(reqCtx, session.User.ID)
		if err != nil {
			if errors.Is(err, repo.ErrAlbumNotFound) {
				logger.Info("user has no personal albums")
			} else {
				panic(err) // TODO 500 page
			}
		}

		personalTAlbums := make([]tAlbum, 0, len(personalAlbums))
		for _, pa := range personalAlbums {
			if user, found := users[pa.OwnerID]; found {
				personalTAlbums = append(personalTAlbums, mapTAbum(pa, user))
			} else {
				logger.WithField("album", fmt.Sprintf("%+v", pa.String())).Warn("owner don't exists anymore")
			}
		}

		// user with canShare true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		if session.User.CanShare {
			sharedAlbums, err := albumRepo.GetByUserID(reqCtx, session.User.ID)
			if err != nil {
				if errors.Is(err, repo.ErrAlbumNotFound) {
					logger.Info("user has no shared albums")
				} else {
					panic(err) // TODO 500 page
				}
			}

			sharedTAlbums := make([]tAlbum, 0, len(personalAlbums))
			for _, pa := range sharedAlbums {
				if user, found := users[pa.OwnerID]; found {
					sharedTAlbums = append(sharedTAlbums, mapTAbum(pa, user))
				} else {
					logger.WithField("album", fmt.Sprintf("%+v", pa.String())).Warn("owner don't exists anymore")
				}
			}

			c.HTML(http.StatusOK, "index.html", gin.H{
				"username":       session.User.Username,
				"user_role":      session.User.Role.String(),
				"personalAlbums": personalTAlbums,
				"sharedAlbums":   sharedTAlbums,
			})

			return
		}

		// TODO add album with permissions given by user's group

		c.HTML(http.StatusOK, "index.html", gin.H{
			"username":       session.User.Username,
			"user_role":      session.User.Role.String(),
			"personalAlbums": personalTAlbums,
		})
	})
}

func getUsers(ctx context.Context, k repo.KeycloakRepo) (map[string]entity.User, error) {
	users, err := k.GetUsers(ctx)
	if err != nil {
		return nil, err

	}

	mappedUsers := make(map[string]entity.User)
	for _, u := range users {
		mappedUsers[u.ID] = u
	}

	return mappedUsers, nil
}

func mapTAbum(a entity.Album, user entity.User) tAlbum {
	var owner string

	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
		owner = user.Username
	} else {
		owner = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	return tAlbum{
		ID:          a.ID,
		Name:        a.Name,
		Date:        a.CreatedAt,
		Location:    a.Location,
		Description: a.Description,
		Owner:       owner,
	}
}

type tAlbum struct {
	ID          int32
	Name        string
	Owner       string
	Date        time.Time
	Description string
	Location    string
}
