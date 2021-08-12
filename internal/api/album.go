package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

func GetAlbums(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	keycloakRepo := repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	r.GET("/api/albums", func(c *gin.Context) {
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

		personalAlbums, err := albumRepo.GetByOwnerID(reqCtx, session.User.ID, nil)
		if err != nil {
			if errors.Is(err, domain.ErrAlbumNotFound) {
				logger.Info("user has no personal albums")
			} else {
				panic(err) // TODO 500 page
			}
		}

		personalTAlbums := mapUsersToAlbums(personalAlbums, users)

		// if I'm an admin show all other albums
		if session.User.Role == entity.RoleAdmin {
			otherAlbums, err := albumRepo.Get(reqCtx, nil)
			if err != nil {
				logger.WithError(err).Error("fetch all albums")
				AbortInternalError(c, err, "")

				return
			}

			sharedAlbums := substractAlbums(otherAlbums, personalAlbums)
			sharedTAlbums := mapUsersToAlbums(sharedAlbums, users)

			c.HTML(http.StatusOK, "index.html", gin.H{
				"username":        session.User.Username,
				"user_role":       session.User.Role.String(),
				"personal_albums": personalTAlbums,
				"shared_albums":   sharedTAlbums,
			})

			return
		}

		// user with canShare true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		if session.User.CanShare {
			sharedUserAlbums, err := albumRepo.GetByUserID(reqCtx, session.User.ID, nil)
			if err != nil {
				if errors.Is(err, domain.ErrAlbumNotFound) {
					logger.Info("user has no shared albums")
				} else {
					panic(err) // TODO 500 page
				}
			}

			// get albums shared by group permissions
			groupSharedAlbums := make([]entity.Album, 0)
			for _, g := range session.User.Groups {
				groupAlbums, err := albumRepo.GetByGroupName(reqCtx, g.Name, nil)
				if err != nil {
					logger.WithError(err).WithField("group name", g.Name).Error("fetch by group name")

					continue
				}

				groupSharedAlbums = append(groupSharedAlbums, groupAlbums...)
			}

			// remove personal groups which are in group shared albums
			groupUniqueAlbums := substractAlbums(groupSharedAlbums, personalAlbums)

			// remove duplicates
			sharedAlbums := addGroups(sharedUserAlbums, groupUniqueAlbums)

			// map them to template album struct
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
