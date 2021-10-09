package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/handlers/album"
	"github.com/tupyy/gophoto/internal/handlers/index"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, as *albumService.Service, us *users.Service) {

	// register index controller
	index.Index(privateGroup, us)
	logutil.GetDefaultLogger().Info("index controller registered")

	// album handlers
	album.GetAlbum(privateGroup, as, us)
	album.GetCreateAlbumForm(privateGroup, us)
	album.GetUpdateAlbumForm(privateGroup, as, us)
	album.CreateAlbum(privateGroup, as)
	album.UpdateAlbum(privateGroup, as)
	album.DeleteAlbum(privateGroup, as)
	logutil.GetDefaultLogger().Info("album handlers registered")
}
