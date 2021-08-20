package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/utils/logutil"
)

func RegisterApi(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, repos domain.Repositories) {

	// register index controller
	GetAlbums(privateGroup, repos)
	logutil.GetDefaultLogger().Info("index controller registered")
}
