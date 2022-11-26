package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/permissions"
)

// (POST /api/gphotos/v1/albums/{album_id}/permissions)
func (server *Server) SetAlbumPermissions(c *gin.Context, albumId apiv1.AlbumId) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, errors.New("not found"), "not found")
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", id, "user", session.User.Username)
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
		zap.S().Errorw("permission denied to set permission for album", "album", id, "user", session.User.Username)
		common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "get album")
		return
	}

	var payload apiv1.AlbumPermissionsRequest
	if err := c.BindJSON(&payload); err != nil {
		zap.S().Errorw("failed to bind to form", "album", id, "payload", payload, "user", session.User.Username)
		common.AbortBadRequestWithJson(c, err, "failed to bind to form")
		return
	}

	perms, err := mappersv1.MapToEntityPermissions(payload)
	if err != nil {
		zap.S().Errorw("failed to map permissions", "album", id, "permissions", payload, "user", session.User.Username)
		common.AbortBadRequestWithJson(c, err, "cannot set permissions")
		return
	}

	if err := server.AlbumService().SetPermissions(c, album, perms); err != nil {
		zap.S().Errorw("failed to set permissions to album", "error", err, "album_id", id, "permissions", perms, "user", session.User.Username)
		common.AbortBadRequestWithJson(c, err, "cannot set permissions")
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// (GET /api/gphotos/v1/albums/{album_id}/permissions)
func (server *Server) GetAlbumPermissions(c *gin.Context, albumId string) {
	session := c.MustGet("session").(entity.Session)

	id, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album id", albumId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, errors.New("not found"), "not found")
		return
	}

	album, err := server.AlbumService().Query().First(c, id)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", id, "user", session.User.Username)
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
		zap.S().Errorw("permission denied to retrieve permissions", "album_id", id, "user", session.User.Username)
		common.AbortForbidden(c, common.NewMissingPermissionError(entity.PermissionEditAlbum, album, session.User), "get album")
		return
	}
	c.JSON(http.StatusOK, mappersv1.MapAlbumPermissions(album))
}

func (s *Server) RemoveAlbumPermissions(c *gin.Context, albumId apiv1.AlbumId) {}
