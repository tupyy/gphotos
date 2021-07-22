package middleware

import (
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/utils/logutil"
)

func HotReloading(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		renderer, err := loadTemplates(conf.GetTemplateFolder())
		if err != nil {
			panic(err)
		}

		r.HTMLRender = renderer
	}
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

	for _, t := range templates {
		layoutCopy := make([]string, len(layouts)+1)

		copy(layoutCopy[1:], layouts)
		layoutCopy[0] = t

		r.AddFromFiles(filepath.Base(t), layoutCopy...)

		logger.WithFields(logrus.Fields{
			"template": filepath.Base(t),
			"files":    layoutCopy,
		}).Debug("add template")
	}

	return r, nil

}
