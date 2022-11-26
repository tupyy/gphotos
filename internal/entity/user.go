package entity

import (
	"fmt"
)

type Role string

const (
	RoleUser   Role = "user"
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
)

func (r Role) String() string {
	return string(r)
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
