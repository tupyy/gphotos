package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/services/album"
	"github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/internal/services/tag"
	"github.com/tupyy/gophoto/internal/services/users"
)

type Server struct {
	albumService     *album.Service
	userService      *users.Service
	tagService       *tag.Service
	mediaService     *media.Service
	encryptionServer EncryptionService
}

func NewServer(a *album.Service, u *users.Service, tag *tag.Service, m *media.Service, e EncryptionService) *Server {
	return &Server{a, u, tag, m, e}
}

func (server *Server) AlbumService() *album.Service {
	return server.albumService
}

func (server *Server) UserService() *users.Service {
	return server.userService
}

func (server *Server) TagService() *tag.Service {
	return server.tagService
}

func (server *Server) MediaService() *media.Service {
	return server.mediaService
}

func (server *Server) EncryptionService() EncryptionService {
	return server.encryptionServer
}

// (GET /api/gphotos/v1)
func (server *Server) GetVersionMetadata(c *gin.Context) {

}
