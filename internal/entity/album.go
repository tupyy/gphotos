package entity

import (
	"fmt"
	"strings"
	"time"
)

type Album struct {
	ID          int32
	Name        string    `validate:"required"`
	CreatedAt   time.Time `validate:"required"`
	OwnerID     string    `validate:"required"`
	Description string
	Location    string
	// UserPermissions holds the list of permissions of other users for this album.
	// The key is the user id.
	UserPermissions Permissions
	// GroupPermissions holds the list of permissions of groups for this album.
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

// HasUserPermissions returns true if user has at least one permission.
func (a Album) HasUserPermissions(userID string) bool {
	_, found := a.UserPermissions[userID]

	return found
}

func (a Album) HasUserPermission(userID string, permission Permission) bool {
	if !a.HasUserPermissions(userID) {
		return false
	}

	for _, p := range a.UserPermissions[userID] {
		if p == permission {
			return true
		}
	}

	return false
}

// GetUserPermissions returns all permission of user.
func (a Album) GetUserPermissions(userID string) (permissions []Permission, found bool) {
	if _, found = a.UserPermissions[userID]; !found {
		return
	}

	return a.UserPermissions[userID], true
}

func (a Album) HasGroupPermission(groupName string, permission Permission) bool {
	if !a.HasGroupPermissions(groupName) {
		return false
	}

	for _, p := range a.GroupPermissions[groupName] {
		if p == permission {
			return true
		}
	}

	return false
}

// HasGroupPermissions returns true if group has at least one permission.
func (a Album) HasGroupPermissions(groupName string) bool {
	_, found := a.GroupPermissions[groupName]

	return found
}

// GetGroupPermissions returns all permission of group.
func (a Album) GetGroupPermissions(groupName string) (permissions []Permission, found bool) {
	if _, found = a.GroupPermissions[groupName]; !found {
		return
	}

	return a.GroupPermissions[groupName], true
}

type AlbumLessFunc func(a1, a2 Album) bool

type AlbumSorter struct {
	Album    []Album
	LessFunc AlbumLessFunc
}

func NewAlbumSorter(albums []Album, lessFunc AlbumLessFunc) *AlbumSorter {
	return &AlbumSorter{albums, lessFunc}
}

func (as *AlbumSorter) Len() int {
	return len(as.Album)
}

func (as *AlbumSorter) Swap(i, j int) {
	as.Album[i], as.Album[j] = as.Album[j], as.Album[i]
}

func (as *AlbumSorter) Less(i, j int) bool {
	return as.LessFunc(as.Album[i], as.Album[j])
}
