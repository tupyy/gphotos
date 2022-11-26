package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	mappersv1 "github.com/tupyy/gophoto/internal/mappers/v1"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"go.uber.org/zap"
)

func (server *Server) CreateTag(c *gin.Context) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to create tags", "user", session.User.Username)
		common.AbortForbiddenWithJson(c, errors.New("only admins or editors can add tags"), "only admins or editors can add tags")
		return
	}

	var tagForm apiv1.TagRequestPayload
	if err := c.BindJSON(&tagForm); err != nil {
		zap.S().Errorw("failed to bind to payload", "payload", tagForm, "user", session.User.Username)
		common.AbortBadRequestWithJson(c, err, "failed to bind to form")
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
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.JSON(http.StatusCreated, mappersv1.MapTagToModel(newTag))
}

func (server *Server) UpdateTag(c *gin.Context, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to update tag", "user", session.User.Username)
		common.AbortForbiddenWithJson(c, errors.New("only admins or editors can update tags"), "only admins or editors can update tags")
		return
	}

	var tagForm apiv1.TagRequestPayload
	if err := c.ShouldBindJSON(&tagForm); err != nil {
		zap.S().Errorw("failed to bind to payload", "payload", tagForm, "user", session.User.Username)
		common.AbortBadRequestWithJson(c, err, "failed to bind to form")
		return
	}

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "")
	}

	// get the old tag
	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		zap.S().Errorw("failed to get tag", "error", err, "tag_id", tagID, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "")
		return
	}

	tag.Name = escapeField(tagForm.Name)
	tag.Color = escapeFieldPtr2(tagForm.Color)

	if err := server.TagService().Update(c, tag); err != nil {
		zap.S().Errorw("failed to update tag", "error", err, "tag_id", tagID, "data", tagForm, "user", session.User.Username)
		common.AbortInternalErrorWithJson(c)
		return
	}

	tag, _ = server.TagService().GetByID(c, session.User.ID, tagID)

	c.JSON(http.StatusOK, mappersv1.MapTagToModel(tag))
}

func (server *Server) DeleteTag(c *gin.Context, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		zap.S().Errorw("permissions denied to delete tag", "user", session.User.Username)
		common.AbortForbiddenWithJson(c, errors.New("only admins or editors can update tags"), "only admins or editors can update tags")
		return
	}

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "")
	}

	// check if a tag exists.
	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		zap.S().Errorw("failed to get tag", "error", err, "tag_id", tagID, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "")
		return
	}

	if err := server.TagService().Delete(c, tag); err != nil {
		zap.S().Errorw("failed to delete tag", "error", err, "tag_id", tagID, "user", session.User.Username)
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) GetTags(c *gin.Context, params apiv1.GetTagsParams) {
	session := c.MustGet("session").(entity.Session)

	tags, err := server.TagService().Get(c, session.User.ID)
	if err != nil {
		zap.S().Errorw("failed to get tags", "error", err, "user", session.User.Username)
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.JSON(http.StatusOK, mappersv1.MapTagsToList(tags))
}

func (server *Server) RemoveTagFromAlbum(c *gin.Context, albumId apiv1.AlbumId, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "")
		return
	}

	albumID, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album_id", albumId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "album not found")
		return
	}

	album, err := server.AlbumService().Query().First(c, albumID)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", albumID, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "dissociate tag from album")
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
		common.AbortForbiddenWithJson(c, errors.New("user has no permission to edit album"), "")
		return
	}

	dissociate := false
	for _, tag := range album.Tags {
		if tag.ID == tagID {
			if err := server.TagService().Dissociate(c, tag, album.ID); err != nil {
				zap.S().Errorw("failed to dissociate tag from album", "error", err, "album_id", albumID, "tag_id", tag.ID, "user", session.User.Username)
				common.AbortInternalErrorWithJson(c)
				return
			}
			dissociate = true
			break
		}
	}

	if !dissociate {
		zap.S().Warnw("tag not associated with album", "album_id", albumID, "tag_id", tagID, "user", session.User.Username)
		common.AbortBadRequestWithJson(c, errors.New("tag not associated with album"), "")
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) SetTagToAlbum(c *gin.Context, albumId apiv1.AlbumId, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	tagID, err := server.EncryptionService().Decrypt(tagId)
	if err != nil {
		zap.S().Errorw("failed to decrypt tag id", "error", err, "tag_id", tagId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "tag not found")
	}

	albumID, err := server.EncryptionService().Decrypt(albumId)
	if err != nil {
		zap.S().Errorw("failed to decrypt album id", "error", err, "album_id", albumId, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "album not found")
		return
	}

	album, err := server.AlbumService().Query().First(c, albumID)
	if err != nil {
		zap.S().Errorw("failed to get album", "error", err, "album_id", albumID, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, err, "associate tag from album")
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
		common.AbortForbiddenWithJson(c, errors.New("user has no permission to edit album"), "")
		return
	}

	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		zap.S().Errorw("failed to get tag", "tag_id", tagID, "user", session.User.Username)
		common.AbortNotFoundWithJson(c, errors.New("tag not found"), "tag not found")
		return
	}

	if err := server.TagService().Associate(c, tag, album.ID); err != nil {
		zap.S().Errorw("failed to associate tag to album", "error", err, "album_id", albumID, "tag_id", tag.ID, "user", session.User.Username)
		common.AbortInternalErrorWithJson(c)
		return
	}

	album.Tags = append(album.Tags, tag)
	c.JSON(http.StatusCreated, mappersv1.MapAlbumToModel(album))
}
