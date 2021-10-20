package services

import "errors"

// Album service errors
var (
	// ErrCreateBucket means that the bucket cannot be created.
	ErrCreateBucket = errors.New("failed to create bucket")
	// ErrDeleteBucket means the removing of the bucket has failed.
	ErrDeleteBucket = errors.New("failed to delete bucket")
	// ErrListBucket means the bucket cannot be read.
	ErrListBucket = errors.New("failed to list bucket")

	// ErrCreateAlbum means that the album cannot be create.
	ErrCreateAlbum = errors.New("failed to create album")
	// ErrUpdateAlbum means the album cannot be updated.
	ErrUpdateAlbum = errors.New("failed to update album")
	// ErrDeleteAlbum means the albums cannot be delete it.
	ErrDeleteAlbum = errors.New("failed to delete album")
	// ErrGetAlbums means that we cannot fetch album from repo.
	ErrGetAlbums = errors.New("failed to get albums")
)
