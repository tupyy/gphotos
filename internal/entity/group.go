package entity

import (
	"fmt"
)

type Group struct {
	Name  string `validate:"required"`
	Users []User
}

func (g Group) Validate() error {
	err := validate.Struct(g)
	if err != nil {
		return fmt.Errorf("%w %s", ErrInvalidEntity, err)
	}

	return nil
}
