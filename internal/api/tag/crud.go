package tag

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/api/dto"
	"github.com/tupyy/gophoto/internal/api/utils"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/permissions"
	"github.com/tupyy/gophoto/internal/services/tag"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Get(r *gin.RouterGroup, tagService *tag.Service) {
	r.GET("/tags", func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		tags, err := tagService.Get(reqCtx, session.User.ID)
		if err != nil {
			logger.WithError(err).WithField("user id", session.User.ID).Error("fetch tags")

			common.AbortInternalErrorWithJson(c)

			return
		}

		dtos := make([]dto.Tag, 0, len(tags))
		for _, tag := range tags {
			dto, err := dto.NewTagDTO(tag)
			if err != nil {
				logger.WithError(err).WithField("tag", tag.String()).Warn("create tag dto")

				continue
			}

			dtos = append(dtos, dto)
		}

		c.HTML(http.StatusOK, "tags.html", gin.H{
			"tags": dtos,
		})
	})
}

func Dissociate(r *gin.RouterGroup, albumService *album.Service, tagService *tag.Service) {
	r.DELETE("/api/albums/:id/tag/:tagID/dissociate", utils.ParseAlbumIDHandler, func(c *gin.Context) {
		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		album, err := albumService.Query().First(reqCtx, int32(c.GetInt("id")))
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

		// decrypt tag id
		gen := encryption.NewGenerator(conf.GetEncryptionKey())

		stagID, err := gen.DecryptData(c.Param("tagID"))
		if err != nil {
			logger.WithError(err).Error("decrypt tag id")
		}

		tagID, err := strconv.Atoi(stagID)
		if err != nil {
			logger.WithError(err).WithField("tag id", tagID).Error("convert to int")
		}

		dissociate := false
		for _, tag := range album.Tags {
			if tag.ID == int32(tagID) {
				if err := tagService.Dissociate(reqCtx, tag, album.ID); err != nil {
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

		return
	})
}
