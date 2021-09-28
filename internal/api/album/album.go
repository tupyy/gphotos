package album

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	albumFilter "github.com/tupyy/gophoto/internal/domain/filters/album"
	albumSort "github.com/tupyy/gophoto/internal/domain/sort/album"
	"github.com/tupyy/gophoto/internal/dto"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/keycloak"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

func GetAlbums(r *gin.RouterGroup, albumService album.Service, keycloakService keycloak.Service) {
	r.GET("/api/albums", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(c)

		// fetch users from keycloak
		users, err := keycloakService.GetUsers(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get users")
			common.AbortInternalErrorWithJson(c)

			return
		}

		// generate the req filters and sorter
		reqParams := bindRequestParams(c)

		albums, err := albumService.Query().
			OwnAlbums(reqParams.FetchPersonalAlbums).
			SharedAlbums(reqParams.FetchSharedAlbums).
			All(ctx, session.User)
		if err != nil {
			logger.WithError(err).Error("failed to get albums")

			common.AbortInternalErrorWithJson(c)
		}

		// sort all of them
		reqParams.Sorter.Sort(albums)

		c.JSON(http.StatusOK, gin.H{
			"user_role": session.User.Role.String(),
			"username":  fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"albums":    dto.NewAlbumDTOs(albums, users),
		})

		return
	})
}

type requestParams struct {
	FetchPersonalAlbums bool
	FetchSharedAlbums   bool
	Filters             albumFilter.Filters
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
func generateFilters(c *gin.Context) albumFilter.Filters {
	albumFilters := make(map[string]albumFilter.Filter)

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

			// hash the date
			key := string(base64.StdEncoding.EncodeToString([]byte(startDate.String())))
			albumFilters[key] = f
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

			key := string(base64.StdEncoding.EncodeToString([]byte(endDate.String())))
			albumFilters[key] = f
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

		keyBytes := []byte{}
		b := bytes.NewBuffer(keyBytes)
		for _, o := range ownerIDs {
			if _, err := b.WriteString(o); err != nil {
				logger.WithError(err).WithField("data", o).Error("failed to write owner id as bytes")
				continue
			}
		}
		key := string(base64.StdEncoding.EncodeToString(keyBytes))
		albumFilters[key] = f
	}

	return nil
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
