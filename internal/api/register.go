package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/api/album"
	"github.com/tupyy/gophoto/internal/api/index"
	"github.com/tupyy/gophoto/internal/api/media"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	mediaService "github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup,
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
	album.Thumbnail(privateGroup, as)
	logutil.GetDefaultLogger().Info("album handlers registered")
}
