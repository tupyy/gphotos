package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/entity"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"go.uber.org/zap"
)

func (server *Server) CreateTag(c *gin.Context) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to create tags", "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "only admins or editors can add tags"))
		return
	}

	var tagForm apiv1.TagRequestPayload
	if err := c.BindJSON(&tagForm); err != nil {
		zap.S().Errorw("failed to bind to payload", "payload", tagForm, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusBadRequest, mappersv1.MapFromStatusf(http.StatusBadRequest, "failed to parse payload: %s", err))
		return
	}

	tag := entity.Tag{
		Name:   escapeField(tagForm.Name),
		Color:  escapeFieldPtr2(tagForm.Color),
		UserID: session.User.ID,
	}

	newTag, err := server.TagService().Create(c, tag)
	if err != nil {
		zap.S().Errorw("failed to create tag", "error", err, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusCreated, mappersv1.MapTagToModel(newTag))
}

func (server *Server) UpdateTag(c *gin.Context, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to update tag", "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "only admins or editors can add tags"))
		return
	}

	var tagForm apiv1.TagRequestPayload
	if err := c.ShouldBindJSON(&tagForm); err != nil {
		zap.S().Errorw("failed to bind to payload", "payload", tagForm, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusBadRequest, mappersv1.MapFromStatusf(http.StatusBadRequest, "failed to parse payload: %s", err))
		return
	}

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
	}

	// get the old tag
	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		zap.S().Errorw("failed to get tag", "error", err, "tag_id", tagID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
		return
	}

	tag.Name = escapeField(tagForm.Name)
	tag.Color = escapeFieldPtr2(tagForm.Color)

	if err := server.TagService().Update(c, tag); err != nil {
		zap.S().Errorw("failed to update tag", "error", err, "tag_id", tagID, "data", tagForm, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	tag, _ = server.TagService().GetByID(c, session.User.ID, tagID)

	c.JSON(http.StatusOK, mappersv1.MapTagToModel(tag))
}

func (server *Server) DeleteTag(c *gin.Context, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to delete tag", "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "only admins or editors can add tags"))
		return
	}

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
	}

	// check if a tag exists.
	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		zap.S().Errorw("failed to get tag", "error", err, "tag_id", tagID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
		return
	}

	if err := server.TagService().Delete(c, tag); err != nil {
		zap.S().Errorw("failed to delete tag", "error", err, "tag_id", tagID, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) GetTags(c *gin.Context, params apiv1.GetTagsParams) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to delete tag", "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "only admins or editors can add tags"))
		return
	}

	tags, err := server.TagService().Get(c, session.User.ID)
	if err != nil {
		zap.S().Errorw("failed to get tags", "error", err, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, mappersv1.MapTagsToList(tags))
}

func (server *Server) RemoveTagFromAlbum(c *gin.Context, albumId apiv1.AlbumId, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
		return
	}

	albumID, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album_id", albumId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumId))
		return
	}

	album, err := server.AlbumService().Query().First(c, albumID)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumId))
		return
	}

	// check permissions to this album
	ats := permissions.NewAlbumPermissionService()
	hasPermission := ats.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionEditAlbum}).
		Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionEditAlbum}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		zap.S().Errorw("user has no permission to edit album", "album_id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "access denied"))
		return
	}

	dissociate := false
	for _, tag := range album.Tags {
		if tag.ID == tagID {
			if err := server.TagService().Dissociate(c, tag, album.ID); err != nil {
				zap.S().Errorw("failed to dissociate tag from album", "error", err, "album_id", albumID, "tag_id", tag.ID, "user", session.User.Username)
				apiErr := mappersv1.MapFromError(err)
				c.AbortWithStatusJSON(apiErr.Code, apiErr)
				return
			}
			dissociate = true
			break
		}
	}

	if !dissociate {
		zap.S().Warnw("tag not associated with album", "album_id", albumID, "tag_id", tagID, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) SetTagToAlbum(c *gin.Context, albumId apiv1.AlbumId, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
	}

	albumID, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album_id", albumId, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumId))
		return
	}

	album, err := server.AlbumService().Query().First(c, albumID)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "album with id '%s' not found", albumId))
		return
	}

	// check permissions to this album
	ats := permissions.NewAlbumPermissionService()
	hasPermission := ats.Policy(permissions.OwnerPolicy{}).
		Policy(permissions.UserPermissionPolicy{Permission: entity.PermissionEditAlbum}).
		Policy(permissions.GroupPermissionPolicy{Permission: entity.PermissionEditAlbum}).
		Strategy(permissions.AtLeastOneStrategy).
		Resolve(album, session.User)

	if !hasPermission {
		zap.S().Errorw("user has no permission to edit album", "album_id", albumID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusForbidden, mappersv1.MapFromStatus(http.StatusForbidden, "access denied"))
		return
	}

	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		zap.S().Errorw("failed to get tag", "tag_id", tagID, "user", session.User.Username)
		c.AbortWithStatusJSON(http.StatusNotFound, mappersv1.MapFromStatusf(http.StatusNotFound, "tag with '%s' not found", tagId))
		return
	}

	if err := server.TagService().Associate(c, tag, album.ID); err != nil {
		zap.S().Errorw("failed to associate tag to album", "error", err, "album_id", albumID, "tag_id", tag.ID, "user", session.User.Username)
		apiErr := mappersv1.MapFromError(err)
		c.AbortWithStatusJSON(apiErr.Code, apiErr)
		return
	}

	album.Tags = append(album.Tags, tag)
	c.JSON(http.StatusCreated, mappersv1.MapAlbumToModel(album))
}
