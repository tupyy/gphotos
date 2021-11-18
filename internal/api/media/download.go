package media

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/utils/logutil"
)

func GetAlbumMedia(r *gin.RouterGroup, albumService *album.Service, mediaService *media.Service) {
	r.GET("/api/albums/:id/album/media", parseAlbumIDHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
		if err != nil {
			common.AbortNotFound(c, err, "download media")

			return
		}

		// check permissions to this album
		ats := permissions.NewAlbumPermissionService()
		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionReadAlbum}).
			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionReadAlbum}).
			Strategy(permissions.AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("user has no permissions to read media")

			common.AbortForbiddenWithJson(c, errors.New("user has no permission to read media"), "")

			return
		}

		media, err := mediaService.ListBucket(ctx, album.Bucket)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"bucket":   album.Bucket,
			}).Error("failed to get media")

		}

		c.JSON(http.StatusOK, media)

		return

	})
}

func DownloadMedia(r *gin.RouterGroup, albumService *album.Service, mediaService *media.Service) {
	r.GET("/api/albums/:id/album/:media/media", parseAlbumIDHandler, parseMediaFilenameHandler, func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
		if err != nil {
			common.AbortNotFound(c, err, "download media")

			return
		}

		// check permissions to this album
		ats := permissions.NewAlbumPermissionService()
		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionReadAlbum}).
			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionReadAlbum}).
			Strategy(permissions.AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("user has no permissions to read media")

			common.AbortForbidden(c, errors.New("user has no permission to read media"), "")

			return
		}

		r, err := mediaService.GetPhoto(ctx, album.Bucket, c.GetString("media"))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"media":    c.GetString("media"),
			}).WithError(err).Error("failed to open media")

			common.AbortInternalError(c)

			return
		}

		fileContent, err := io.ReadAll(r)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"media":    c.GetString("media"),
			}).WithError(err).Error("failed to read from media")
			common.AbortInternalError(c)

			return
		}

		w := c.Writer
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))

		if _, err := w.Write(fileContent); err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"media":    c.GetString("media"),
			}).WithError(err).Error("failed to write media")
			common.AbortInternalError(c)

			return
		}

		return
	})
}
