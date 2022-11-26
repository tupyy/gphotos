package entity

// HasUserPermissions returns true if user has at least one permission.
func HasUserPermissions(a Album, userID string) bool {
	found := false
	for _, perm := range a.UserPermissions {
		if perm.OwnerID == userID {
			return true
		}
	}
	return found
}

// HasUserPermission returns true if user has permission set for the album.
func HasUserPermission(a Album, userID string, permission Permission) bool {
	if !HasUserPermissions(a, userID) {
		return false
	}

	perms, _ := GetUserPermissions(a, userID)
	for _, p := range perms {
		if p == permission {
			return true
		}
	}

	return false
}

// GetUserPermissions returns all permission of user.
func GetUserPermissions(a Album, userID string) (permissions []Permission, found bool) {
	for _, perm := range a.UserPermissions {
		if perm.OwnerID == userID {
			return perm.Permissions, true
		}
	}

	return []Permission{}, false
}

// HasGroupPermission returns true if group has permission set for album.
func HasGroupPermission(a Album, groupName string, permission Permission) bool {
	if !HasGroupPermissions(a, groupName) {
		return false
	}

	perms, _ := GetGroupPermissions(a, groupName)
	for _, p := range perms {
		if p == permission {
			return true
		}
	}

	return false
}

// HasGroupPermissions returns true if group has at least one permission.
func HasGroupPermissions(a Album, groupName string) bool {
	found := false
	for _, perm := range a.GroupPermissions {
		if perm.OwnerID == groupName {
			return true
		}
	}
	return found
}

// GetGroupPermissions returns all permission of group.
func GetGroupPermissions(a Album, groupName string) (permissions []Permission, found bool) {
	for _, perm := range a.GroupPermissions {
		if perm.OwnerID == groupName {
			return perm.Permissions, true
		}
	}

	return []Permission{}, false
}
