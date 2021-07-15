package postgres

import "errors"

var (
	// ErrInternalError means that something bad happened.
	ErrInternalError = errors.New("internal error")

	// ErrUserNotFound means that user was not found.
	ErrUserNotFound = errors.New("user not found")

	// Album repos errors
	ErrAlbumNotFound = errors.New("album not found")
	ErrCreateAlbum   = errors.New("cannot create album")
	ErrUpdateAlbum   = errors.New("cannot update album")
)
