package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	domainFilter "github.com/tupyy/gophoto/internal/domain/filters"
	domainSort "github.com/tupyy/gophoto/internal/domain/sort"
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

		// fetch users from keycloak
		users, err := keycloakRepo.GetUsers(reqCtx, nil)
		if err != nil {
			logger.WithError(err).Error("index fetch users")
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		personalAlbums, err := albumRepo.GetByOwnerID(reqCtx, session.User.ID, domainSort.NewAlbumSorterByDate(domainSort.ReverseOrder))
		if err != nil {
			if errors.Is(err, domain.ErrAlbumNotFound) {
				logger.Info("user has no personal albums")
			} else {
				panic(err) // TODO 500 page
			}
		}

		notInPersonalAlbumsFilter, err := domainFilter.GenerateAlbumFilterFuncs(domainFilter.FilterNotInList, personalAlbums)
		if err != nil {
			logutil.GetLogger(c).WithError(err).Error("filter error")
			AbortWithJson(c, http.StatusInternalServerError, err, "")

			return
		}

		// if I'm an admin show all other albums
		if session.User.Role == entity.RoleAdmin {
			sharedAlbums, err := albumRepo.Get(reqCtx, domainSort.NewAlbumSorterByDate(domainSort.ReverseOrder), notInPersonalAlbumsFilter)
			if err != nil {
				logger.WithError(err).Error("fetch all albums")
				AbortInternalError(c, err, "")

				return
			}

			c.JSON(http.StatusOK, gin.H{
				"username":        session.User.Username,
				"user_role":       session.User.Role.String(),
				"personal_albums": newSimpleAlbums(personalAlbums, users),
				"shared_albums":   newSimpleAlbums(sharedAlbums, users),
			})

			return
		}

		// user with canShare true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		if session.User.CanShare {
			sharedUserAlbums, err := albumRepo.GetByUserID(reqCtx, session.User.ID, domainSort.NewAlbumSorterByDate(domainSort.ReverseOrder))
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
				// get albums by group name. Sort them by date in reverse order and filter out the personal albums.
				// It can happens that the user shared albums with its own groups.
				notInPersonalAlbumsFilter, err := domainFilter.GenerateAlbumFilterFuncs(domainFilter.FilterNotInList, personalAlbums)
				if err != nil {
					logutil.GetLogger(c).WithError(err).Error("filter error")
					AbortWithJson(c, http.StatusInternalServerError, err, "")

					return
				}

				// In case when the album is shared by user AND group, keep only one copy.
				notInSharedAlbumsFilter, err := domainFilter.GenerateAlbumFilterFuncs(domainFilter.FilterNotInList, sharedUserAlbums)
				if err != nil {
					logutil.GetLogger(c).WithError(err).Error("filter error")
					AbortWithJson(c, http.StatusInternalServerError, err, "")

					return
				}

				groupAlbums, err := albumRepo.GetByGroupName(reqCtx,
					g.Name,
					domainSort.NewAlbumSorterByDate(domainSort.ReverseOrder),
					notInPersonalAlbumsFilter,
					notInSharedAlbumsFilter)
				if err != nil {
					logger.WithError(err).WithField("group name", g.Name).Error("fetch by group name")

					continue
				}

				groupSharedAlbums = append(groupSharedAlbums, groupAlbums...)
			}

			sharedAlbums := joinAlbums(sharedUserAlbums, groupSharedAlbums)

			c.JSON(http.StatusOK, gin.H{
				"username":        session.User.Username,
				"user_role":       session.User.Role.String(),
				"personal_albums": newSimpleAlbums(personalAlbums, users),
				"shared_albums":   newSimpleAlbums(sharedAlbums, users),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":        session.User.Username,
			"user_role":       session.User.Role.String(),
			"personal_albums": newSimpleAlbums(personalAlbums, users),
		})
	})
}
