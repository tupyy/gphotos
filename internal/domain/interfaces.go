package domain

import (
	"context"

	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/filters"
	"github.com/tupyy/gophoto/internal/domain/sort"
)

type Repositories map[RepoName]interface{}

type RepoName int

const (
	KeycloakRepoName RepoName = iota
	AlbumRepoName
)

type KeycloakRepo interface {
	// Get returns all the users.
	GetUsers(ctx context.Context, sorter sort.UserSorter, filters ...filters.UserFilter) ([]entity.User, error)
	// GetByID return the user by id.
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	// // GetByUsername return the user by username
	// GetByUsername(ctx context.Context, username string) (entity.User, error)
	// // GetByGroupID returns all the users belonging to groupID.
	// GetByGroupID(ctx context.Context, groupID int32) ([]entity.User, error)
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
	Get(ctx context.Context, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error)
	// GetByID return an album by id.
	GetByID(ctx context.Context, id int32) (entity.Album, error)
	// GetByOwner return all albums of a user for which he is the owner.
	GetByOwnerID(ctx context.Context, ownerID string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error)
	// GetByUserID returns a list of albums for which the user has at least one permission set.
	GetByUserID(ctx context.Context, userID string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error)
	// GetByGroup returns a list of albums for which the group has at least one permission.
	GetByGroupName(ctx context.Context, groupName string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error)
}
