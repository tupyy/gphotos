package v1

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

const (
	baseV1URL = "/api/gphotos/v1"
)

func MapAlbumToModel(album entity.Album) apiv1.Album {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedUsername, _ := gen.EncryptData(album.Owner)

	albumRef := mapAlbumRef(album)
	permissions := MapAlbumPermissions(album)

	model := apiv1.Album{
		Id:          albumRef.Id,
		Href:        albumRef.Href,
		Kind:        albumRef.Kind,
		Bucket:      album.Bucket,
		Name:        album.Name,
		Description: &album.Description,
		Location:    &album.Location,
		CreatedAt:   album.CreatedAt,
		Thumbnail:   &album.Thumbnail,
		Owner: &apiv1.ObjectReference{
			Kind: "User",
			Href: fmt.Sprintf("%s/users/%s", baseV1URL, encryptedUsername),
			Id:   encryptedUsername,
		},
		Permissions: &permissions,
	}

	return model
}

func mapAlbumRef(album entity.Album) apiv1.ObjectReference {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, _ := gen.EncryptData(album.ID)

	return apiv1.ObjectReference{
		Href: fmt.Sprintf("%s/album/%s", baseV1URL, encryptedID),
		Id:   encryptedID,
		Kind: "Album",
	}
}

func MapAlbumPermissions(album entity.Album) apiv1.AlbumPermissions {
	mapPermissions := func(permissions entity.Permissions, kind string) []apiv1.Permissions {
		apiPermissions := []apiv1.Permissions{}
		for owner, perms := range permissions {
			up := apiv1.Permissions{
				Owner: apiv1.ObjectReference{
					Kind: kind,
					Href: fmt.Sprintf("%s/%s/%s", baseV1URL, kind, owner),
					Id:   owner,
				},
			}
			for _, permission := range perms {
				up.Permissions = append(up.Permissions, permission.String())
			}
			apiPermissions = append(apiPermissions, up)
		}
		return apiPermissions
	}

	userPermissions := mapPermissions(album.UserPermissions, "user")
	groupPermissions := mapPermissions(album.GroupPermissions, "group")
	albumRef := mapAlbumRef(album)

	albumPermissions := apiv1.AlbumPermissions{
		Kind:   "AlbumPermissionsList",
		Id:     albumRef.Id,
		Href:   fmt.Sprintf("%s/album/%s/permissions", baseV1URL, albumRef.Id),
		Users:  &userPermissions,
		Groups: &groupPermissions,
		Album:  &albumRef,
	}

	return albumPermissions
}
