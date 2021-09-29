package index

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/dto"
	"github.com/tupyy/gophoto/internal/services/keycloak"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, keycloakService *keycloak.Service) {
	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		ctx := context.WithValue(c.Request.Context(), "username", session.User.Username)
		logger := logutil.GetLogger(ctx)

		users, err := keycloakService.Query().
			Where(keycloak.NotUsername(session.User.Username)).
			Where(keycloak.CanShare(true)).
			Where(keycloak.Roles([]entity.Role{entity.RoleEditor, entity.RoleUser})).
			AllUsers(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get users")
			common.AbortInternalError(c)

			return
		}

		// return all users if current user is an admin
		if session.User.Role == entity.RoleAdmin {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
				"user_role": session.User.Role.String(),
				"can_share": session.User.CanShare,
				"users":     dto.NewUserDTOs(users),
			})

			return
		}

		// if current user can share get all users that share an album with the current one.
		// if session.User.CanShare {
		// 	// get all shared albums in order to filtered users which don't share albums with the current user
		// 	ids, err := userRepo.GetRelatedUsers(reqCtx, session.User)
		// 	if err != nil {
		// 		logger.WithError(err).WithField("user id", session.User.ID).Error("fetch related users")
		// 		common.AbortInternalError(c, err, "")

		// 		return
		// 	}

		// 	users = mapUsers(users, ids)
		// }

		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"user_role": session.User.Role.String(),
			"can_share": session.User.CanShare,
			"users":     dto.NewUserDTOs(users),
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
