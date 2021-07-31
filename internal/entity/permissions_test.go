package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePermissions(t *testing.T) {
	testData := []struct {
		PermString    string
		CapturedPerms map[string][]Permission
	}{
		{
			PermString: "(bob#rw)(jane#ed)(joe#rwed)",
			CapturedPerms: map[string][]Permission{
				"bob":  {PermissionReadAlbum, PermissionWriteAlbum},
				"jane": {PermissionEditAlbum, PermissionDeleteAlbum},
				"joe":  {PermissionReadAlbum, PermissionWriteAlbum, PermissionEditAlbum, PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#rw)(jane#)(joe#rwed)",
			CapturedPerms: map[string][]Permission{
				"bob": {PermissionReadAlbum, PermissionWriteAlbum},
				"joe": {PermissionReadAlbum, PermissionWriteAlbum, PermissionEditAlbum, PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#rw)(jane#)(joe#rwed)",
			CapturedPerms: map[string][]Permission{
				"bob": {PermissionReadAlbum, PermissionWriteAlbum},
				"joe": {PermissionReadAlbum, PermissionWriteAlbum, PermissionEditAlbum, PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#rw)jane#re",
			CapturedPerms: map[string][]Permission{
				"bob": {PermissionReadAlbum, PermissionWriteAlbum},
			},
		},
	}

	var perms Permissions
	for _, td := range testData {
		perms.Decode(td.PermString, false)

		assert.EqualValues(t, td.CapturedPerms, perms)
	}
}
