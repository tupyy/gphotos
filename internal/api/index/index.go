package index

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	goI18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tupyy/gophoto/i18n"
	"github.com/tupyy/gophoto/internal/api/dto"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, usersService *users.Service) {
	r.GET("/", func(c *gin.Context) {
		ss := sessions.Default(c)

		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		// return all users if current user is an admin
		if session.User.Role == entity.RoleAdmin {
			users, err := usersService.Query().
				Where(users.NotUsername(session.User.Username)).
				Where(users.CanShare(true)).
				Where(users.Roles([]entity.Role{entity.RoleEditor, entity.RoleUser})).
				All(ctx)
			if err != nil {
				logger.WithError(err).Error("failed to get users")
				common.AbortInternalError(c)

				return
			}

			c.HTML(http.StatusOK, "index.html", gin.H{
				"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
				"user_role": session.User.Role.String(),
				"can_share": session.User.CanShare,
				"users":     dto.NewUserDTOs(users),
			})

			return
		}

		var relatedUsers []entity.User
		// if current user can share get all users that share an album with the current one.
		if session.User.CanShare {
			var err error
			relatedUsers, err = usersService.Query().AllRelatedUsers(ctx, session.User)
			if err != nil {
				logger.WithError(err).Error("failed to get related users")
				common.AbortInternalError(c)

				return
			}
		}

		if len(session.Alerts) > 0 {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
				"user_role": session.User.Role.String(),
				"can_share": session.User.CanShare,
				"users":     dto.NewUserDTOs(relatedUsers),
				"alerts":    session.Alerts,
			})

			session.ClearAlerts()
			ss.Set(session.SessionID, session)
			ss.Save()

			return
		}

		// localizer
		accept := c.GetHeader("Accept-Language")
		localizer := goI18n.NewLocalizer(i18n.Bundle, accept)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":           fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"user_role":      session.User.Role.String(),
			"can_share":      session.User.CanShare,
			"users":          dto.NewUserDTOs(relatedUsers),
			"Title":          i18n.GetTranslation(localizer, "Title"),
			"CreateAlbum":    i18n.GetTranslation(localizer, "NavBarCreateAlbum"),
			"Tags":           i18n.GetTranslation(localizer, "NavBarTags"),
			"Filters":        i18n.GetTranslation(localizer, "IndexFilters"),
			"AlbumType":      i18n.GetTranslation(localizer, "IndexAlbumType"),
			"Personal":       i18n.GetTranslation(localizer, "IndexAlbumTypePersonal"),
			"Shared":         i18n.GetTranslation(localizer, "IndexAlbumTypeShared"),
			"Date":           i18n.GetTranslation(localizer, "IndexFilterDate"),
			"To":             i18n.GetTranslation(localizer, "IndexFilterTo"),
			"SharedByOthers": i18n.GetTranslation(localizer, "IndexFilterSharedByOthers"),
			"Albums":         i18n.GetTranslation(localizer, "Albums"),
			"SortBy":         i18n.GetTranslation(localizer, "IndexSort"),
			"Name":           i18n.GetTranslation(localizer, "IndexSortByName"),
			"DateNormal":     i18n.GetTranslation(localizer, "IndexSortByDate"),
			"DateRev":        i18n.GetTranslation(localizer, "IndexSortByDateRev"),
			"Location":       i18n.GetTranslation(localizer, "IndexSortByLocation"),
		})
	})
}
