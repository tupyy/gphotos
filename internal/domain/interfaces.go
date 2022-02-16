package domain

import (
	"context"
	"io"

	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
	"github.com/tupyy/gophoto/internal/entity"
)

type Repositories map[RepoName]interface{}

type RepoName int

const (
	KeycloakRepoName RepoName = iota
	AlbumRepoName
	UserRepoName
	MinioRepoName
	TagRepoName
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
	Get(ctx context.Context) ([]entity.Album, error)
	// GetByID return an album by id.
	GetByID(ctx context.Context, id int32) (entity.Album, error)
	// GetByOwner return all albums of a user for which he is the owner.
	GetByOwnerID(ctx context.Context, ownerID string) ([]entity.Album, error)
	// GetByUserID returns a list of albums for which the user has at least one permission set.
	GetByUserID(ctx context.Context, userID string) ([]entity.Album, error)
	// GetByGroup returns a list of albums for which the group has at least one permission.
	GetByGroupName(ctx context.Context, groupName string) ([]entity.Album, error)
	// GetByGroups returns a list of albums with at least one persmission for at least on group in the list.
	GetByGroups(ctx context.Context, groups []string) ([]entity.Album, error)
}

// Postgres repo to handler relationships between users
type User interface {
	GetRelatedUsers(ctx context.Context, user entity.User) (ids []string, err error)
}

// Store describe photo store operations
type Store interface {
	// GetFile returns a reader to file.
	GetFile(ctx context.Context, bucket, filename string) (io.ReadSeeker, map[string]string, error)
	// PutFile save a file to a bucket.
	PutFile(ctx context.Context, bucket, filename string, size int64, r io.Reader, metadata map[string]string) error
	// ListFiles list the content of a bucket
	ListBucket(ctx context.Context, bucket string) ([]entity.Media, error)
	// DeleteFile deletes a file from a bucket.
	DeleteFile(ctx context.Context, bucket, filename string) error
	// CreateBucket create a bucket.
	CreateBucket(ctx context.Context, bucket string) error
	// DeleteBucket removes bucket.
	DeleteBucket(ctx context.Context, bucket string) error
}

type Tag interface {
	// Create -- create the tag.
	Create(ctx context.Context, tag entity.Tag) (int32, error)
	// Update -- update the tag.
	Update(ctx context.Context, tag entity.Tag) error
	// Delete -- delete the tag. it does not cascade.
	Delete(ctx context.Context, id int32) error
	// GetByUser -- fetch all user's tags
	GetByUser(ctx context.Context, userID string) ([]entity.Tag, error)
	// GetByName -- fetch the tag by name and user id.
	GetByName(ctx context.Context, userID, name string) (entity.Tag, error)
	// GetByID -- fetch the tag by id
	GetByID(ctx context.Context, userID string, id int32) (entity.Tag, error)
	// GetByAlbum -- fetch all user's tag for the album
	GetByAlbum(ctx context.Context, albumID int32) ([]entity.Tag, error)
	// AssociateTag -- associates a tag with an album.
	Associate(ctx context.Context, albumID, tagID int32) error
	// Dissociate -- removes a tag from an album.
	Dissociate(ctx context.Context, albumID, tagID int32) error
}
