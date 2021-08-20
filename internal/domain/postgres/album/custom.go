package album

import (
	"fmt"
	"time"

	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/models"
	"github.com/tupyy/gophoto/utils/logutil"
)

// custom struct to map the join
type customAlbum struct {
	ID               int32                `gorm:"column_name:id;type:INT4"`
	Name             string               `gorm:"column:name;type:TEXT;"`
	CreatedAt        time.Time            `gorm:"column:created_at;type:TIMESTAMP;default:timezone('UTC');"`
	OwnerID          string               `gorm:"column:owner_id;type:INT4;"`
	Description      *string              `gorm:"column:description;type:TEXT;"`
	Location         *string              `gorm:"column:location;type:TEXT;"`
	UserPermissions  models.PermissionIDs `gorm:"column:user_permissions;type:_PERMISSION_ID;"`
	GroupPermissions models.PermissionIDs `gorm:"column:group_permissions;type:_PERMISSION_ID;"`
	UserID           string               `gorm:"column:user_id;type:TEXT;"`
	GroupName        string               `gorm:"column:group_name;type:TEXT;"`
}

func (ca customAlbum) ToEntity() (entity.Album, error) {
	var emptyAlbum entity.Album

	album := entity.Album{
		ID:        ca.ID,
		Name:      ca.Name,
		CreatedAt: ca.CreatedAt,
		OwnerID:   ca.OwnerID,
	}

	if ca.Description != nil {
		album.Description = *ca.Description
	}

	if ca.Location != nil {
		album.Location = *ca.Location
	}

	if len(ca.UserPermissions) > 0 {
		album.UserPermissions = make(map[string][]entity.Permission)

		permissions := make([]entity.Permission, 0, len(ca.UserPermissions))

		for _, perm := range ca.UserPermissions {
			if p, err := entity.NewPermission(string(perm)); err != nil {
				logutil.GetDefaultLogger().WithField("permission", perm).Warn("error parsing permission")

				return emptyAlbum, fmt.Errorf("%w cannot parse permission", err)
			} else {
				permissions = append(permissions, p)
			}
		}

		album.UserPermissions[ca.UserID] = permissions
	}

	if len(ca.GroupPermissions) > 0 {
		album.GroupPermissions = make(map[string][]entity.Permission)

		permissions := make([]entity.Permission, 0, len(ca.GroupPermissions))

		for _, perm := range ca.GroupPermissions {
			if p, err := entity.NewPermission(string(perm)); err != nil {
				logutil.GetDefaultLogger().WithField("permission", perm).Warn("error parsing permission")

				return emptyAlbum, fmt.Errorf("%w cannot parse permission", err)
			} else {
				permissions = append(permissions, p)
			}
		}

		album.GroupPermissions[ca.GroupName] = permissions
	}

	return album, nil
}

type customAlbums []customAlbum

// mergeAlbums merge a list of customAlbums into a list of distinct albums.
func (albums customAlbums) Merge() []entity.Album {
	entitiesMap := make(map[int32]entity.Album)

	for _, ca := range albums {
		if e, err := ca.ToEntity(); err != nil {
			logutil.GetDefaultLogger().WithError(err).Warn("cannot create entity")
		} else {
			if ent, found := entitiesMap[e.ID]; found {
				// merge permissions
				if len(e.UserPermissions) > 0 {
					for k, v := range e.UserPermissions {
						ent.UserPermissions[k] = v
					}
				}

				if len(e.GroupPermissions) > 0 {
					for k, v := range e.GroupPermissions {
						ent.GroupPermissions[k] = v
					}
				}
			} else {
				entitiesMap[e.ID] = e
			}
		}
	}

	entities := make([]entity.Album, 0, len(entitiesMap))
	for _, v := range entitiesMap {
		entities = append(entities, v)
	}

	return entities

}
