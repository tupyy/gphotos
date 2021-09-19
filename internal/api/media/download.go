package media

import (
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/permissions"
	"github.com/tupyy/gophoto/utils/logutil"
)

func DownloadMedia(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	minioRepo := repos[domain.MinioRepoName].(domain.Store)
	bucketRepo := repos[domain.BucketRepoName].(domain.Bucket)

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

		// get the bucket of this album
		bucket, err := bucketRepo.Get(reqCtx, album.ID)
		if err != nil {
			logger.WithField("album id", album.ID).WithError(err).Error("failed to get bucket for album")
			common.AbortInternalError(c, err, "failed to get bucket for album")

			return
		}

		r, err := minioRepo.GetFile(reqCtx, bucket.Urn, c.GetString("media"))
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
			log.Println("unable to write image.")
		}

		return
	})
}
