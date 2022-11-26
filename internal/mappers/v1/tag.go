package v1

import (
	"fmt"

	apiv1 "github.com/tupyy/gophoto/api/v1"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/encryption"
)

func MapTagToModel(tag entity.Tag) apiv1.Tag {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, _ := gen.EncryptData(tag.ID)
	encryptedOwner, _ := gen.EncryptData(tag.UserID)

	model := apiv1.Tag{
		Id:    encryptedID,
		Href:  fmt.Sprintf("%s/%s", baseV1URL, encryptedID),
		Kind:  TagKind,
		Name:  tag.Name,
		Color: tag.Color,
		User: apiv1.ObjectReference{
			Kind: UserKind,
			Href: fmt.Sprintf("%s/users/%s", baseV1URL, encryptedOwner),
			Id:   encryptedOwner,
		},
		Albums: make([]apiv1.ObjectReference, 0, len(tag.Albums)),
	}

	return model
}

func MapTagsToList(tags []entity.Tag) apiv1.TagList {
	list := apiv1.TagList{
		Kind:  TagListKind,
		Page:  1,
		Size:  len(tags),
		Total: len(tags),
		Items: make([]apiv1.Tag, 0, len(tags)),
	}

	for _, tag := range tags {
		list.Items = append(list.Items, MapTagToModel(tag))
	}

	return list
}
