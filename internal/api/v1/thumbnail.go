package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/services/album"
)

func Thumbnail(r *gin.RouterGroup, albumService *album.Service) {
	// s, _ := c.Get("sessionData")
	// session := s.(entity.Session)

	// ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	// logger := logutil.GetLogger(ctx)

	// album, err := albumService.Query().First(ctx, int32(c.GetInt("id")))
	// if err != nil {
	// 	logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
	// 	common.AbortNotFoundWithJson(c, err, "update album")

	// 	return
	// }

	// // only editors and admins have the right to create albums
	// apr := permissions.NewAlbumPermissionService()
	// hasPermission := apr.Policy(permissions.OwnerPolicy{}).
	// 	Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
	// 	Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionEditAlbum}).
	// 	Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionEditAlbum}).
	// 	Strategy(permissions.AtLeastOneStrategy).
	// 	Resolve(album, session.User)

	// if !hasPermission {
	// 	common.AbortForbiddenWithJson(c, errors.New("user has no edit permission"), "user role forbids editing the album")

	// 	return
	// }

	// var thumbnailForm form.AlbumThumbnail
	// if err := c.ShouldBind(&thumbnailForm); err != nil {
	// 	logger.WithError(err).WithField("query parameters", fmt.Sprintf("%v", thumbnailForm)).Error("failed to bind query parameters to form")
	// 	common.AbortBadRequestWithJson(c, err, "fail to bind to form")

	// 	return
	// }

	// gen := encryption.NewGenerator(conf.GetEncryptionKey())

	// decryptedImageName, err := gen.DecryptData(thumbnailForm.Image)
	// if err != nil {
	// 	logger.WithError(err).WithField("album", album.String()).Error("failed to set thumbnail")

	// 	common.AbortInternalErrorWithJson(c)

	// 	return
	// }

	// parts := strings.Split(decryptedImageName, "/")
	// thumbnailImageName := fmt.Sprintf("thumbnail/%s", parts[1])

	// logger.WithField("thumbnail", thumbnailImageName).Debug("set thumbnail")

	// encryptedThumbnailImageName, _ := gen.EncryptData(thumbnailImageName)
	// album.Thumbnail = encryptedThumbnailImageName

	// if _, err := albumService.Update(ctx, album); err != nil {
	// 	logger.WithError(err).WithField("album", album.String()).Error("failed to set thumbnail")

	// 	common.AbortInternalErrorWithJson(c)

	// 	return
	// }

	// c.JSON(http.StatusOK, "ok")
}
