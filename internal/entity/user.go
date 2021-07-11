package entity

import (
	"fmt"
)

type Role int

const (
	RoleUser Role = iota
	RoleAdmin
	RoleEditor
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	case RoleEditor:
		return "editor"
	}

	return "unknown"
}

type User struct {
	ID       *int32 `validate:"required"`
	UserID   string `validate:"required"`
	Username string `validate:"required"`
	Role     Role
	CanShare bool
	Groups   []Group
}

func (u User) Validate() error {
	err := validate.Struct(u)
	if err != nil {
		return fmt.Errorf("%w %v", ErrInvalidEntity, err)
	}

	return nil
}
