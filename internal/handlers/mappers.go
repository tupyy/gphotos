package handlers

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

func mapAlbumToModel(album entity.Album, users []entity.User) (apiv1.UserAlbum, error) {
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
		return apiv1.UserAlbum{}, err
	}

	encryptedUsername, err := gen.EncryptData(owner.Username)
	if err != nil {
		return apiv1.UserAlbum{}, err
	}

	model := apiv1.UserAlbum{
		Id:          encryptedID,
		Href:        fmt.Sprintf("/api/v1/albums/%s", encryptedID),
		Kind:        "UserAlbum",
		Bucket:      album.Bucket,
		Name:        album.Name,
		Description: &album.Description,
		Location:    &album.Location,
		CreatedAt:   album.CreatedAt.Unix(),
		Thumbnail:   &album.Thumbnail,
		Owner: &apiv1.User{
			Kind:    "User",
			Href:    fmt.Sprintf("/api/v1/users/%s", encryptedUsername),
			UserId:  &encryptedUsername,
			Id:      encryptedID,
			Name:    &owner.FirstName,
			Surname: &owner.LastName,
		},
	}

	return model, nil
}
