package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/handlers/album"
	"github.com/tupyy/gophoto/internal/handlers/index"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/keycloak"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, albumService *albumService.Service, keycloakService *keycloak.Service) {

	// register index controller
	index.Index(privateGroup, keycloakService)
	logutil.GetDefaultLogger().Info("index controller registered")

	// album handlers
	album.GetAlbum(privateGroup, albumService, keycloakService)
	album.GetCreateAlbumForm(privateGroup, keycloakService)
	album.GetUpdateAlbumForm(privateGroup, albumService, keycloakService)
	album.CreateAlbum(privateGroup, albumService)
	album.UpdateAlbum(privateGroup, albumService)
	album.DeleteAlbum(privateGroup, albumService)
	logutil.GetDefaultLogger().Info("album handlers registered")
}
