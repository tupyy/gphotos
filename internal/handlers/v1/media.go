package v1

import (
	"context"
	"errors"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

const (
	FILENAME_MAX_LENGTH = 100
)

var (
	filenameReg = regexp.MustCompile(`^[^±!@£$%&*+§¡€#¢§¶•ªº«\\/<>?:;|=,]*$`)
)

func (server *Server) GetAlbumPhotos(c *gin.Context, albumID string, params apiv1.GetAlbumPhotosParams) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := server.EncryptionService().Decrypt(albumID)
	if err != nil {
		logger.WithError(err).WithField("album id", albumID).Error("failed to decrypt album id")
		common.AbortInternalError(c)
		return
	}

	album, err := server.AlbumService().Query().First(ctx, id)
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

	media, err := server.MediaService().ListBucket(ctx, album.Bucket)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"album_id": album.ID,
			"bucket":   album.Bucket,
		}).Error("failed to get media")
	}

	total := len(media)
	page, size := 0, 0

	if params.Page != nil {
		page = int(*params.Page)
	}
	if params.Size != nil {
		size = int(*params.Size)
	}

	model := mappersv1.MapMediaListToModel(album, paginate(media, page, size))
	model.Total = total
	c.JSON(http.StatusOK, model)
}

func (server *Server) GetPhoto(c *gin.Context, albumId apiv1.AlbumId, photoId apiv1.PhotoId) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		logger.WithError(err).WithField("album id", albumId).Error("failed to decrypt album id")
		common.AbortInternalError(c)
		return
	}

	album, err := server.AlbumService().Query().First(ctx, id)
	if err != nil {
		common.AbortNotFoundWithJson(c, err, "download media")
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

	pID, err := server.EncryptionService().Decrypt(photoId)
	if err != nil {
		logger.WithError(err).WithField("photo id", photoId).Error("failed to decrypt photo id")
		common.AbortBadRequestWithJson(c, err, "bad request")
		return
	}

	r, _, err := server.MediaService().GetPhoto(ctx, album.Bucket, pID)
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
}

// (DELETE /api/gphotos/v1/album/{album_id}/photo/{photo_id})
func (server *Server) DeletePhoto(c *gin.Context, albumId apiv1.AlbumId, photoId apiv1.PhotoId) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		logger.WithError(err).WithField("album id", albumId).Error("failed to decrypt album id")
		common.AbortInternalError(c)
		return
	}

	album, err := server.AlbumService().Query().First(ctx, id)
	if err != nil {
		common.AbortNotFoundWithJson(c, err, "delete media")
		return
	}

	// check permissions to this album
	ats := permissions.NewAlbumPermissionService()
	hasPermission := ats.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
		Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		logger.WithFields(logrus.Fields{
			"album_id": album.ID,
			"user_id":  session.User.ID,
		}).Error("user has no permissions to delete media")

		common.AbortForbidden(c, errors.New("user has no permission to read media"), "")

		return
	}

	pID, err := server.EncryptionService().Decrypt(photoId)
	if err != nil {
		logger.WithError(err).WithField("photo id", photoId).Error("failed to decrypt photo id")
		common.AbortBadRequestWithJson(c, err, "bad request")
		return
	}

	err = server.MediaService().Delete(ctx, album.Bucket, pID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"album_id": album.ID,
			"media":    c.GetString("media"),
		}).WithError(err).Error("failed to remove media")

		common.AbortInternalError(c)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// (POST /api/gphotos/v1/albums/{album_id}/photos)
func (server *Server) UploadPhoto(c *gin.Context, albumId apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		logger.WithError(err).WithField("album id", albumId).Error("failed to decrypt album id")
		common.AbortNotFoundWithJson(c, err, "album not found")
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		common.AbortNotFoundWithJson(c, err, "update album")
		return
	}

	// check permissions to this album
	ats := permissions.NewAlbumPermissionService()
	hasPermission := ats.Policy(permissions.OwnerPolicy{}).
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

	if err := server.MediaService().Save(c, album.Bucket, sanitizedFilename, src, media.Photo); err != nil {
		logger.WithError(err).Error("failed to upload media")
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.JSON(http.StatusCreated, mappersv1.MapMediaToModel(album, entity.Media{
		MediaType: entity.Photo,
		Bucket:    album.Bucket,
		Filename:  sanitizedFilename,
	}))
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

func paginate(photos []entity.Media, pageNumber, size int) []entity.Media {
	// pagination
	var page []entity.Media

	if pageNumber <= 0 || size <= 0 {
		return photos
	}

	offset := (pageNumber - 1) * size
	limit := offset + size

	if limit >= len(photos) {
		limit = len(photos)
	}

	if offset >= len(photos) {
		return []entity.Media{}
	}

	page = append(page, photos[offset:limit]...)

	return page
}
