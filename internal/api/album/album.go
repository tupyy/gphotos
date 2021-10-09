package album

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/dto"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

func GetAlbums(r *gin.RouterGroup, albumService *album.Service, usersService *users.Service) {
	r.GET("/api/albums", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(c)

		// fetch users from keycloak
		users, err := usersService.Query().AllUsers(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get users")
			common.AbortInternalErrorWithJson(c)

			return
		}

		reqParams := bindRequestParams(c)

		q := albumService.Query().
			OwnAlbums(reqParams.FetchPersonalAlbums).
			SharedAlbums(reqParams.FetchSharedAlbums)

		for _, f := range reqParams.Filters {
			q.Where(f)
		}

		// setup sort
		switch c.Query("sort") {
		case "name":
			q.Sort(album.SortByName, album.NormalOrder)
		case "location":
			q.Sort(album.SortByLocation, album.NormalOrder)
		case "date-normal":
			q.Sort(album.SortByDate, album.NormalOrder)
		default:
			q.Sort(album.SortByDate, album.ReverseOrder)
		}

		albums, err := q.All(ctx, session.User)
		if err != nil {
			logger.WithError(err).Error("failed to get albums")

			common.AbortInternalErrorWithJson(c)
		}

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
	Filters             []album.Predicate
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

	return reqParams

}

// GenerateAlbumFilters generates a list of filters from the query parameters.
func generateFilters(c *gin.Context) []album.Predicate {
	predicates := make([]album.Predicate, 0, 3)

	logger := logutil.GetLogger(c)

	if c.Query("start_date") != "" {
		if startDate, err := time.Parse("02/01/2006", c.Query("start_date")); err != nil {
			logger.WithError(err).Error("cannot parse start_date query param")
		} else {
			f := album.AfterDate(startDate)
			predicates = append(predicates, f)
		}
	}

	if c.Query("end_date") != "" {
		if endDate, err := time.Parse("02/01/2006", c.Query("end_date")); err != nil {
			logger.WithError(err).Error("cannot parse end_date query param")
		} else {
			f := album.BeforeDate(endDate)
			predicates = append(predicates, f)
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

		f := album.Owner(ownerIDs)
		predicates = append(predicates, f)
	}

	return predicates
}
