package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/api/album"
	"github.com/tupyy/gophoto/internal/api/index"
	"github.com/tupyy/gophoto/internal/api/media"
	"github.com/tupyy/gophoto/internal/api/tag"
	albumService "github.com/tupyy/gophoto/internal/services/album"
	mediaService "github.com/tupyy/gophoto/internal/services/media"
	tagService "github.com/tupyy/gophoto/internal/services/tag"
	"github.com/tupyy/gophoto/internal/services/users"
	"github.com/tupyy/gophoto/utils/logutil"
)

func RegisterIndexHandler(privateGroup *gin.RouterGroup, us *users.Service) {
	// register index controller
	index.Index(privateGroup, us)
	logutil.GetDefaultLogger().Info("index controller registered")
}

func RegisterAlbumHandler(privateGroup *gin.RouterGroup, as *albumService.Service, us *users.Service, t *tagService.Service) {

	// register album api
	album.GetAlbums(privateGroup, as, us)
	album.GetAlbumsTags(privateGroup, as, t)
	logutil.GetDefaultLogger().Info("api album registered")

	// album handlers
	album.GetAlbum(privateGroup, as, us, t)
	album.GetCreateAlbumForm(privateGroup, us)
	album.GetUpdateAlbumForm(privateGroup, as, us)
	album.CreateAlbum(privateGroup, as)
	album.UpdateAlbum(privateGroup, as)
	album.DeleteAlbum(privateGroup, as)
	album.Thumbnail(privateGroup, as)
	logutil.GetDefaultLogger().Info("album handlers registered")
}

func RegisterMediaHandler(privateGroup *gin.RouterGroup, as *albumService.Service, ms *mediaService.Service) {
	media.UploadMedia(privateGroup, as, ms)
	media.DownloadMedia(privateGroup, as, ms)
	media.GetAlbumMedia(privateGroup, as, ms)

	media.DeleteMedia(privateGroup, as, ms)

	logutil.GetDefaultLogger().Info("api media registered")
}

func RegisterTagHandler(privateGroup *gin.RouterGroup, as *albumService.Service, t *tagService.Service) {
	tag.Get(privateGroup, t)
	tag.Crud(privateGroup, t)

	tag.Dissociate(privateGroup, as, t)
	tag.Associate(privateGroup, as, t)

	logutil.GetDefaultLogger().Info("api tag registered")
}
