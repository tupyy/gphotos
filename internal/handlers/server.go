package handlers

import (
	"github.com/gin-gonic/gin"
	apiv1 "github.com/tupyy/gophoto/api/v1"
)

type Server struct {
	services map[string]interface{}
}

func NewServer(services map[string]interface{}) *Server {
	return &Server{services: services}
}

// (GET /api/gphotos/v1)
func (s *Server) GetVersionMetadata(c *gin.Context) {

}

// (GET /api/gphotos/v1/albums)
func (s *Server) GetAlbums(c *gin.Context, params apiv1.GetAlbumsParams) {}

// (GET /api/gphotos/v1/albums/groups/{group_id})
func (s *Server) GetAlbumsByGroup(c *gin.Context, groupId string, params apiv1.GetAlbumsByGroupParams) {
}

// (GET /api/gphotos/v1/albums/users/{user_id})
func (s *Server) GetAlbumsByUser(c *gin.Context, userId string, params apiv1.GetAlbumsByUserParams) {}

// (GET /api/gphotos/v1/albums/{album_id}/permissions)
func (s *Server) GetAlbumPermissions(c *gin.Context, albumId string, params apiv1.GetAlbumPermissionsParams) {
}

// (GET /api/gphotos/v1/auth/callback)
func (s *Server) GetApiGphotosV1AuthCallback(c *gin.Context, params apiv1.GetApiGphotosV1AuthCallbackParams) {
}
