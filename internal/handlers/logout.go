package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/auth"
	"github.com/tupyy/gophoto/internal/entity"
)

func Logout(r *gin.RouterGroup, k auth.Authenticator) {
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)

		cookie, err := c.Request.Cookie("sessionID")
		if err != nil {
			logrus.WithError(err).Error("error reading cookie")

			http.Redirect(c.Writer, c.Request, "/", http.StatusFound)

			return
		}

		s := session.Get(cookie.Value)
		if s == nil {
			logrus.WithField("sessionID", cookie.Value).Error("no session with this id")

			http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
			return
		}

		sessionData, _ := s.(entity.Session)

		// logout from keycloak
		if err := k.Logout(c, sessionData.User.Username, sessionData.Token.RefreshToken); err != nil {
			logrus.WithError(err).Error("error logging out user")
		}

		session.Delete(cookie.Value)
		session.Save()

		c.SetCookie("sessionID", "", 3600, "/", "localhost", true, true)

		logrus.WithField("username", c.GetString("username")).Debug("user logged out")
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
