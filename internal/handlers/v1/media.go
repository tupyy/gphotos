package v1

// import (
// 	"context"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sirupsen/logrus"
// 	"github.com/tupyy/gophoto/internal/api/utils"
// 	"github.com/tupyy/gophoto/internal/common"
// 	"github.com/tupyy/gophoto/internal/entity"
// 	"github.com/tupyy/gophoto/internal/services/album"
// 	"github.com/tupyy/gophoto/internal/services/media"
// 	"github.com/tupyy/gophoto/internal/services/permissions"
// 	"github.com/tupyy/gophoto/utils/logutil"
// )

// func GetAlbumMedia(r *gin.RouterGroup, albumService *album.Service, mediaService *media.Service) {
// 	r.GET("/api/albums/:id/album/media", utils.ParseAlbumIDHandler, func(c *gin.Context) {
// 		s, _ := c.Get("sessionData")
// 		session := s.(entity.Session)

// 		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
// 		logger := logutil.GetLogger(ctx)

// 		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
// 		if err != nil {
// 			common.AbortNotFound(c, err, "download media")

// 			return
// 		}

// 		// check permissions to this album
// 		ats := permissions.NewAlbumPermissionService()
// 		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
// 			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionReadAlbum}).
// 			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionReadAlbum}).
// 			Strategy(permissions.AtLeastOneStrategy).
// 			Resolve(album, session.User)

// 		if !hasPermission {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"user_id":  session.User.ID,
// 			}).Error("user has no permissions to read media")

// 			common.AbortForbiddenWithJson(c, errors.New("user has no permission to read media"), "")

// 			return
// 		}

// 		media, err := mediaService.ListBucket(ctx, album.Bucket)
// 		if err != nil {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"bucket":   album.Bucket,
// 			}).Error("failed to get media")

// 		}

// 		c.JSON(http.StatusOK, media)

// 		return

// 	})
// }

// func DownloadMedia(r *gin.RouterGroup, albumService *album.Service, mediaService *media.Service) {
// 	r.GET("/api/albums/:id/album/:media/media", utils.ParseAlbumIDHandler, parseMediaFilenameHandler, func(c *gin.Context) {
// 		s, _ := c.Get("sessionData")
// 		session := s.(entity.Session)

// 		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
// 		logger := logutil.GetLogger(ctx)

// 		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
// 		if err != nil {
// 			common.AbortNotFound(c, err, "download media")

// 			return
// 		}

// 		// check permissions to this album
// 		ats := permissions.NewAlbumPermissionService()
// 		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
// 			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionReadAlbum}).
// 			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionReadAlbum}).
// 			Strategy(permissions.AtLeastOneStrategy).
// 			Resolve(album, session.User)

// 		if !hasPermission {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"user_id":  session.User.ID,
// 			}).Error("user has no permissions to read media")

// 			common.AbortForbidden(c, errors.New("user has no permission to read media"), "")

// 			return
// 		}

// 		r, _, err := mediaService.GetPhoto(ctx, album.Bucket, c.GetString("media"))
// 		if err != nil {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"media":    c.GetString("media"),
// 			}).WithError(err).Error("failed to open media")

// 			common.AbortInternalError(c)

// 			return
// 		}

// 		fileContent, err := io.ReadAll(r)
// 		if err != nil {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"media":    c.GetString("media"),
// 			}).WithError(err).Error("failed to read from media")
// 			common.AbortInternalError(c)

// 			return
// 		}

// 		w := c.Writer
// 		w.Header().Set("Content-Type", "image/jpeg")
// 		w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))

// 		if _, err := w.Write(fileContent); err != nil {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"media":    c.GetString("media"),
// 			}).WithError(err).Error("failed to write media")
// 			common.AbortInternalError(c)

// 			return
// 		}

// 		return
// 	})
// }

// func DeleteMedia(r *gin.RouterGroup, albumService *album.Service, mediaService *media.Service) {
// 	r.DELETE("/api/albums/:id/album/:media/media", utils.ParseAlbumIDHandler, parseMediaFilenameHandler, func(c *gin.Context) {
// 		s, _ := c.Get("sessionData")
// 		session := s.(entity.Session)

// 		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
// 		logger := logutil.GetLogger(ctx)

// 		album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
// 		if err != nil {
// 			common.AbortNotFound(c, err, "delete media")

// 			return
// 		}

// 		// check permissions to this album
// 		ats := permissions.NewAlbumPermissionService()
// 		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
// 			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
// 			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
// 			Strategy(permissions.AtLeastOneStrategy).
// 			Resolve(album, session.User)

