package domain

import "errors"

var (
	ErrNotImplementated = errors.New("not implementated")

	// ErrInternalError means that something bad happened.
	ErrInternalError = errors.New("internal error")

	// ErrNotFound resource not found
	ErrNotFound = errors.New("resource not found")

	// Album repos errors
	ErrCreateAlbum = errors.New("cannot create album")
	ErrUpdateAlbum = errors.New("cannot update album")
	ErrDeleteAlbum = errors.New("error deleting album")

	// Minio
	ErrBucketAlreadyExists = errors.New("bucket exists")
)
