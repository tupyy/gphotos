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
	ID        string `validate:"required"`
	Username  string `validate:"required"`
	FirstName string
	LastName  string
	Role      Role
	CanShare  bool
	Groups    []Group
}

func (u User) Validate() error {
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("%w %v", ErrInvalidEntity, err)
	}

	return nil
}
