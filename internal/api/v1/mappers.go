package v1

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

func mapAlbumToModel(album entity.Album, users []entity.User) (apiv1.Album, error) {
	owner := entity.User{}
	for _, user := range users {
		if user.ID == album.OwnerID {
			owner = user
			break
		}
	}

	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", album.ID))
	if err != nil {
		return apiv1.Album{}, err
	}

	encryptedUsername, err := gen.EncryptData(owner.Username)
	if err != nil {
		return apiv1.Album{}, err
	}

	model := apiv1.Album{
		Id:          encryptedID,
		Href:        fmt.Sprintf("/api/v1/albums/%s", encryptedID),
		Kind:        "UserAlbum",
		Bucket:      album.Bucket,
		Name:        album.Name,
		Description: &album.Description,
		Location:    &album.Location,
		CreatedAt:   album.CreatedAt.Unix(),
		Thumbnail:   &album.Thumbnail,
		Owner: &apiv1.ObjectReference{
			Kind: "User",
			Href: fmt.Sprintf("/api/v1/users/%s", encryptedUsername),
			Id:   encryptedUsername,
		},
	}

	return model, nil
}
