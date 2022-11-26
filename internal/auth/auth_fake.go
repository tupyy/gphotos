package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/entity"
	"go.uber.org/zap"
)

// FakeAuthMiddleware reads the cookie and unmarshall it into session.
// The cookie must be encoded base64.
func FakeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		sessionEncoded := c.Request.Header.Get("SESSIONID")

		payload, err := base64.StdEncoding.DecodeString(sessionEncoded)
		if err != nil {
			zap.S().Errorw("cannot decode cookie", "error", err)
			c.Abort()
			return
		}
		var se entity.Session
		if err := json.Unmarshal(payload, &se); err != nil {
			zap.S().Errorw("cannot unmarshal cookie to session", "error", err)
			c.Abort()
			return
		}
		session.Set(sessionEncoded, se)
		session.Save()

		c.Set("session", se)
		c.Next()
	}
}
