package router

import (
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/conf"
)

type PhotoRouter struct {
	PrivateGroup *gin.RouterGroup
	PublicGroup  *gin.RouterGroup
}

func InitEngine(server *gin.Engine, store sessions.Store, authenticator auth.Authenticator, middlewares ...gin.HandlerFunc) {
	server.Use(sessions.Sessions("gophoto", store))

	if gin.Mode() == "debug" {
		server.Static("/static", conf.GetStaticsFolder())
	}

	server.LoadHTMLFiles("static/index.html")

	server.Use(gin.Recovery())
	server.Use(auth.FakeAuthMiddleware())

	// set auth callback
	url, err := url.Parse(conf.GetServerAuthCallback())
	if err != nil {
		panic(err)
	}

	server.GET(url.RequestURI(), authenticator.Callback())
}
