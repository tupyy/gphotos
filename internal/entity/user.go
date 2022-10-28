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
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Role      Role    `json:"role"`
	CanShare  bool    `json:"can_share"`
	Groups    []Group `json:"groups"`
}

func (u User) Validate() error {
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("%w %v", ErrInvalidEntity, err)
	}

	return nil
}
