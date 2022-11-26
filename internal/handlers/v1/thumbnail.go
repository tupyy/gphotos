package v1

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

func (server *Server) GetAlbumThumbnail(c *gin.Context, albumId apiv1.AlbumId) {
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
		logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
		common.AbortNotFound(c, err, "update album")

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
		logger.WithFields(logrus.Fields{
			"request user id": session.User.ID,
			"album owner id":  album.Owner,
		}).Error("current user has no permission of this album")
		common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "get album")
		return
	}

	thumbnail, _, err := server.MediaService().GetPhoto(c, album.Bucket, album.Thumbnail)
	if err != nil {
		common.AbortNotFoundWithJson(c, err, "thumbnail not found")
		return
	}

	content, err := ioutil.ReadAll(thumbnail)
	if err != nil {
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.Data(http.StatusOK, "image/jpeg", content)
}
