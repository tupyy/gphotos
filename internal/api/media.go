package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain"
)

func UploadMedia(r *gin.RouterGroup, repos domain.Repositories) {
	r.POST("/api/albums/:id/album/upload", func(c *gin.Context) {

	})
}
