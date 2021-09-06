package entity

// HasUserPermissions returns true if user has at least one permission.
func HasUserPermissions(a Album, userID string) bool {
	_, found := a.UserPermissions[userID]

	return found
}

// HasUserPermission returns true if user has permission set for the album.
func HasUserPermission(a Album, userID string, permission Permission) bool {
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
func GetUserPermissions(a Album, userID string) (permissions []Permission, found bool) {
	if _, found = a.UserPermissions[userID]; !found {
		return
	}

	return a.UserPermissions[userID], true
}

// HasGroupPermission returns true if group has permission set for album.
func HasGroupPermission(a Album, groupName string, permission Permission) bool {
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
func HasGroupPermissions(a Album, groupName string) bool {
	_, found := a.GroupPermissions[groupName]

	return found
}

// GetGroupPermissions returns all permission of group.
func GetGroupPermissions(a Album, groupName string) (permissions []Permission, found bool) {
	if _, found = a.GroupPermissions[groupName]; !found {
		return
	}

	return a.GroupPermissions[groupName], true
}
