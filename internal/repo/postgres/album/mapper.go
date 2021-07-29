package album

import (
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/models"
)

func toModel(e entity.Album) models.Album {
	return models.Album{
		Name:        e.Name,
		CreatedAt:   e.CreatedAt,
		OwnerID:     e.OwnerID,
		Description: &e.Description,
		Location:    &e.Location,
	}
}

func fromModel(m models.Album) entity.Album {
	e := entity.Album{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		OwnerID:   m.OwnerID,
	}

	if m.Description != nil {
		e.Description = *m.Description
	}

	if m.Location != nil {
		e.Location = *m.Location
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
