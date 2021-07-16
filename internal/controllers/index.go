package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
)

func Index(r *gin.RouterGroup, albumRepo repo.AlbumRepo) {
	r.GET("/", func(c *gin.Context) {
		s, _ := c.Get("sessionData")

		session := s.(entity.Session)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"username": session.User.Username,
		})

	})
}
