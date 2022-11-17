package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/tag"
	"github.com/tupyy/gophoto/internal/services/users"
)

type Server struct {
	albumService *album.Service
	userService  *users.Service
	tagService   *tag.Service
}

func NewServer(a *album.Service, u *users.Service, tag *tag.Service) *Server {
	return &Server{a, u, tag}
}

func (server *Server) GetAlbumService() *album.Service {
	return server.albumService
}

func (server *Server) GetUserService() *users.Service {
	return server.userService
}

func (server *Server) TagService() *tag.Service {
	return server.tagService
}

// (GET /api/gphotos/v1)
func (server *Server) GetVersionMetadata(c *gin.Context) {

}

func (Server *Server) GetAlbumPhotos(c *gin.Context, albumID string) {
}
