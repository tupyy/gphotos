package utils

import "github.com/tupyy/gophoto/internal/domain/entity"

// HasUserPermissions returns true if user has at least one permission.
func HasUserPermissions(a entity.Album, userID string) bool {
	_, found := a.UserPermissions[userID]

	return found
}

// HasUserPermission returns true if user has permission set for the album.
func HasUserPermission(a entity.Album, userID string, permission entity.Permission) bool {
	if !HasUserPermissions(a, userID) {
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
func GetUserPermissions(a entity.Album, userID string) (permissions []entity.Permission, found bool) {
	if _, found = a.UserPermissions[userID]; !found {
		return
	}

	return a.UserPermissions[userID], true
}

// HasGroupPermission returns true if group has permission set for album.
func HasGroupPermission(a entity.Album, groupName string, permission entity.Permission) bool {
	if !HasGroupPermissions(a, groupName) {
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
func HasGroupPermissions(a entity.Album, groupName string) bool {
	_, found := a.GroupPermissions[groupName]

	return found
}

// GetGroupPermissions returns all permission of group.
func GetGroupPermissions(a entity.Album, groupName string) (permissions []entity.Permission, found bool) {
	if _, found = a.GroupPermissions[groupName]; !found {
		return
	}

	return a.GroupPermissions[groupName], true
}
