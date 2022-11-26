package v1

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"go.uber.org/zap"
)

func (server *Server) GetAlbumThumbnail(c *gin.Context, albumId apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album_id", albumId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "album not found")
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", id, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "")
		return
	}

	// only users with editPermission set for this album or one of user's group with the same permission
	// can edit this album
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
		Policy(permissions.AnyUserPermissionPolicty{}).
		Policy(permissions.AnyGroupPermissionPolicy{}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		zap.S().Errorw("user has no read permissions on the album", "album_id", id, "user", session.User.Username)
		common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "get album")
		return
	}

	thumbnail, _, err := server.MediaService().GetPhoto(c, album.Bucket, album.Thumbnail)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", id, "thumbnail_filename", album.Thumbnail, "bucket", album.Bucket, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "thumbnail not found")
		return
	}

	content, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		zap.S().Errorw("failed to read thumbnail", "error", err, "album_id", id, "thumbnail_filename", album.Thumbnail, "bucket", album.Bucket, "user", session.User.Username)
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.Data(http.StatusOK, "image/jpeg", content)
}
