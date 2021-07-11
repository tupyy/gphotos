package repo

import (
	"context"

	"github.com/tupyy/gophoto/internal/entity"
)

type KeycloakRepo interface {
	// Logout logs out the user from ALL sessions.
	Logout(accessToken, refreshToken, userID string) error
}

type UserRepo interface {
	// Add insert a new user into db.
	Create(ctx context.Context, user entity.User) (int, error)
	// Update updates a user entry.
	Update(ctx context.Context, user entity.User) (entity.User, error)
	// Get returns the user if found.
	Get(ctx context.Context, username string) (entity.User, error)
}

type GroupRepo interface {
	// FirstOrCreate returns the first group found. If not found it creates the group and return the entity.
	FirstOrCreate(ctx context.Context, name string) (entity.Group, bool, error)
}
