package entity

import (
	"errors"
)

type Permission int

const (
	// PermissionReadAlbum gives the user the right to view the album.
	PermissionReadAlbum Permission = iota
	// PermissionWriteAlbum gives the user the right to update photos.
	PermissionWriteAlbum
	// PermissionEditAlbum gives the user the right to edit album information.
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

type AlbumPermission struct {
	OwnerID     string
	OwnerKind   string
	Permissions []Permission
}
