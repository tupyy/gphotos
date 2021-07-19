package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/repo"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Register(privateGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, repos repo.Repositories) {

	// register index controller
	Index(privateGroup, repos[repo.AlbumRepoName].(repo.AlbumRepo))
	logutil.GetDefaultLogger().Info("index controller registered")

	CreateAlbum(privateGroup, repos)
	logutil.GetDefaultLogger().Info("album controller registered")
}
