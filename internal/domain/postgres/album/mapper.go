package album

import (
	"database/sql"

	"github.com/tupyy/gophoto/internal/domain/models"
	"github.com/tupyy/gophoto/internal/entity"
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

func toUserPermissionsModels(albumID int32, permissions map[string][]entity.Permission) []models.AlbumUserPermissions {
	if len(permissions) == 0 {
		return []models.AlbumUserPermissions{}
	}

	permModels := make([]models.AlbumUserPermissions, 0)

	for k, v := range permissions {
		mm := make(models.PermissionIDs, 0, len(v))

		for _, p := range v {
			mm = append(mm, models.PermissionID(p.String()))
		}

		permModel := models.AlbumUserPermissions{
			UserID:      k,
			AlbumID:     albumID,
			Permissions: mm,
		}

		permModels = append(permModels, permModel)
	}

	return permModels
}

func toGroupPermissionsModels(albumID int32, permissions map[string][]entity.Permission) []models.AlbumGroupPermissions {
	if len(permissions) == 0 {
		return []models.AlbumGroupPermissions{}
	}

	permModels := make([]models.AlbumGroupPermissions, 0)

	for k, v := range permissions {
		mm := make(models.PermissionIDs, 0, len(v))

		for _, p := range v {
			mm = append(mm, models.PermissionID(p.String()))
		}

		permModel := models.AlbumGroupPermissions{
			GroupName:   k,
			AlbumID:     albumID,
			Permissions: mm,
		}

		permModels = append(permModels, permModel)
	}

	return permModels
}
