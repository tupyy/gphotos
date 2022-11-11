package album

import (
	"database/sql"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repos/models"
)

func toModel(e entity.Album) models.Album {
	m := models.Album{
		Name:        e.Name,
		CreatedAt:   e.CreatedAt,
		OwnerID:     e.Owner,
		Description: &e.Description,
		Location:    &e.Location,
		Bucket:      e.Bucket,
	}

	if len(e.Thumbnail) == 0 {
		m.Thumbnail = sql.NullString{Valid: false}
	} else {
		m.Thumbnail = sql.NullString{String: e.Thumbnail, Valid: true}
	}

	return m
}

func fromModel(m models.Album) entity.Album {
	e := entity.Album{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Owner:     m.OwnerID,
		Bucket:    m.Bucket,
	}

	if m.Description != nil {
		e.Description = *m.Description
	}

	if m.Location != nil {
		e.Location = *m.Location
	}

	if m.Thumbnail.Valid {
		e.Thumbnail = m.Thumbnail.String
	}

	return e
}
