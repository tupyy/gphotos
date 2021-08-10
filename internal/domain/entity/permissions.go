package entity

import (
	"encoding/json"
	"errors"
	"html"
	"strings"

	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
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

// Permissions represents a mapping between an user/group and a set of permissions
type Permissions map[string][]Permission

func (pp Permissions) Parse(permString map[string][]string, encrypted bool) {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	for k, v := range permString {
		if len(k) == 0 {
			continue
		}

		var id string
		if encrypted {
			var err error
			decryptedID, err := gen.DecryptData(k)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("key", k).Error("decrypt id")

				continue
			}

			id = decryptedID
		} else {
			// sanitize key
			id = html.EscapeString(k)
		}

		entities := make([]Permission, 0, len(v))

		for i := range v {
			switch string(v[i]) {
			case "r":
				entities = append(entities, PermissionReadAlbum)
			case "w":
				entities = append(entities, PermissionWriteAlbum)
			case "e":
				entities = append(entities, PermissionEditAlbum)
			case "d":
				entities = append(entities, PermissionDeleteAlbum)
			}
		}

		if len(entities) > 0 {
			pp[id] = entities
		}
	}
}

func (pp Permissions) Json() (string, error) {
	permForm := make(map[string][]string)

	for k, v := range pp {
		vv := make([]string, 0, len(v))

		for _, permission := range v {
			parts := strings.Split(permission.String(), ".")
			vv = append(vv, strings.ToLower(string(parts[1][0])))
		}

		if len(vv) > 0 {
			permForm[k] = vv
		}
	}

	j, err := json.Marshal(permForm)
	if err != nil {
		return "", err
	}

	return string(j), nil
}
