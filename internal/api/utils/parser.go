package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// parseAlbumIDHandler decrypt the album id passes as parameters and set the id in the context.
func ParseAlbumIDHandler(c *gin.Context) {
	logger := logutil.GetLogger(c)

	// decrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	decryptedID, err := gen.DecryptData(c.Param("id"))
	if err != nil {
		logger.WithError(err).Error("cannot decrypt album id")
		c.AbortWithError(http.StatusNotFound, err) // explicit return not found here

		return
	}

	id, err := strconv.Atoi(decryptedID)
	if err != nil {
		logger.WithError(err).WithField("id", decryptedID).Error("cannot parse album id")
		c.AbortWithError(http.StatusNotFound, err)

		return
	}

	c.Set("id", id)
}
