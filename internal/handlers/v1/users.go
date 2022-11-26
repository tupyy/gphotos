package v1

import (
	"github.com/gin-gonic/gin"

	apiv1 "github.com/tupyy/gophoto/api/v1"
)

// (GET /api/gphotos/v1/users)
func (server *Server) GetUsers(c *gin.Context) {}

// (GET /api/gphotos/v1/users/{user_id}/groups/related)
func (server *Server) GetRelatedGroups(c *gin.Context, userId apiv1.UserId) {}

// (GET /api/gphotos/v1/users/{user_id}/related)
func (server *Server) GetRelatedUsers(c *gin.Context, userId apiv1.UserId) {}

// (GET /api/gphotos/v1/groups)
func (server *Server) GetGroups(c *gin.Context) {}
