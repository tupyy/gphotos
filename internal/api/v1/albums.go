package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/filter"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

// (GET /api/gphotos/v1/albums)
func (server *Server) GetAlbums(c *gin.Context, params apiv1.GetAlbumsParams) {
	s, _ := c.Get("sessionData")
	session := s.(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(c)

	userService := server.GetUserService()
	albumService := server.GetAlbumService()

	// fetch users from keycloak
	users, err := userService.Query().All(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to get users")
		common.AbortInternalErrorWithJson(c)

		return
	}

	q := albumService.Query()
	if params.Personal != nil {
		q.OwnAlbums(*params.Personal)
	}
	if params.Shared != nil {
		q.SharedAlbums(*params.Shared)
	}
	if params.Filter != nil {
		filter, err := filter.New(*params.Filter)
		if err != nil {
			logger.WithError(err).Error("failed to create filter engine")
			common.AbortBadRequestWithJson(c, err, err.Error())

			return
		}

		q.Filter(filter)
	}
	// setup sort
	if params.Sort != nil {
		switch *params.Sort {
		case "name":
			q.Sort(album.SortByName, album.NormalOrder)
		case "location":
			q.Sort(album.SortByLocation, album.NormalOrder)
		case "date-normal":
			q.Sort(album.SortByDate, album.NormalOrder)
		default:
			q.Sort(album.SortByDate, album.ReverseOrder)
		}
	}

	// paginate
	if params.Offset != nil {
		q.Offset(int(*params.Offset))
	}

	if params.Limit != nil {
		q.Limit(int(*params.Limit))
	}

	albums, _, err := q.All(ctx, session.User)
	if err != nil {
		logger.WithError(err).Error("failed to get albums")
		common.AbortInternalErrorWithJson(c)
		return
	}

	albumModels := make([]apiv1.Album, 0, len(albums))
	for _, album := range albums {
		model, err := mapAlbumToModel(album, users)
		if err != nil {
			logutil.GetLogger(ctx).Error(err)
			continue
		}
		albumModels = append(albumModels, model)
	}

	c.JSON(http.StatusOK, &apiv1.Albums{
		Albums: &albumModels,
	})

	return
}

// (GET /api/gphotos/v1/albums/groups/{group_id})
func (server *Server) GetAlbumsByGroup(c *gin.Context, groupId string) {
}

// (GET /api/gphotos/v1/albums/users/{user_id})
func (server *Server) GetAlbumsByUser(c *gin.Context, userId string) {
}

// (GET /api/gphotos/v1/albums/{album_id}/permissions)
func (server *Server) GetAlbumPermissions(c *gin.Context, albumId string) {
}

// (GET /api/gphotos/v1/auth/callback)
func (server *Server) GetApiGphotosV1AuthCallback(c *gin.Context) {
}
