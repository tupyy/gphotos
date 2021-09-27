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
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/permissions"
	"github.com/tupyy/gophoto/utils/logutil"
)

func GetAlbumMedia(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	minioRepo := repos[domain.MinioRepoName].(domain.Store)

	r.GET("/api/albums/:id/album/media", parseAlbumIDHandler, func(c *gin.Context) {
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		// add username to context for logging
		reqCtx := context.WithValue(c.Request.Context(), "username", session.User.Username)

		album, err := albumRepo.GetByID(reqCtx, int32(c.GetInt("id")))
		if err != nil {
			common.AbortNotFound(c, err, "download media")

			return
		}

		// check permissions to this album
		atr := permissions.NewAlbumPermissionResolver()
		hasPermission := atr.Policy(permissions.OwnerPolicy{}).
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

		medias, err := minioRepo.ListBucket(reqCtx, album.Bucket)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("failed to read bucket")

			common.AbortInternalErrorWithJson(c, errors.New("failed to read bucket"), "")

			return
		}

		c.JSON(http.StatusOK, medias)

		return

	})
}

func DownloadMedia(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	minioRepo := repos[domain.MinioRepoName].(domain.Store)

	r.GET("/api/albums/:id/album/:media/media", parseAlbumIDHandler, parseMediaFilenameHandler, func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		album, err := albumRepo.GetByID(reqCtx, int32(c.GetInt("id")))
		if err != nil {
			common.AbortNotFound(c, err, "download media")

			return
		}

		// check permissions to this album
		atr := permissions.NewAlbumPermissionResolver()
		hasPermission := atr.Policy(permissions.OwnerPolicy{}).
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

		r, err := minioRepo.GetFile(reqCtx, album.Bucket, c.GetString("media"))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"media":    c.GetString("media"),
			}).WithError(err).Error("failed to open media")
			common.AbortInternalError(c, err, "failed to open media")

			return
		}

		fileContent, err := io.ReadAll(r)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"media":    c.GetString("media"),
			}).WithError(err).Error("failed to read from media")
			common.AbortInternalError(c, err, "failed to read media")

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
			common.AbortInternalError(c, err, "failed to write media")

			return
		}

		return
	})
}
