package album

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/api/dto"
	"github.com/tupyy/gophoto/internal/api/utils"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/services/tag"
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
		users, err := usersService.Query().All(ctx)
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

		var offset int
		if len(c.Query("offset")) > 0 {
			if o, err := strconv.Atoi(c.Query("offset")); err == nil {
				offset = o
			}
		}

		var limit int
		if len(c.Query("limit")) > 0 {
			if l, err := strconv.Atoi(c.Query("limit")); err == nil {
				limit = l
			}
		}

		albums, count, err := q.Offset(offset).Limit(limit).All(ctx, session.User)
		if err != nil {
			logger.WithError(err).Error("failed to get albums")

			common.AbortInternalErrorWithJson(c)
		}

		c.JSON(http.StatusOK, gin.H{
			"user_role": session.User.Role.String(),
			"username":  fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"albums":    dto.NewAlbumDTOs(albums, users),
			"count":     count,
		})

		return
	})
}

func GetAlbumsTags(r *gin.RouterGroup, albumService *album.Service, tagService *tag.Service) {
	r.GET("/api/albums/:id/tags", utils.ParseAlbumIDHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
		if err != nil {
			logger.WithError(err).WithField("id", c.GetInt("id")).Error("album not found")
			common.AbortNotFound(c, err, "failed to album")

			return
		}

		// check permissions to this album
		ats := permissions.NewAlbumPermissionService()
		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
			Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
			Policy(permissions.AnyUserPermissionPolicty{}).
			Policy(permissions.AnyGroupPermissionPolicy{}).
			Strategy(permissions.AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("user has no permissions to access this album")

			common.AbortForbiddenWithJson(c, common.NewMissingPermissionError(entity.PermissionReadAlbum, album, session.User), "user has no permission for the album")

			return
		}

		tags, err := tagService.GetByAlbum(ctx, album.ID)
		if err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("fetch tags for album")

			common.AbortInternalErrorWithJson(c)

			return
		}

		dtos := make([]dto.Tag, 0, len(tags))
		for _, tag := range tags {
			dto, err := dto.NewTagDTO(tag)
			if err != nil {
				logger.WithError(err).WithField("tag", tag.String()).Warn("create tag dto")

				continue
			}

			dtos = append(dtos, dto)
		}

		c.JSON(http.StatusOK, gin.H{
			"tags": dtos,
		})
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
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

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

	tagParameter := c.QueryArray("tag")
	if len(tagParameter) > 0 {
		tags := make([]string, 0, len(tagParameter))

		for _, tag := range tagParameter {
			tags = append(tags, strings.ToLower(tag))
		}

		f := album.Tags(tags)
		predicates = append(predicates, f)
	}

	return predicates
}
