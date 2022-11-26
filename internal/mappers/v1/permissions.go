package v1

import (
	"fmt"
	"strings"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/encryption"
)

func MapToEntityPermissions(form apiv1.AlbumPermissionsRequest) ([]entity.AlbumPermission, error) {
	encryption, _ := encryption.New() // must not fail here
	mapToPermissionList := func(perms []string) []entity.Permission {
		pperms := make([]entity.Permission, 0, len(perms))
		for _, pp := range perms {
			perm, err := entity.NewPermission(pp)
			if err == nil {
				pperms = append(pperms, perm)
			}
		}
		return pperms
	}

	albumPermissions := []entity.AlbumPermission{}
	for _, p := range form {
		id, err := encryption.Decrypt(p.Owner.Id)
		if err != nil {
			id = p.Owner.Id
		}
		perms := mapToPermissionList(p.Permissions)
		if strings.ToLower(p.Owner.Kind) != "user" && strings.ToLower(p.Owner.Kind) != "group" {
			return []entity.AlbumPermission{}, fmt.Errorf("invalid error kind: '%s'", p.Owner.Kind)
		}
		albumPermissions = append(albumPermissions, entity.AlbumPermission{
			OwnerID:     id,
			OwnerKind:   strings.ToLower(p.Owner.Kind),
			Permissions: perms,
		})
	}
	return albumPermissions, nil
}
