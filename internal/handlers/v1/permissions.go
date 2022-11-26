package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

// (POST /api/gphotos/v1/albums/{album_id}/permissions)
func (server *Server) SetAlbumPermissions(c *gin.Context, albumId apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		logger.WithError(err).WithField("album id", albumId).Error("failed to decrypt album id")
		common.AbortNotFoundWithJson(c, errors.New("not found"), "not found")
		return
	}

	album, err := server.AlbumService().Query().First(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
		common.AbortNotFoundWithJson(c, err, "update album")

		return
	}

	// only the admin or owner can visualize the permissions
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
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

	var payload apiv1.AlbumPermissionsRequest
	if err := c.BindJSON(&payload); err != nil {
		common.AbortBadRequestWithJson(c, err, "failed to bind to form")
		return
	}

	perms, err := mappersv1.MapToEntityPermissions(payload)
	if err != nil {
		common.AbortBadRequestWithJson(c, err, "cannot set permissions")
		return
	}

	if err := server.AlbumService().SetPermissions(ctx, album, perms); err != nil {
		common.AbortBadRequestWithJson(c, err, "cannot set permissions")
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// (GET /api/gphotos/v1/albums/{album_id}/permissions)
func (server *Server) GetAlbumPermissions(c *gin.Context, albumId string) {
	session := c.MustGet("session").(entity.Session)

	ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
	logger := logutil.GetLogger(ctx)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		logger.WithError(err).WithField("album id", albumId).Error("failed to decrypt album id")
		common.AbortNotFoundWithJson(c, errors.New("not found"), "not found")
		return
	}

	album, err := server.AlbumService().Query().First(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("album id", c.GetInt("id")).Error("failed to get album")
		common.AbortNotFound(c, err, "update album")

		return
	}

	// only the admin or owner can visualize the permissions
	apr := permissions.NewAlbumPermissionService()
	hasPermission := apr.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.RolePolicy{Role: entity.RoleAdmin}).
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
	c.JSON(http.StatusOK, mappersv1.MapAlbumPermissions(album))
}

func (s *Server) RemoveAlbumPermissions(c *gin.Context, albumId apiv1.AlbumId) {}
