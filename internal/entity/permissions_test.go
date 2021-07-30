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
			PermString: "(bob#r,w)(jane#e,d)(joe#r,w,e,d)",
			CapturedPerms: map[string][]Permission{
				"bob":  {PermissionReadAlbum, PermissionWriteAlbum},
				"jane": {PermissionEditAlbum, PermissionDeleteAlbum},
				"joe":  {PermissionReadAlbum, PermissionWriteAlbum, PermissionEditAlbum, PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#r,w)(jane#)(joe#r,w,e,d)",
			CapturedPerms: map[string][]Permission{
				"bob": {PermissionReadAlbum, PermissionWriteAlbum},
				"joe": {PermissionReadAlbum, PermissionWriteAlbum, PermissionEditAlbum, PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#rw)(jane#)(joe#r,w,e,d)",
			CapturedPerms: map[string][]Permission{
				"joe": {PermissionReadAlbum, PermissionWriteAlbum, PermissionEditAlbum, PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#r,w)jane#r,e",
			CapturedPerms: map[string][]Permission{
				"bob": {PermissionReadAlbum, PermissionWriteAlbum},
			},
		},
	}

	var perms Permissions
	for _, td := range testData {
		p := perms.Decode(td.PermString, false)

		assert.EqualValues(t, td.CapturedPerms, p)
	}
}
