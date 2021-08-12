package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain/entity"
)

func Index(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":      fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName),
			"user_role": session.User.Role.String(),
		})
	})
}
