package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/utils/logutil"
)

func AbortWithJson(c *gin.Context, status int, err error, msg string) {

	logutil.GetLogger(c).WithError(err).Errorf("abort %s with code %d: %s", c.FullPath(), status, msg)
	c.AbortWithStatusJSON(status, err)
}

func AbortBadRequest(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusBadRequest, err, msg)
}

func AbortInternalError(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusInternalServerError, err, msg)
}

func AbortForbidden(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusForbidden, err, msg)
}

func AbortNotFound(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusNotFound, err, msg)
}
