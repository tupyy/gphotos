package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

// AlbumIDMiddleware decrypt the album id passes as parameters and set the id in the context.
func DecryptID(c *gin.Context) {
	// decrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	for _, param := range c.Params {
		if strings.HasSuffix(param.Key, "id") {
			decryptedID, err := gen.DecryptData(param.Value)
			if err != nil {
				c.Set(param.Key, param.Value)
				continue
			}
			c.Set(param.Key, decryptedID)
		}
	}
}
