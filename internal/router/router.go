package router

import (
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

type PhotoRouter struct {
	r            *gin.Engine
	PrivateGroup *gin.RouterGroup
	PublicGroup  *gin.RouterGroup
}

// NewRouter returns a new gin router.
func NewRouter(store sessions.Store, authenticator auth.Authenticator) *PhotoRouter {
	r := gin.Default()

	r.Use(sessions.Sessions("gophoto", store))

	if gin.Mode() == "debug" {
		logutil.GetDefaultLogger().Debug("loading statics")
		r.Static("/static", conf.GetStaticsFolder())
	}

	// setup authentication for the priate group.
	private := r.Group("/", authenticator.AuthMiddleware())

	// create a public group.
	public := r.Group("/public")

	// set auth callback
	url, err := url.Parse(conf.GetServerAuthCallback())
	if err != nil {
		panic(err)
	}

	r.GET(url.RequestURI(), authenticator.Callback())

	return &PhotoRouter{r: r, PrivateGroup: private, PublicGroup: public}
}

func (p *PhotoRouter) Run() {
	p.r.Run(":8080")
}
