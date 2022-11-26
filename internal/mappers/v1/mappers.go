package v1

import (
	"fmt"
	"strings"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/encryption"
)

const (
	baseV1URL = "/api/gphotos/v1"
)

func MapAlbumToModel(album entity.Album) apiv1.Album {
	encryption, _ := encryption.New() // must not fail here
	encryptedUsername, _ := encryption.Encrypt(album.Owner)

	albumRef := mapAlbumRef(album)
	tags := make([]apiv1.Tag, 0, len(album.Tags))
	for _, tag := range album.Tags {
		tags = append(tags, MapTagToModel(tag))
	}

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
		Owner: apiv1.ObjectReference{
			Kind: UserKind,
			Href: fmt.Sprintf("%s/users/%s", baseV1URL, encryptedUsername),
			Id:   encryptedUsername,
		},
		Permissions: apiv1.ObjectReference{
			Kind: AlbumPermissionsKind,
			Href: fmt.Sprintf("%s/permissions", albumRef.Href),
			Id:   albumRef.Id,
		},
		Photos: apiv1.ObjectReference{
			Kind: PhotoListKind,
			Href: fmt.Sprintf("%s/photos", albumRef.Href),
			Id:   albumRef.Id,
		},
		Tags: &tags,
	}

	return model
}

func mapAlbumRef(album entity.Album) apiv1.ObjectReference {
	encryption, _ := encryption.New() // must not fail here
	encryptedID, _ := encryption.Encrypt(album.ID)

	return apiv1.ObjectReference{
		Href: fmt.Sprintf("%s/album/%s", baseV1URL, encryptedID),
		Id:   encryptedID,
		Kind: AlbumKind,
	}
}

func MapAlbumPermissions(album entity.Album) apiv1.AlbumPermissions {
	encryption, _ := encryption.New() // must not fail here. TODO find a better way

	mapPermissions := func(permissions []entity.AlbumPermission, kind string) []apiv1.Permissions {
		apiPermissions := []apiv1.Permissions{}
		for _, perms := range permissions {
			encryptedID, _ := encryption.Encrypt(perms.OwnerID)
			up := apiv1.Permissions{
				Owner: apiv1.ObjectReference{
					Kind: kind,
					Href: fmt.Sprintf("%s/%s/%s", baseV1URL, strings.ToLower(kind), encryptedID),
					Id:   encryptedID,
				},
			}
			for _, permission := range perms.Permissions {
				up.Permissions = append(up.Permissions, permission.String())
			}
			apiPermissions = append(apiPermissions, up)
		}
		return apiPermissions
	}

	userPermissions := mapPermissions(album.UserPermissions, UserKind)
	groupPermissions := mapPermissions(album.GroupPermissions, GroupKind)
	albumRef := mapAlbumRef(album)

	albumPermissions := apiv1.AlbumPermissions{
		Kind:   AlbumPermissionsKind,
		Id:     albumRef.Id,
		Href:   fmt.Sprintf("%s/album/%s/permissions", baseV1URL, albumRef.Id),
		Users:  &userPermissions,
		Groups: &groupPermissions,
		Album:  &albumRef,
	}

	return albumPermissions
}
