package v1

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/entity"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"go.uber.org/zap"
)

const (
	FILENAME_MAX_LENGTH = 100
)

var (
	filenameReg = regexp.MustCompile(`^[^±!@£$%&*+§¡€#¢§¶•ªº«\\/<>?:;|=,]*$`)
)

func (server *Server) GetAlbumPhotos(c *gin.Context, albumID string, params apiv1.GetAlbumPhotosParams) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumID)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
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
		zap.S().Errorw("permission denied to access album", "album", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, "access denied")
		return
	}

	photos, err := server.MediaService().ListBucket(c, album.Bucket)
	if err != nil {
		zap.S().Errorw("failed to get photos", "error", err, "album id", id, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	total := len(photos)
	page, size := 0, 0

	if params.Page != nil {
		page = int(*params.Page)
	}
	if params.Size != nil {
		size = int(*params.Size)
	}

	model := mappersv1.MapMediaListToModel(album, paginate(photos, page, size))
	model.Total = total
	c.JSON(http.StatusOK, model)
}

func (server *Server) GetPhoto(c *gin.Context, albumId apiv1.AlbumId, photoId apiv1.PhotoId) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
		return
	}

	album, err := server.AlbumService().Query().First(ctx, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "album", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
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
		c.AbortWithStatusJSON(http.StatusForbidden, "access denied")
		return
	}

	pID, err := server.EncryptionService().Decrypt(photoId)
	if err != nil {
		zap.S().Errorw("failed to decrypt photo id", "error", err, "photo id", photoId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("photo with id '%s' not found", id))
		return
	}

	r, _, err := server.MediaService().GetPhoto(ctx, album.Bucket, pID)
	if err != nil {
		zap.S().Errorw("failed to open photo", "error", err, "album id", id, "photo id", pID, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	fileContent, err := io.ReadAll(r)
	if err != nil {
		zap.S().Errorw("failed to read photo", "error", err, "album id", id, "photo id", pID, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	w := c.Writer
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))

	if _, err := w.Write(fileContent); err != nil {
		zap.S().Errorw("failed to write photo to response", "error", err, "album id", id, "photo id", pID, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
	}
}

// (DELETE /api/gphotos/v1/album/{album_id}/photo/{photo_id})
func (server *Server) DeletePhoto(c *gin.Context, albumId apiv1.AlbumId, photoId apiv1.PhotoId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
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
		zap.S().Errorw("user has no permission to write photo", "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, "access denied")
		return
	}

	pID, err := server.EncryptionService().Decrypt(photoId)
	if err != nil {
		zap.S().Errorw("failed to decrypt photo id", photoId, "album id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("photo with id '%s' not found", id))
		return
	}

	err = server.MediaService().Delete(c, album.Bucket, pID)
	if err != nil {
		zap.S().Errorw("failed to delete photo", "error", err, "photo id", pID, "album id", id, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// (POST /api/gphotos/v1/albums/{album_id}/photos)
func (server *Server) UploadPhoto(c *gin.Context, albumId apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album_id", albumId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("album with id '%s' not found", id))
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
		zap.S().Errorw("user hos no permission to update photo to album", "album_id", id, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, "access denied")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		zap.S().Errorw("failed to bind to file from request", "error", err, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to get file from request")
		return
	}

	// validate filename
	if err := validate(file.Filename); err != nil {
		zap.S().Errorw("failed to valdidate filename", "error", err, "filename", file.Filename, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("invalid filename '%s'", err))
		return
	}

	src, err := file.Open()
	if err != nil {
		zap.S().Errorw("failed to open file from request", "error", err, "filename", file.Filename, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusBadRequest, "failed to open file from request")
		return
	}
	defer src.Close()

	sanitizedFilename := html.EscapeString(file.Filename)

	if err := server.MediaService().Save(c, album.Bucket, sanitizedFilename, src, media.Photo); err != nil {
		zap.S().Errorw("failed to upload photo to repo", "error", err, "album_id", id, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
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
