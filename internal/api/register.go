package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/api/album"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/keycloak"
	"github.com/tupyy/gophoto/utils/logutil"
)

func RegisterApi(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, albumService *albumService.Service, keycloakService *keycloak.Service) {

	// register album api
	album.GetAlbums(privateGroup, albumService, keycloakService)
	logutil.GetDefaultLogger().Info("api album registered")

	// register media upload api
	// media.UploadMedia(privateGroup, albumService)
	// logutil.GetDefaultLogger().Info("api media registered")

	// media.DownloadMedia(privateGroup, al
}
