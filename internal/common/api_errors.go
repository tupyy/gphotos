package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonError struct {
	Message string `json:"message"`
}

func AbortWithJson(c *gin.Context, status int, err error, msg string) {
	c.AbortWithStatusJSON(status, jsonError{msg})
}

func AbortBadRequestWithJson(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusBadRequest, err, msg)
}

func AbortInternalErrorWithJson(c *gin.Context) {
	AbortWithJson(c, http.StatusInternalServerError, errors.New("internal error"), "internal error")
}

func AbortForbiddenWithJson(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusForbidden, err, msg)
}

func AbortNotFoundWithJson(c *gin.Context, err error, msg string) {
	AbortWithJson(c, http.StatusNotFound, err, msg)
}
