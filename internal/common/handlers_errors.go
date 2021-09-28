package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

var (
	ErrDeleteAlbum    = errors.New("error deleting album")
	ErrEncryptionData = errors.New("error encrypting data")
)

type MissingPermissionError struct {
	error
	requiredPermission entity.Permission
	album              entity.Album
	user               entity.User
}

func NewMissingPermissionError(requiredPermission entity.Permission, album entity.Album, user entity.User) *MissingPermissionError {
	return &MissingPermissionError{
		requiredPermission: requiredPermission,
		album:              album,
		user:               user,
	}
}

func (p *MissingPermissionError) Error() string {
	return fmt.Sprintf("User %s is missing permission %s to access album %v", p.user.ID, p.requiredPermission.String(), p.album)
}

func Abort(c *gin.Context, status int, err error, msg string) {

	logutil.GetLogger(c).WithError(err).Errorf("abort %s with code %d: %s", c.FullPath(), status, msg)
	c.AbortWithError(status, err)
}

func AbortBadRequest(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusBadRequest, err, msg)
}

func AbortInternalError(c *gin.Context) {
	Abort(c, http.StatusInternalServerError, errors.New("internal error"), "")
}

func AbortForbidden(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusForbidden, err, msg)
}

func AbortNotFound(c *gin.Context, err error, msg string) {
	Abort(c, http.StatusNotFound, err, msg)
}
