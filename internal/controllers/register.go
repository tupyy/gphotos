package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, repos domain.Repositories) {

	// register index controller
	Index(privateGroup, repos)
	logutil.GetDefaultLogger().Info("index controller registered")

	// album controllers
	GetAlbum(privateGroup, repos)
	GetCreateAlbumForm(privateGroup, repos)
	GetUpdateAlbumForm(privateGroup, repos)
	CreateAlbum(privateGroup, repos)
	UpdateAlbum(privateGroup, repos)
	DeleteAlbum(privateGroup, repos)
	logutil.GetDefaultLogger().Info("album controllers registered")
}
