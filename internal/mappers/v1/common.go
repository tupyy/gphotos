package v1

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/common"
)

type EncryptionService interface {
	// Encrypt data in a deterministic way.
	Encrypt(data string) (string, error)
	// Decrypt data.
	Decrypt(data string) (string, error)
}

const (
	AlbumKind            string = "Album"
	AlbumListKind        string = "AlbumList"
	AlbumPermissionsKind string = "AlbumPermissionsList"
	PhotoKind            string = "Photo"
	PhotoListKind        string = "PhotoList"
	UserKind             string = "User"
	GroupKind            string = "Group"
	TagKind              string = "Tag"
	TagListKind          string = "TagList"
)

func MapFromError(err error) apiv1.Error {
	apiError := apiv1.Error{
		Kind: "Error",
		Href: "api/v1/errors/",
	}
	switch v := err.(type) {
	case common.ServiceError:
		apiError.Reason = &v.Message
		switch v.Cause {
		case common.EntityNotFound:
			apiError.Code = 404
		default:
			apiError.Code = 500
		}
	default:
		apiError.Id = "xxx"
		apiError.Code = 500
		reason := err.Error()
		apiError.Reason = &reason
	}
	return apiError
}

func MapFromStatus(code int, msg string) apiv1.Error {
	return apiv1.Error{
		Id:     "xxx",
		Kind:   "Error",
		Href:   "api/v1/errors/",
		Code:   code,
		Reason: &msg,
	}
}

func MapFromStatusf(code int, format string, args ...interface{}) apiv1.Error {
	msg := fmt.Sprintf(format, args...)
	return apiv1.Error{
		Id:     "xxx",
		Kind:   "Error",
		Href:   "api/v1/errors/",
		Code:   code,
		Reason: &msg,
	}
}
