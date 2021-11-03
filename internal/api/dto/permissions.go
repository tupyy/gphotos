package dto

import "github.com/tupyy/gophoto/internal/entity"

type PermissionDTO struct {
	UserPermissions  map[string][]string `json:"user_permissions"`
	GroupPermissions map[string][]string `json:"group_permissions"`
}

func NewPermissionDTO(album entity.Album, users []User) PermissionDTO {
	permissions := PermissionDTO{
		UserPermissions:  make(map[string][]string),
		GroupPermissions: make(map[string][]string),
	}

	for id, p := range album.UserPermissions {
		var userEncryptedID string

		for _, user := range users {
			if user.ID == id {
				userEncryptedID = user.EncryptedID
			}
		}

		if len(userEncryptedID) == 0 {
			continue
		}

		perms := []string{}
		for _, pp := range p {
			perms = append(perms, toString(pp))
		}

		permissions.UserPermissions[userEncryptedID] = perms
	}

	for name, p := range album.GroupPermissions {
		perms := []string{}
		for _, pp := range p {
			perms = append(perms, toString(pp))
		}

		permissions.GroupPermissions[name] = perms
	}

	return permissions
}

func toString(p entity.Permission) string {
	switch p {
	case entity.PermissionReadAlbum:
		return "r"
	case entity.PermissionWriteAlbum:
		return "w"
	case entity.PermissionEditAlbum:
		return "e"
	case entity.PermissionDeleteAlbum:
		return "d"
	default:
		return "unknown"
	}
}
