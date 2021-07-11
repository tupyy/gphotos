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
