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
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("%w %v", ErrInvalidEntity, err)
	}

	return nil
}

type UserFilter struct {
	users []User
}

type UserFilterFunc func(u User) bool

func NewUserFilter(users []User) *UserFilter {
	return &UserFilter{users}
}

func (uf *UserFilter) Filter(filterFunc UserFilterFunc) []User {
	filteredUsers := make([]User, 0, len(uf.users))

	for _, u := range uf.users {
		if filterFunc(u) {
			filteredUsers = append(filteredUsers, u)
		}
	}

	return filteredUsers
}
