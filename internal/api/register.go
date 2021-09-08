package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/api/album"
	"github.com/tupyy/gophoto/internal/api/media"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/utils/logutil"
)

func RegisterApi(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, repos domain.Repositories) {

	// register album api
	album.GetAlbums(privateGroup, repos)
	logutil.GetDefaultLogger().Info("api album registered")

	// register media upload api
	media.UploadMedia(privateGroup, repos)
	logutil.GetDefaultLogger().Info("api media registered")
}
