package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/users"
)

type Server struct {
	services map[string]interface{}
}

func NewServer(services map[string]interface{}) *Server {
	return &Server{services: services}
}

func (server *Server) GetAlbumService() album.Service {
	return server.services["album"].(album.Service)
}

func (server *Server) GetUserService() users.Service {
	return server.services["user"].(users.Service)
}

// (GET /api/gphotos/v1)
func (server *Server) GetVersionMetadata(c *gin.Context) {

}

func (Server *Server) GetAlbumPhotos(c *gin.Context, albumID string) {
}
