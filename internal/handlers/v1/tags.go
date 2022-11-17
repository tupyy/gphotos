package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	presentersv1 "github.com/tupyy/gophoto/internal/presenters/v1"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/utils/encryption"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

func (server *Server) CreateTag(c *gin.Context) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		common.AbortForbiddenWithJson(c, errors.New("only admins or editors can add tags"), "only admins or editors can add tags")
		return
	}

	var tagForm apiv1.TagRequestPayload
	if err := c.BindJSON(&tagForm); err != nil {
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
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.JSON(http.StatusCreated, presentersv1.MapTagToModel(newTag))
}

func (server *Server) UpdateTag(c *gin.Context, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		common.AbortForbiddenWithJson(c, errors.New("only admins or editors can update tags"), "only admins or editors can update tags")

		return
	}

	var tagForm apiv1.TagRequestPayload
	if err := c.ShouldBindJSON(&tagForm); err != nil {
		common.AbortBadRequestWithJson(c, err, "failed to bind to form")

		return
	}

	logger := logutil.GetLogger(c)
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	tagID, err := gen.DecryptData(tagId)
	if err != nil {
		logger.WithError(err).Error("decrypt tag id")
	}

	// get the old tag
	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		common.AbortNotFoundWithJson(c, errors.New("tag does not exists"), "tag does not exists")

		return
	}

	tag.Name = escapeField(tagForm.Name)
	tag.Color = escapeFieldPtr2(tagForm.Color)

	if err := server.TagService().Update(c, tag); err != nil {
		common.AbortInternalErrorWithJson(c)

		return
	}

	tag, _ = server.TagService().GetByID(c, session.User.ID, tagID)

	c.JSON(http.StatusOK, presentersv1.MapTagToModel(tag))
}

func (server *Server) DeleteTag(c *gin.Context, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	if session.User.Role != entity.RoleAdmin && session.User.Role != entity.RoleEditor {
		common.AbortForbiddenWithJson(c, errors.New("only admins or editors can update tags"), "only admins or editors can update tags")

		return
	}

	logger := logutil.GetLogger(c)
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	tagID, err := gen.DecryptData(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("decrypt tag id")
	}

	// check if a tag exists.
	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		common.AbortBadRequestWithJson(c, errors.New("tag does not exists"), "tag does not exists")

		return
	}

	if err := server.TagService().Delete(c, tag); err != nil {
		logger.WithError(err).WithField("tag", tag.String()).Error("delete tag")
		common.AbortInternalErrorWithJson(c)

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (server *Server) GetTags(c *gin.Context, params apiv1.GetTagsParams) {
	session := c.MustGet("session").(entity.Session)

	logger := logutil.GetLogger(c)
	tags, err := server.TagService().Get(c, session.User.ID)
	if err != nil {
		logger.WithError(err).WithField("user id", session.User.ID).Error("fetch tags")
		common.AbortInternalErrorWithJson(c)
		return
	}

	c.JSON(http.StatusOK, presentersv1.MapTagsToList(tags))
}

func (server *Server) RemoveTagFromAlbum(c *gin.Context, albumId apiv1.AlbumId, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	logger := logutil.GetLogger(c)
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	tagID, err := gen.DecryptData(tagId)
	if err != nil {
		logger.WithError(err).Error("decrypt tag id")
	}

	albumID, _ := gen.DecryptData(albumId)

	album, err := server.GetAlbumService().Query().First(c, albumID)
	if err != nil {
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
		logger.WithFields(logrus.Fields{
			"album_id": album.ID,
			"user_id":  session.User.ID,
		}).Error("user has no permissions to edit album")

		common.AbortForbiddenWithJson(c, errors.New("user has no permission to edit album"), "")

		return
	}

	dissociate := false
	for _, tag := range album.Tags {
		if tag.ID == tagID {
			if err := server.TagService().Dissociate(c, tag, album.ID); err != nil {
				logger.WithFields(logrus.Fields{
					"tag":      tag.String(),
					"album_id": album.ID,
				}).WithError(err).Error("dissociate tag from album")
				common.AbortInternalErrorWithJson(c)
				return
			}

			dissociate = true

			break
		}
	}

	if !dissociate {
		logger.WithFields(logrus.Fields{
			"tag_id":   tagID,
			"album_id": album.ID,
		}).Error("tag not associated with album")

		common.AbortBadRequestWithJson(c, errors.New("tag not associated with album"), "")

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
func (server *Server) SetTagToAlbum(c *gin.Context, albumId apiv1.AlbumId, tagId apiv1.TagId) {
	session := c.MustGet("session").(entity.Session)

	logger := logutil.GetLogger(c)
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	tagID, err := gen.DecryptData(tagId)
	if err != nil {
		logger.WithError(err).Error("decrypt tag id")
	}

	albumID, _ := gen.DecryptData(albumId)
	album, err := server.GetAlbumService().Query().First(c, albumID)
	if err != nil {
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
		logger.WithFields(logrus.Fields{
			"album_id": album.ID,
			"user_id":  session.User.ID,
		}).Error("user has no permissions to edit album")

		common.AbortForbiddenWithJson(c, errors.New("user has no permission to edit album"), "")

		return
	}

	tag, err := server.TagService().GetByID(c, session.User.ID, tagID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"tag_id":  tagID,
			"user_id": session.User.ID,
		}).WithError(err).Error("user has no such tag")

		common.AbortBadRequestWithJson(c, errors.New("tag not found"), "tag not found")
	}

	if err := server.TagService().Associate(c, tag, album.ID); err != nil {
		logger.WithFields(logrus.Fields{
			"tag":      tag.String(),
			"album_id": album.ID,
		}).WithError(err).Error("associate tag with album")

		common.AbortInternalErrorWithJson(c)

		return
	}

	album.Tags = append(album.Tags, tag)
	c.JSON(http.StatusCreated, presentersv1.MapAlbumToModel(album))
}
