package common

import (
	"errors"
	"fmt"
)

const (
	InternalError        string = "internal error"
	EntityNotFound       string = "entity not found"
	PostgresNotAvailable string = "postgres not available"
	MinioNotAvailable    string = "minio not available"
	KeycloakNotAvailable string = "keycloak not available"
	Forbidden            string = "Forbidden"
)

type ServiceError struct {
	error
	Cause   string
	Message string
}

func (s ServiceError) Unwrap() error {
	return s.error
}

func (s ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", s.Cause, s.Message)
}

func NewInternalError(err error, msg string) ServiceError {
	return ServiceError{
		error:   err,
		Cause:   InternalError,
		Message: msg,
	}
}

func NewEntityNotFound(msg string) ServiceError {
	return ServiceError{
		error:   errors.New("entity not found"),
		Cause:   EntityNotFound,
		Message: msg,
	}
}

func NewPostgresNotAvailableError(msg string) ServiceError {
	return ServiceError{
		error:   errors.New("postgres not available"),
		Cause:   PostgresNotAvailable,
		Message: msg,
	}
}
