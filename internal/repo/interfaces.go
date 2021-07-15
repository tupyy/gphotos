package repo

import (
	"context"

	"github.com/tupyy/gophoto/internal/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user entity.User) (int, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Get(ctx context.Context, username string) (entity.User, error)
}

type GroupRepo interface {
	// FirstOrCreate returns the first group found. If not found it creates the group and return the entity.
	FirstOrCreate(ctx context.Context, name string) (entity.Group, bool, error)
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
	// GetAlbumsByOwner return all albums of a user for which he is the owner.
	GetAlbumsByOwnerID(ctx context.Context, ownerID int32) ([]entity.Album, error)
	// GetAlbumsByUser returns a list of album for which the user ether has a permission set or it is the owner.
	GetAlbumsByUserID(ctx context.Context, userID int32) ([]entity.Album, error)
	// GetAlbumsByGroup returns a list of albums for which the group has at least one permission.
	GetAlbumsByGroupID(ctx context.Context, groupID int32) ([]entity.Album, error)
}
