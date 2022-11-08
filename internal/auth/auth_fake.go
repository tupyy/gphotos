package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

// FakeAuthMiddleware reads the cookie and unmarshall it into session.
// The cookie must be encoded base64.
func FakeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		sessionEncoded := c.Request.Header.Get("SESSIONID")

		logger := logutil.GetLogger(c)
		payload, err := base64.StdEncoding.DecodeString(sessionEncoded)
		if err != nil {
			logger.WithError(err).Error("cannot decode cookie")
			c.Abort()
			return
		}
		var se entity.Session
		if err := json.Unmarshal(payload, &se); err != nil {
			logger.WithError(err).Error("cannot unmarshal cookie to session")
			c.Abort()
			return
		}
		session.Set(sessionEncoded, se)
		session.Save()

		c.Set("session", se)
		c.Next()
	}
}
