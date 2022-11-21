package v1

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

func MapMediaToModel(album entity.Album, photo entity.Media) apiv1.Photo {
	id := fmt.Sprintf("%s/%s", photo.Bucket, photo.Filename)
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, _ := gen.EncryptData(id)
	model := apiv1.Photo{
		Album:     mapAlbumRef(album),
		Id:        encryptedID,
		Href:      fmt.Sprintf("%s/photo/%s", baseV1URL, encryptedID),
		Kind:      PhotoKind,
		Thumbnail: fmt.Sprintf("%s/photo/%s/thumbnail", baseV1URL, encryptedID),
	}
	return model
}

func MapMediaListToModel(album entity.Album, photos []entity.Media) apiv1.PhotoList {
	model := apiv1.PhotoList{
		Items: make([]apiv1.Photo, 0, len(photos)),
		Kind:  PhotoListKind,
	}
	for _, photo := range photos {
		model.Items = append(model.Items, MapMediaToModel(album, photo))
	}
	return model
}
