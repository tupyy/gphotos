package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/api/album"
	"github.com/tupyy/gophoto/internal/api/media"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	mediaService "github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
)

func RegisterApi(privateGroup *gin.RouterGroup,
	publicGroup *gin.RouterGroup,
	as *albumService.Service,
	ms *mediaService.Service,
	us *users.Service) {

	// register album api
	album.GetAlbums(privateGroup, as, us)
	logutil.GetDefaultLogger().Info("api album registered")

	media.UploadMedia(privateGroup, as, ms)
	media.DownloadMedia(privateGroup, as, ms)
	media.GetAlbumMedia(privateGroup, as, ms)

	logutil.GetDefaultLogger().Info("api media registered")
}
