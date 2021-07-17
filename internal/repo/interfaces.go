package repo

import (
	"context"

	"github.com/tupyy/gophoto/internal/entity"
)

type Repositories map[RepoName]interface{}

type RepoName int

const (
	UserRepoName RepoName = iota
	GroupRepoName
	AlbumRepoName
)

type UserRepo interface {
	// Create creates the user the return the id of created user or error.
	Create(ctx context.Context, user entity.User) (int32, error)
	// Update updates the user.
	Update(ctx context.Context, user entity.User) error
	// Delete removes the user.
	Delete(ctx context.Context, userID int32) error
	// Get returns all the users.
	Get(ctx context.Context) ([]entity.User, error)
	// GetByID return the user by id.
	GetByID(ctx context.Context, id int32) (entity.User, error)
	// GetByUsername return the user by username
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	// GetByGroupID returns all the users belonging to groupID.
	GetByGroupID(ctx context.Context, groupID int32) ([]entity.User, error)
}

type GroupRepo interface {
	// Create creates the group and return the id of created group or error.
	Create(ctx context.Context, group entity.Group) (int32, error)
	// Updates the group.
	Update(ctx context.Context, group entity.Group) error
	// Delete removes the group.
	Delete(ctx context.Context, groupID int32) error
	// Get returns all the groups.
	Get(ctx context.Context) ([]entity.Group, error)
	// GetByID return a group by id.
	GetByID(ctx context.Context, id int32) (entity.Group, error)
	// GetByName returns a group by name.
	GetByName(ctx context.Context, name string) (entity.Group, error)
	// GetByUserID returns all the groups of the user.
	GetByUserID(ctx context.Context, userID string) ([]entity.Group, error)
}

type AlbumRepo interface {
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
	GetByOwnerID(ctx context.Context, ownerID int32) ([]entity.Album, error)
	// GetByUser returns a list of album for which the user ether has a permission set or it is the owner.
	GetByUserID(ctx context.Context, userID int32) ([]entity.Album, error)
	// GetByGroup returns a list of albums for which the group has at least one permission.
	GetByGroupID(ctx context.Context, groupID int32) ([]entity.Album, error)
}
