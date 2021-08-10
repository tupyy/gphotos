package entity

import (
	"fmt"
	"strings"
	"time"
)

type Album struct {
	// ID - id of the album
	ID int32
	// Name - name of the album
	Name string `validate:"required"`
	// CreateAt - creation date
	CreatedAt time.Time `validate:"required"`
	// OwnerID - owner-s id
	OwnerID string `validate:"required"`
	// Description - description of the album
	Description string
	// Location - location of the albul
	Location string
	// UserPermissions - holds the list of permissions of other users for this album.
	// The key is the user id.
	UserPermissions Permissions
	// GroupPermissions - holds the list of permissions of groups for this album.
	// The key is the group name.
	GroupPermissions Permissions
}

func (a Album) Validate() error {
	if err := validate.Struct(a); err != nil {
		return fmt.Errorf("%w album not valid %v", ErrInvalidEntity, err)
	}

	return nil
}

func (a Album) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "id = %d\n", a.ID)
	fmt.Fprintf(&sb, "name = %s\n", a.Name)
	fmt.Fprintf(&sb, "created_at = %+v\n", a.CreatedAt)
	fmt.Fprintf(&sb, "description = %s\n", a.Description)
	fmt.Fprintf(&sb, "location = %s\n", a.Location)

	for k, v := range a.UserPermissions {
		fmt.Fprintf(&sb, "user = %s, permisions = %+v\n", k, v)
	}

	for k, v := range a.GroupPermissions {
		fmt.Fprintf(&sb, "group = %s, permisions = %+v\n", k, v)
	}

	return sb.String()

}
