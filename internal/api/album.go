package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	albumFilter "github.com/tupyy/gophoto/internal/domain/filters/album"
	albumSort "github.com/tupyy/gophoto/internal/domain/sort/album"
	"github.com/tupyy/gophoto/utils/encryption"
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
		users, err := keycloakRepo.GetUsers(reqCtx)
		if err != nil {
			logger.WithError(err).Error("index fetch users")
			c.AbortWithError(http.StatusInternalServerError, err)

			return
		}

		// generate the req filters and sorter
		reqParams := bindRequestParams(c)

		var personalAlbums []entity.Album
		if reqParams.FetchPersonalAlbums {
			var err error
			personalAlbums, err = albumRepo.GetByOwnerID(reqCtx, session.User.ID, reqParams.Filters...)
			if err != nil {
				logger.WithError(err).Error("error fetching personal albums")
				AbortWithJson(c, http.StatusInternalServerError, err, "")

				return
			}
		}

		// user with canShare=true can share albums with other users
		// fetch all albums for which the user has at least one permissions
		var sharedAlbums []entity.Album
		if reqParams.FetchSharedAlbums {
			notOwnerFilter, err := albumFilter.GenerateFilterFuncs(albumFilter.NotFilterByOwnerID, []string{session.User.ID})
			if err != nil {
				logger.WithError(err).Error("error generate notOwnerFilter")
				AbortWithJson(c, http.StatusInternalServerError, err, "")

				return
			}

			// if i'm an admin fetch all other albums regardless of permissions
			if session.User.Role == entity.RoleAdmin {
				filters := append(reqParams.Filters, notOwnerFilter)

				var err error
				sharedAlbums, err = albumRepo.Get(reqCtx, filters...)
				if err != nil {
					logger.WithError(err).Error("error fetching all albums")
					AbortWithJson(c, http.StatusInternalServerError, err, "")

					return
				}
			} else if session.User.CanShare {
				var err error
				sharedAlbums, err = albumRepo.GetByUserID(reqCtx, session.User.ID, reqParams.Filters...)
				if err != nil {
					logger.WithError(err).Error("error fetching shared albums")
					AbortWithJson(c, http.StatusInternalServerError, err, "")

					return
				}

				filters := append(reqParams.Filters, notOwnerFilter)
				// get albums shared by the user's groups
				groupSharedAlbum, err := albumRepo.GetByGroups(reqCtx, groupsToList(session.User.Groups), filters...)
				if err != nil {
					logger.
						WithError(err).
						WithFields(logrus.Fields{
							"user_id": session.User.ID,
							"groups":  session.User.Groups,
						}).Error("cannot fetch albums by group name")
					AbortWithJson(c, http.StatusInternalServerError, err, "")

					return
				}
				// join and remove the duplicates
				sharedAlbums = join(sharedAlbums, groupSharedAlbum)
			}
		}

		albums := merge(personalAlbums, sharedAlbums)

		// sort all of them
		reqParams.Sorter.Sort(albums)

		c.JSON(http.StatusOK, gin.H{
			"user_role": session.User.Role.String(),
			"username":  fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"albums":    newSimpleAlbums(albums, users),
		})

		return
	})
}

type requestParams struct {
	FetchPersonalAlbums bool
	FetchSharedAlbums   bool
	Filters             []albumFilter.Filter
	Sorter              albumSort.Sorter
}

// bindRequestParams returns a struct with filters and a sorter generated from query parameters
func bindRequestParams(c *gin.Context) requestParams {
	logger := logutil.GetLogger(c)

	reqParams := requestParams{
		FetchPersonalAlbums: true,
		FetchSharedAlbums:   true,
	}

	if c.Query("personal") != "" {
		personalAlbumsFilterValue, err := strconv.ParseBool(c.Query("personal"))
		if err != nil {
			logger.WithError(err).WithField("personal", c.Query("personal")).Warn("cannot parse personal filter value")
		} else {
			reqParams.FetchPersonalAlbums = personalAlbumsFilterValue
		}
	}

	if c.Query("shared") != "" {
		sharedAlbumsFilterValue, err := strconv.ParseBool(c.Query("shared"))
		if err != nil {
			logger.WithError(err).WithField("shared", c.Query("shared")).Warn("cannot parse shared filter value")
		} else {
			reqParams.FetchSharedAlbums = sharedAlbumsFilterValue
		}
	}

	reqParams.Filters = generateFilters(c)
	reqParams.Sorter = generateSort(c)

	return reqParams

}

// GenerateAlbumFilters generates a list of filters from the query parameters.
func generateFilters(c *gin.Context) []albumFilter.Filter {
	albumFilters := make([]albumFilter.Filter, 0, 5)

	logger := logutil.GetLogger(c)

	if c.Query("start_date") != "" {
		if startDate, err := time.Parse("02/01/2006", c.Query("start_date")); err != nil {
			logger.WithError(err).Error("cannot parse start_date query param")
		} else {
			f, err := albumFilter.GenerateFilterFuncs(albumFilter.FilterAfterDate, startDate)
			if err != nil {
				logger.WithError(err).Error("error create FilterAfterDate filter")
			}

			logger.WithField("start date", startDate).Debug("filter start date created")
			albumFilters = append(albumFilters, f)
		}
	}

	if c.Query("end_date") != "" {
		if endDate, err := time.Parse("02/01/2006", c.Query("end_date")); err != nil {
			logger.WithError(err).Error("cannot parse end_date query param")
		} else {
			f, err := albumFilter.GenerateFilterFuncs(albumFilter.FilterBeforeDate, endDate)
			if err != nil {
				logger.WithError(err).Error("error create FilterBeforeDate filter")
			}

			logger.WithField("end date", endDate).Debug("filter end date created")
			albumFilters = append(albumFilters, f)
		}
	}

	owners := c.QueryArray("owner")
	if len(owners) > 0 {
		gen := encryption.NewGenerator(conf.GetEncryptionKey())

		ownerIDs := make([]string, 0, len(owners))
		for _, o := range owners {
			ownerID, err := gen.DecryptData(o)
			if err != nil {
				logger.WithError(err).WithField("data", o).Error("error decrypt owner id")

				continue
			}

			ownerIDs = append(ownerIDs, ownerID)
			logger.WithField("owner_id", ownerID).Debug("filter by owner id created")
		}

		f, err := albumFilter.GenerateFilterFuncs(albumFilter.FilterByOwnerID, ownerIDs)
		if err != nil {
			logger.WithError(err).Error("error create FilterOwnerID filter")
		}

		albumFilters = append(albumFilters, f)
	}

	return albumFilters
}

func generateSort(c *gin.Context) albumSort.Sorter {
	switch c.Query("sort") {
	case "name":
		return albumSort.NewSorter(albumSort.SortByName, albumSort.NormalOrder)
	case "location":
		return albumSort.NewSorter(albumSort.SortByLocation, albumSort.NormalOrder)
	case "date-normal":
		return albumSort.NewSorter(albumSort.SortByDate, albumSort.NormalOrder)
	default:
		return albumSort.NewSorter(albumSort.SortByDate, albumSort.ReverseOrder)
	}
}
