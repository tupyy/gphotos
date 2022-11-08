package v1

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

func mapAlbumToModel(album entity.Album) apiv1.Album {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, _ := gen.EncryptData(fmt.Sprintf("%d", album.ID))

	encryptedUsername, _ := gen.EncryptData(album.Owner)

	model := apiv1.Album{
		Id:          encryptedID,
		Href:        fmt.Sprintf("/api/v1/albums/%s", encryptedID),
		Kind:        "Album",
		Bucket:      album.Bucket,
		Name:        album.Name,
		Description: &album.Description,
		Location:    &album.Location,
		CreatedAt:   album.CreatedAt,
		Thumbnail:   &album.Thumbnail,
		Owner: &apiv1.ObjectReference{
			Kind: "User",
			Href: fmt.Sprintf("/api/v1/users/%s", encryptedUsername),
			Id:   encryptedUsername,
		},
	}

	return model
}