// 		if !hasPermission {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"user_id":  session.User.ID,
// 			}).Error("user has no permissions to delete media")

// 			common.AbortForbidden(c, errors.New("user has no permission to read media"), "")

// 			return
// 		}

// 		err = mediaService.Delete(ctx, album.Bucket, c.GetString("media"))
// 		if err != nil {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"media":    c.GetString("media"),
// 			}).WithError(err).Error("failed to remove media")

// 			common.AbortInternalError(c)

// 			return
// 		}

// 		return
// 	})
// }

// import (
// 	"errors"
// 	"html"
// 	"net/http"
// 	"regexp"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sirupsen/logrus"
// 	"github.com/tupyy/gophoto/internal/api/utils"
// 	"github.com/tupyy/gophoto/internal/common"
// 	"github.com/tupyy/gophoto/internal/conf"
// 	"github.com/tupyy/gophoto/internal/entity"
// 	"github.com/tupyy/gophoto/internal/services/album"
// 	"github.com/tupyy/gophoto/internal/services/media"
// 	"github.com/tupyy/gophoto/internal/services/permissions"
// 	"github.com/tupyy/gophoto/utils/encryption"
// 	"github.com/tupyy/gophoto/utils/logutil"
// )

// const (
// 	FILENAME_MAX_LENGTH = 100
// )

// var (
// 	filenameReg = regexp.MustCompile(`^[^±!@£$%&*+§¡€#¢§¶•ªº«\\/<>?:;|=,]*$`)
// )

// func UploadMedia(r *gin.RouterGroup, albumService *album.Service, mediaService *media.Service) {
// 	r.POST("/api/albums/:id/album/upload", utils.ParseAlbumIDHandler, func(c *gin.Context) {
// 		reqCtx := c.Request.Context()
// 		logger := logutil.GetLogger(c)

// 		s, _ := c.Get("sessionData")
// 		session := s.(entity.Session)

// 		album, err := albumService.Query().First(reqCtx, int32(c.GetInt("id")))
// 		if err != nil {
// 			common.AbortNotFound(c, err, "update album")

// 			return
// 		}

// 		// check permissions to this album
// 		ats := permissions.NewAlbumPermissionService()
// 		hasPermission := ats.Policy(permissions.OwnerPolicy{}).
// 			Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
// 			Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionWriteAlbum}).
// 			Strategy(permissions.AtLeastOneStrategy).
// 			Resolve(album, session.User)

// 		if !hasPermission {
// 			logger.WithFields(logrus.Fields{
// 				"album_id": album.ID,
// 				"user_id":  session.User.ID,
// 			}).Error("user has no permissions to upload media to this album")

// 			common.AbortForbidden(c, errors.New("user has no permission to upload media"), "")

// 			return
// 		}

// 		file, err := c.FormFile("file")
// 		if err != nil {
// 			logger.WithField("album id", album.ID).WithError(err).Error("failed to file from request")
// 			common.AbortInternalError(c)

// 			return
// 		}

// 		// validate filename
// 		if err := validate(file.Filename); err != nil {
// 			logger.WithField("filename", file.Filename).WithError(err).Error("failed to validate filename")
// 			common.AbortBadRequestWithJson(c, err, "invalid filename")

// 			return
// 		}

// 		src, err := file.Open()
// 		if err != nil {
// 			logger.WithField("album id", album.ID).WithError(err).Error("failed to open file from request")
// 			common.AbortInternalError(c)

// 			return
// 		}
// 		defer src.Close()

// 		sanitizedFilename := html.EscapeString(file.Filename)

// 		if err := mediaService.Save(reqCtx, album.Bucket, sanitizedFilename, src, media.Photo); err != nil {
// 			logger.WithError(err).Error("failed to upload media")

// 			common.AbortInternalErrorWithJson(c)

// 			return
// 		}
// 	})
// }

// func validate(filename string) error {
// 	if len(filename) > FILENAME_MAX_LENGTH {
// 		return errors.New("filename length exceeds max length")
// 	}

// 	if !filenameReg.MatchString(filename) {
// 		return errors.New("filename container forbidden characters")
// 	}

// 	return nil
// }

// func parseMediaFilenameHandler(c *gin.Context) {
// 	logger := logutil.GetLogger(c)

// 	// decrypt album id
// 	gen := encryption.NewGenerator(conf.GetEncryptionKey())

// 	decryptedMedia, err := gen.DecryptData(c.Param("media"))
// 	if err != nil {
// 		logger.WithError(err).Error("cannot decrypt media filename")
// 		c.AbortWithError(http.StatusNotFound, err) // explicit return not found here

// 		return
// 	}

// 	c.Set("media", decryptedMedia)
// }