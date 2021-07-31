package entity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
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

// Encode will return a string like "(username#rw)(username2#ed)" to be used in forms
// If encrypted is set to true username is encrypted with server encryption key.
func (p Permissions) Encode(encrypted bool) string {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	permStr := ""
	for key, permissions := range p {
		encryptedKey := key
		if encrypted {
			var err error

			encryptedKey, err = gen.EncryptData(key)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).Error("cannot ecrypt name")

				continue
			}
		}

		permList := ""
		for _, p := range permissions {
			str := strings.Split(p.String(), ".")
			permList = permList + strings.ToLower(string(str[1][0]))
		}

		permStr = permStr + fmt.Sprintf("(%s#%s)", encryptedKey, permList)
	}

	return permStr
}

func (p Permissions) Decode(encodedPermissions string, encrypted bool) {
	permRe := regexp.MustCompile(`(\((\w+)#(([rwed])+)\))`)
	permissions := make(map[string][]Permission)

	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	for matchIdx, match := range permRe.FindAllStringSubmatch(encodedPermissions, -1) {
		logutil.GetDefaultLogger().WithFields(logrus.Fields{"idx": matchIdx, "match": fmt.Sprintf("%+v", match)}).Debug("permission matched")
		// get 2nd and 3rd groups only
		name := match[2]

		if encrypted {
			var err error

			name, err = gen.DecryptData(match[2])
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("data", match[2]).Error("decrypt name")

				continue
			}
		}

		permList := match[3]
		entities := make([]Permission, 0, 4)

		for i := range permList {
			switch string(permList[i]) {
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
			permissions[name] = entities
		}
	}

	p = permissions
}
