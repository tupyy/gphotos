package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, repos Repositories) {

	// register index controller
	Index(privateGroup, repos)
	logutil.GetDefaultLogger().Info("index controller registered")

	CreateAlbum(privateGroup, repos)
	logutil.GetDefaultLogger().Info("create album controller registered")

	UpdateAlbum(privateGroup, repos)
	logutil.GetDefaultLogger().Info("update album controller registered")
}
