package dto

import (
	"fmt"

	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/encryption"
)

type Tag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func NewTagDTO(tag entity.Tag) (Tag, error) {
	dto := Tag{}

	if tag.Color != nil {
		dto.Color = *tag.Color
	}

	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	// encrypt tag id
	id, err := gen.EncryptData(fmt.Sprintf("%d", tag.ID))
	if err != nil {
		return dto, fmt.Errorf("encrypt tag id '%d': %+v", tag.ID, err)
	}

	dto.ID = id
	dto.Name = tag.Name

	return dto, nil
}
