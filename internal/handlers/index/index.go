package index

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/handlers/common"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Index(r *gin.RouterGroup, repos domain.Repositories) {
	keycloakRepo := repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		filters, err := generateFilters(session.User)
		if err != nil {
			logutil.GetLogger(c).WithError(err).Error("create user filters")
			common.AbortInternalError(c, err, "")

			return
		}

		users, err := keycloakRepo.GetUsers(c.Request.Context(), nil, filters...)
		if err != nil {
			logutil.GetLogger(c).WithError(err).Error("fetch user filters")
			common.AbortInternalError(c, err, "")

			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"user_role": session.User.Role.String(),
			"can_share": session.User.CanShare,
			"users":     serialize(users),
		})
	})
}
