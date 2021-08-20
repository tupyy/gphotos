package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/handlers/album"
	"github.com/tupyy/gophoto/internal/handlers/index"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, repos domain.Repositories) {

	// register index controller
	index.Index(privateGroup, repos)
	logutil.GetDefaultLogger().Info("index controller registered")

	// album handlers
	album.GetAlbum(privateGroup, repos)
	album.GetCreateAlbumForm(privateGroup, repos)
	album.GetUpdateAlbumForm(privateGroup, repos)
	album.CreateAlbum(privateGroup, repos)
	album.UpdateAlbum(privateGroup, repos)
	album.DeleteAlbum(privateGroup, repos)
	logutil.GetDefaultLogger().Info("album handlers registered")
}
