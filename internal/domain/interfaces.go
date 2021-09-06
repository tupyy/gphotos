package domain

import (
	"context"
	"io"

	"github.com/tupyy/gophoto/internal/domain/entity"
	albumFilters "github.com/tupyy/gophoto/internal/domain/filters/album"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
)

type Repositories map[RepoName]interface{}

type RepoName int

const (
	KeycloakRepoName RepoName = iota
	AlbumRepoName
	UserRepoName
	BucketRepoName
	MinioRepoName
)

type KeycloakRepo interface {
	// Get returns all the users.
	GetUsers(ctx context.Context, filters userFilters.Filters) ([]entity.User, error)
	// GetByID return the user by id.
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	// GetGroups returns all groups.
	GetGroups(ctx context.Context) ([]entity.Group, error)
}

type Album interface {
	// Create creates an album.
	Create(ctx context.Context, album entity.Album) (albumID int32, err error)
	// Update an album.
	Update(ctx context.Context, album entity.Album) error
	// Delete removes an album from postgres.
	Delete(ctx context.Context, id int32) error
	// Get return all the albums.
	Get(ctx context.Context, filters albumFilters.Filters) ([]entity.Album, error)
	// GetByID return an album by id.
	GetByID(ctx context.Context, id int32) (entity.Album, error)
	// GetByOwner return all albums of a user for which he is the owner.
	GetByOwnerID(ctx context.Context, ownerID string, filters albumFilters.Filters) ([]entity.Album, error)
	// GetByUserID returns a list of albums for which the user has at least one permission set.
	GetByUserID(ctx context.Context, userID string, filters albumFilters.Filters) ([]entity.Album, error)
	// GetByGroup returns a list of albums for which the group has at least one permission.
	GetByGroupName(ctx context.Context, groupName string, filters albumFilters.Filters) ([]entity.Album, error)
	// GetByGroups returns a list of albums with at least one persmission for at least on group in the list.
	GetByGroups(ctx context.Context, groups []string, filters albumFilters.Filters) ([]entity.Album, error)
}

// Bucket describe postgres operation on bucket table.
type Bucket interface {
	// Create a bucket
	Create(ctx context.Context, bucket entity.Bucket) error
	// Delete bucket from postgres
	Delete(ctx context.Context, bucket entity.Bucket) error
}

// Postgres repo to handler relationships between users
type User interface {
	GetRelatedUsers(ctx context.Context, user entity.User) (ids []string, err error)
}

// Store describe photo store operations
type Store interface {
	// GetFile returns a reader to file.
	// GetFile(ctx context.Context, bucket, filename string) (io.Reader, error)
	// PutFile save a file to a bucket.
	PutFile(ctx context.Context, bucket, filename string, size int64, r io.Reader) error
	// // ListFiles list the content of a bucket
	// ListFiles(ctx context.Context, bucket string) ([]string, error)
	// // DeleteFile deletes a file from a bucket.
	// DeleteFile(ctx context.Context, bucket, filename string) error
	// CreateBucket create a bucket.
	CreateBucket(ctx context.Context, bucket string) error
	// DeleteBucket removes bucket.
	DeleteBucket(ctx context.Context, bucket string) error
}
