package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Permission int

const (
	// PermissionReadAlbum gives the user the right to view the album.
	PermissionReadAlbum Permission = iota
	// PermissionWriteAlbum gives the user the right to update photos.
	PermissionWriteAlbum
	// PermissionEditAlbum gives the user the right to edit album informations.
	PermissionEditAlbum
	// PermissionDeleteAlbum gives the user the right to delete the album.
	PermissionDeleteAlbum
	// Permission unknown
	PermissionUnknown
)

var ErrInvalidPermission = errors.New("invalid permission")

func (p Permission) String() string {
	switch p {
	case PermissionReadAlbum:
		return "album.read"
	case PermissionWriteAlbum:
		return "album.write"
	case PermissionEditAlbum:
		return "album.edit"
	case PermissionDeleteAlbum:
		return "album.delete"
	}

	return "unknown"
}

func NewPermission(perm string) (Permission, error) {
	switch perm {
	case "album.read":
		return PermissionReadAlbum, nil
	case "album.write":
		return PermissionWriteAlbum, nil
	case "album.edit":
		return PermissionEditAlbum, nil
	case "album.delete":
		return PermissionDeleteAlbum, nil
	default:
		return PermissionUnknown, ErrInvalidPermission
	}
}

type Album struct {
	ID          int32
	Name        string    `validate:"required"`
	CreatedAt   time.Time `validate:"required"`
	OwnerID     string    `validate:"required"`
	Description *string
	Location    *string
	// UserPermissions holds the list of permissions of other users for this album.
	// The key is the user id.
	UserPermissions map[string][]Permission
	// GroupPermissions holds the list of permissions of groups for this album.
	// The key is the group name.
	GroupPermissions map[string][]Permission
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

	if a.Description != nil {
		fmt.Fprintf(&sb, "description = %s\n", *a.Description)
	}

	if a.Location != nil {
		fmt.Fprintf(&sb, "location = %s\n", *a.Location)
	}

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
func (a Album) HasGroupPermissions(groupID string) bool {
	_, found := a.GroupPermissions[groupID]

	return found
}

// GetGroupPermissions returns all permission of group.
func (a Album) GetGroupPermissions(groupID string) (permissions []Permission, found bool) {
	if _, found = a.GroupPermissions[groupID]; !found {
		return
	}

	return a.GroupPermissions[groupID], true
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
