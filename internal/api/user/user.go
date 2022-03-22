package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/api/dto"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
)

func GetAccount(r *gin.RouterGroup) {
	r.GET("/api/v1/account", func(c *gin.Context) {
		s, _ := c.Get("sessionData")
		session := s.(entity.Session)

		userDTO, err := dto.NewUserDTO(session.User)
		if err != nil {
			common.AbortInternalErrorWithJson(c)
		}

		c.JSON(http.StatusOK, userDTO)
	})
}
