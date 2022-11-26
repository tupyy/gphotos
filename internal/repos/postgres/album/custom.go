package album

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repos/models"
	"go.uber.org/zap"
)

// custom struct to map the join
type albumJoinRow struct {
	ID                  string               `gorm:"column_name:id;type:TEXT"`
	Name                string               `gorm:"column:name;type:TEXT;"`
	CreatedAt           time.Time            `gorm:"column:created_at;type:TIMESTAMP;default:timezone('UTC');"`
	OwnerID             string               `gorm:"column:owner_id;type:TEXT;"`
	Description         *string              `gorm:"column:description;type:TEXT;"`
	Location            *string              `gorm:"column:location;type:TEXT;"`
	Bucket              string               `gorm:"column:bucket;type:TEXT;"`
	TagID               string               `gorm:"column:tag_id;type:TEXT"`
	TagName             *string              `gorm:"column:tag_name;type:TEXT;"`
	TagColor            *string              `gorm:"column:tag_color;tape:TEXT"`
	Thumbnail           sql.NullString       `gorm:"column:thumbnail;type:VARCHAR;size:100;"`
	Permissions         models.PermissionIDs `gorm:"column:permissions;type:_PERMISSION_ID;"`
	PermissionOwnerID   string               `gorm:"column:permission_owner_id;type:TEXT;"`
	PermissionOwnerKind string               `gorm:"column:permission_owner_kind;type:TEXT;"`
}

func (ca albumJoinRow) ToEntity() (entity.Album, error) {
	var emptyAlbum entity.Album

	album := entity.Album{
		ID:               ca.ID,
		Name:             ca.Name,
		CreatedAt:        ca.CreatedAt,
		Owner:            ca.OwnerID,
		Bucket:           ca.Bucket,
		UserPermissions:  make([]entity.AlbumPermission, 0),
		GroupPermissions: make([]entity.AlbumPermission, 0),
	}

	if ca.Description != nil {
		album.Description = *ca.Description
	}

	if ca.Location != nil {
		album.Location = *ca.Location
	}

	if ca.Thumbnail.Valid {
		album.Thumbnail = ca.Thumbnail.String
	}

	if len(ca.Permissions) > 0 {
		permissions := []entity.Permission{}
		for _, perm := range ca.Permissions {
			if p, err := entity.NewPermission(string(perm)); err != nil {
				return emptyAlbum, fmt.Errorf("%w cannot parse permission", err)
			} else {
				permissions = append(permissions, p)
			}
		}

		switch ca.PermissionOwnerKind {
		case "user":
			album.UserPermissions = append(album.UserPermissions, entity.AlbumPermission{
				OwnerID:     ca.PermissionOwnerID,
				OwnerKind:   ca.PermissionOwnerKind,
				Permissions: permissions,
			})
		case "group":
			album.GroupPermissions = append(album.GroupPermissions, entity.AlbumPermission{
				OwnerID:     ca.PermissionOwnerID,
				OwnerKind:   ca.PermissionOwnerKind,
				Permissions: permissions,
			})
		default:
			return emptyAlbum, fmt.Errorf("wrong owner kind '%s'", ca.PermissionOwnerKind)
		}

	}

	if ca.TagName != nil {
		album.Tags = make([]entity.Tag, 0, 1)
		if ca.TagColor != nil {
			album.Tags = append(album.Tags, entity.Tag{ID: ca.TagID, Name: *ca.TagName, Color: ca.TagColor})
		} else {
			album.Tags = append(album.Tags, entity.Tag{ID: ca.TagID, Name: *ca.TagName})
		}
	}

	return album, nil
}

type albumJoinRows []albumJoinRow

// mergeAlbums merge a list of customAlbums into a list of distinct albums.
func (albums albumJoinRows) Merge() []entity.Album {
	entitiesMap := make(map[string]entity.Album)

	for _, ca := range albums {
		e, err := ca.ToEntity()
		if err != nil {
			zap.S().Warnw("cannot map to entity", "error", err, "model", ca)
			continue
		}

		if ent, found := entitiesMap[e.ID]; found {
			// merge permissions
			if len(e.UserPermissions) > 0 {
				ent.UserPermissions = append(ent.UserPermissions, e.UserPermissions...)
			}

			if len(e.GroupPermissions) > 0 {
				ent.GroupPermissions = append(ent.GroupPermissions, e.GroupPermissions...)
			}

			// merge tags
			if len(e.Tags) > 0 {
				for _, t := range e.Tags {
					ent.Tags = append(ent.Tags, t)
				}
			}
			delete(entitiesMap, ent.ID)
			entitiesMap[ent.ID] = ent
		} else {
			entitiesMap[e.ID] = e
		}
	}

	entities := make([]entity.Album, 0, len(entitiesMap))
	for _, v := range entitiesMap {
		entities = append(entities, v)
	}

	return entities
}
