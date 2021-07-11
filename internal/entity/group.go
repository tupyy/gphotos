package entity

import (
	"fmt"
)

type Group struct {
	ID   *int32 `validate:"required"`
	Name string `validate:"required"`
}

func (g Group) Validate() error {
	err := validate.Struct(g)
	if err != nil {
		return fmt.Errorf("%w %s", ErrInvalidEntity, err)
	}

	return nil
}
