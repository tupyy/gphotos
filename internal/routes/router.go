package router

import (
	"html/template"
	"net/url"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/templates/funcs"
	"github.com/tupyy/gophoto/utils/logutil"
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

	// load templates
	logutil.GetDefaultLogger().Debug("loading templates")

	renderer, err := loadTemplates(conf.GetTemplateFolder())
	if err != nil {
		panic(err)
	}

	r.HTMLRender = renderer

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

func loadTemplates(templateDir string) (multitemplate.Renderer, error) {
	logger := logutil.GetDefaultLogger()
	logger.WithField("template dir", templateDir).Debug("load templates")

	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templateDir + "/layout/*.html")
	if err != nil {
		return r, err
	}

	templates, err := filepath.Glob(templateDir + "/*.html")
	if err != nil {
		return r, err
	}

	templateFuncs := template.FuncMap{
		"day":              funcs.Day,
		"month":            funcs.Month,
		"year":             funcs.Year,
		"perm_name":        funcs.PermissionName,
		"date":             funcs.Date,
		"date_photo":       funcs.DatePhoto,
		"extract_metadata": funcs.ExtractMetadata,
	}

	for _, t := range templates {
		layoutCopy := make([]string, len(layouts)+1)

		copy(layoutCopy[1:], layouts)
		layoutCopy[0] = t

		r.AddFromFilesFuncs(filepath.Base(t), templateFuncs, layoutCopy...)

		logger.WithFields(logrus.Fields{
			"template": filepath.Base(t),
			"files":    layoutCopy,
		}).Debug("add template")
	}

	return r, nil

}
