package media

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/image"
	"github.com/tupyy/gophoto/internal/permissions"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

const (
	FILENAME_MAX_LENGTH = 100
)

var (
	filenameReg = regexp.MustCompile(`^[^±!@£$%&*+§¡€#¢§¶•ªº«\\/<>?:;|=,]*$`)
)

func UploadMedia(r *gin.RouterGroup, repos domain.Repositories) {
	albumRepo := repos[domain.AlbumRepoName].(domain.Album)
	minioRepo := repos[domain.MinioRepoName].(domain.Store)

	r.POST("/api/albums/:id/album/upload", parseAlbumIDHandler, func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		album, err := albumRepo.GetByID(reqCtx, int32(c.GetInt("id")))
		if err != nil {
			common.AbortNotFound(c, err, "update album")

			return
		}

		// check permissions to this album
		atr := permissions.NewAlbumPermissionResolver()
		hasPermission := atr.Policy(permissions.OwnerPolicy{}).
			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
			Strategy(permissions.AtLeastOneStrategy).
			Resolve(album, session.User)

		if !hasPermission {
			logger.WithFields(logrus.Fields{
				"album_id": album.ID,
				"user_id":  session.User.ID,
			}).Error("user has no permissions to upload media to this album")

			common.AbortForbidden(c, errors.New("user has no permission to upload media"), "")

			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			logger.WithField("album id", album.ID).WithError(err).Error("failed to file from request")
			common.AbortInternalError(c)

			return
		}

		// validate filename
		if err := validate(file.Filename); err != nil {
			logger.WithField("filename", file.Filename).WithError(err).Error("failed to validate filename")
			common.AbortBadRequestWithJson(c, err, "invalid filename")

			return
		}

		src, err := file.Open()
		if err != nil {
			logger.WithField("album id", album.ID).WithError(err).Error("failed to open file from request")
			common.AbortInternalError(c)

			return
		}
		defer src.Close()

		sanitizedFilename := html.EscapeString(file.Filename)

		err = minioRepo.PutFile(reqCtx, conf.GetMinioTemporaryBucket(), sanitizedFilename, file.Size, src)
		if err != nil {
			logger.WithField("filename", file.Filename).WithError(err).Error("failed to put file into the bucket")
			common.AbortInternalError(c)

			return
		}

		// do image processing
		var imgBuffer bytes.Buffer
		var imgThumbnailBuffer bytes.Buffer
		if err := image.Process(src, &imgBuffer, &imgThumbnailBuffer); err != nil {
			logger.WithError(err).Error("failed to process image")
			common.AbortInternalError(c)

			return
		}

		logger.Info("image processing done")
		// save images

		basename := strings.Split(sanitizedFilename, ".")[0]

		if err := minioRepo.PutFile(reqCtx, album.Bucket, fmt.Sprintf("%s.jpg", basename), int64(imgBuffer.Len()), &imgBuffer); err != nil {
			logger.WithError(err).Errorf("failed to write image to bucket %s", album.Bucket)
			common.AbortInternalError(c)

			return
		}

		if err := minioRepo.PutFile(reqCtx, album.Bucket, fmt.Sprintf("%s_thumbnail.jpg", basename), int64(imgThumbnailBuffer.Len()), &imgThumbnailBuffer); err != nil {
			logger.WithError(err).Errorf("failed to write thumbnail to bucket %s", album.Bucket)
			common.AbortInternalError(c)

			return
		}

	})
}

func validate(filename string) error {
	if len(filename) > FILENAME_MAX_LENGTH {
		return errors.New("filename length exceeds max length")
	}

	if !filenameReg.MatchString(filename) {
		return errors.New("filename container forbidden characters")
	}

	return nil
}

// parseAlbumIDHandler decrypt the album id passes as parameters and set the id in the context.
func parseAlbumIDHandler(c *gin.Context) {
	logger := logutil.GetLogger(c)

	// decrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	decryptedID, err := gen.DecryptData(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("cannot decrypt album id")
		c.AbortWithError(http.StatusNotFound, err) // explicit return not found here

		return
	}

	id, err := strconv.Atoi(decryptedID)
	if err != nil {
		logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
		c.AbortWithError(http.StatusNotFound, err)

		return
	}

	c.Set("id", id)
}

func parseMediaFilenameHandler(c *gin.Context) {
	logger := logutil.GetLogger(c)

	// decrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	decryptedMedia, err := gen.DecryptData(c.Param("media"))
	if err != nil {
		logger.WithError(err).Error("cannot decrypt media filename")
		c.AbortWithError(http.StatusNotFound, err) // explicit return not found here

		return
	}

	c.Set("media", decryptedMedia)
}
