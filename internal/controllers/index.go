package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/encryption"
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

		personalTAlbums := mapUsersToAlbums(personalAlbums, users)

		// user with canShare true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		if session.User.CanShare {
			sharedUserAlbums, err := albumRepo.GetByUserID(reqCtx, session.User.ID)
			if err != nil {
				if errors.Is(err, repo.ErrAlbumNotFound) {
					logger.Info("user has no shared albums")
				} else {
					panic(err) // TODO 500 page
				}
			}

			// get albums shared by group permissions
			groupSharedAlbums := make([]entity.Album, 0)
			for _, g := range session.User.Groups {
				groupAlbums, err := albumRepo.GetByGroupName(reqCtx, g.Name)
				if err != nil {
					logger.WithError(err).WithField("group name", g.Name).Error("fetch by group name")

					continue
				}

				groupSharedAlbums = append(groupSharedAlbums, groupAlbums...)
			}

			// remove duplicates
			sharedAlbums := removeDuplicates(sharedUserAlbums, groupSharedAlbums)

			sharedTAlbums := mapUsersToAlbums(sharedAlbums, users)

			c.HTML(http.StatusOK, "index.html", gin.H{
				"username":        session.User.Username,
				"user_role":       session.User.Role.String(),
				"personal_albums": personalTAlbums,
				"shared_albums":   sharedTAlbums,
			})

			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"username":        session.User.Username,
			"user_role":       session.User.Role.String(),
			"personal_albums": personalTAlbums,
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

func removeDuplicates(albums1, albums2 []entity.Album) []entity.Album {
	album1Map := make(map[int32]entity.Album)
	album2Map := make(map[int32]entity.Album)

	var iLimit int
	if len(albums1) > len(albums2) {
		iLimit = len(albums1)
	} else {
		iLimit = len(albums2)
	}

	for i := 0; i < iLimit; i++ {
		if i < len(albums1) {
			album1Map[albums1[i].ID] = albums1[i]
		}

		if i < len(albums2) {
			album2Map[albums2[i].ID] = albums2[i]
		}
	}

	distinctAlbums := make([]entity.Album, 0, len(albums1)+len(albums2))
	for k, v := range album1Map {
		if _, found := album2Map[k]; found {
			distinctAlbums = append(distinctAlbums, v)
			delete(album1Map, k)
			delete(album2Map, k)
		}
	}

	// put the rest from album2Map
	for _, v := range album2Map {
		distinctAlbums = append(distinctAlbums, v)
	}

	return distinctAlbums
}

func mapUsersToAlbums(albums []entity.Album, users map[string]entity.User) []tAlbum {
	tAlbums := make([]tAlbum, 0, len(albums))
	for _, pa := range albums {
		if user, found := users[pa.OwnerID]; found {
			talbum, err := mapTAbum(pa, user)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("album", fmt.Sprintf("%+v", pa)).Error("cannot map shared album")

				continue
			}

			tAlbums = append(tAlbums, talbum)
		} else {
			logutil.GetDefaultLogger().WithField("album", fmt.Sprintf("%+v", pa.String())).Warn("owner don't exists anymore")
		}
	}

	return tAlbums
}

func mapTAbum(a entity.Album, user entity.User) (tAlbum, error) {
	var owner string

	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
		owner = user.Username
	} else {
		owner = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	// encrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", a.ID))
	if err != nil {
		return tAlbum{}, err
	}

	logutil.GetDefaultLogger().WithFields(logrus.Fields{
		"id":           a.ID,
		"encrypted_id": encryptedID,
	}).Trace("encrypt album id")

	return tAlbum{
		ID:          encryptedID,
		Name:        a.Name,
		Date:        a.CreatedAt,
		Location:    a.Location,
		Description: a.Description,
		Owner:       owner,
	}, nil
}

type tAlbum struct {
	ID          string
	Name        string
	Owner       string
	Date        time.Time
	Description string
	Location    string
}
