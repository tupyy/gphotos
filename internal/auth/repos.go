package auth

import (
	"context"

	"github.com/tupyy/gophoto/internal/entity"
)

// UserRepo represents the interface for user repo used in auth packages.
type UserRepo interface {
	// Create creates the user the return the id of created user or error.
	Create(ctx context.Context, user entity.User) (int32, error)
	// Update updates the user.
	Update(ctx context.Context, user entity.User) error
	// GetByUsername return the user by username
	GetByUsername(ctx context.Context, username string) (entity.User, error)
}

// GroupRepo represents the interface for group repo used in auth packages.
type GroupRepo interface {
	// Create creates the group and return the id of created group or error.
	Create(ctx context.Context, group entity.Group) (int32, error)
	// Updates the group.
	Update(ctx context.Context, group entity.Group) error
	// Delete removes the group.
	Delete(ctx context.Context, groupID int32) error
	// GetByName returns a group by name.
	GetByName(ctx context.Context, name string) (entity.Group, error)
}
