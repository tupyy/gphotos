package index

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/handlers/common"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, repos domain.Repositories) {
	keycloakRepo := repos[domain.KeycloakRepoName].(domain.KeycloakRepo)
	userRepo := repos[domain.UserRepoName].(domain.User)

	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		reqCtx := c.Request.Context()
		logger := logutil.GetLogger(c)

		filters, err := generateFilters(session.User)
		if err != nil {
			logger.WithError(err).Error("create user filters")
			common.AbortInternalError(c, err, "")

			return
		}

		users, err := keycloakRepo.GetUsers(reqCtx, nil, filters...)
		if err != nil {
			logger.WithError(err).Error("fetch user filters")
			common.AbortInternalError(c, err, "")

			return
		}

		if session.User.Role != entity.RoleAdmin {
			// get all shared albums in order to filtered users which don't share albums with the current user
			ids, err := userRepo.GetRelatedUsers(reqCtx, session.User)
			if err != nil {
				logger.WithError(err).WithField("user id", session.User.ID).Error("fetch related users")
				common.AbortInternalError(c, err, "")

				return
			}

			users = mapUsers(users, ids)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"user_role": session.User.Role.String(),
			"can_share": session.User.CanShare,
			"users":     serialize(users),
		})
	})
}

func mapUsers(users []entity.User, ids []string) []entity.User {
	relatedUsers := make([]entity.User, 0, len(ids))

	// remove users which are not relevant for albums found.
	addedUsers := make(map[string]interface{})
	for _, id := range ids {
		for _, u := range users {
			_, alreadyAdded := addedUsers[u.ID]

			if u.ID == id && !alreadyAdded {
				relatedUsers = append(relatedUsers, u)
				addedUsers[u.ID] = true
			}
		}
	}

	return relatedUsers
}
