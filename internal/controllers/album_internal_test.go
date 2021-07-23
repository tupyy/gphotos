package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/gophoto/internal/entity"
)

func TestParsePermissions(t *testing.T) {
	testData := []struct {
		PermString    string
		CapturedPerms map[string][]entity.Permission
	}{
		{
			PermString: "(bob#r,w)(jane#e,d)(joe#r,w,e,d)",
			CapturedPerms: map[string][]entity.Permission{
				"bob":  {entity.PermissionReadAlbum, entity.PermissionWriteAlbum},
				"jane": {entity.PermissionEditAlbum, entity.PermissionDeleteAlbum},
				"joe":  {entity.PermissionReadAlbum, entity.PermissionWriteAlbum, entity.PermissionEditAlbum, entity.PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#r,w)(jane#)(joe#r,w,e,d)",
			CapturedPerms: map[string][]entity.Permission{
				"bob": {entity.PermissionReadAlbum, entity.PermissionWriteAlbum},
				"joe": {entity.PermissionReadAlbum, entity.PermissionWriteAlbum, entity.PermissionEditAlbum, entity.PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#rw)(jane#)(joe#r,w,e,d)",
			CapturedPerms: map[string][]entity.Permission{
				"joe": {entity.PermissionReadAlbum, entity.PermissionWriteAlbum, entity.PermissionEditAlbum, entity.PermissionDeleteAlbum},
			},
		},
		{
			PermString: "(bob#r,w)jane#r,e",
			CapturedPerms: map[string][]entity.Permission{
				"bob": {entity.PermissionReadAlbum, entity.PermissionWriteAlbum},
			},
		},
	}

	for _, td := range testData {
		p := parsePermissions(td.PermString)

		assert.EqualValues(t, td.CapturedPerms, p)
	}
}
