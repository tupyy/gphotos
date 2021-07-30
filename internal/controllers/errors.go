package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/utils/logutil"
)

func Abort(c *gin.Context, status int, err error, msg string) {

	logutil.GetLogger(c).WithError(err).Errorf("abort %s with code %d: %s", c.FullPath(), status, msg)
	c.AbortWithError(status, err)
}

func AbortBadRequest(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusBadRequest, err, msg)
}

func AbortInternalError(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusInternalServerError, err, msg)
}

func AbortForbidden(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusForbidden, err, msg)
}

func AbortNotFound(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusNotFound, err, msg)
}
