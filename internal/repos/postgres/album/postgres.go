package album

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	pgclient "github.com/tupyy/gophoto/internal/clients/pg"
	"github.com/tupyy/gophoto/internal/common"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repos/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AlbumPostgresRepo struct {
	db             *gorm.DB
	client         pgclient.Client
	circuitBreaker pgclient.CircuitBreaker
}

func NewPostgresRepo(client pgclient.Client) (*AlbumPostgresRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &AlbumPostgresRepo{}, err
	}

	return &AlbumPostgresRepo{gormDB, client, client.GetCircuitBreaker()}, nil
}

func (a *AlbumPostgresRepo) Create(ctx context.Context, album entity.Album) (entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return entity.Album{}, common.NewPostgresNotAvailableError("pg not available while creating album")
	}

	m := toModel(album)
	m.ID = xid.New().String()
	album.ID = m.ID

	result := a.db.WithContext(ctx).Create(&m)
	if result.Error != nil {
		if a.checkNetworkError(result.Error) {
			return entity.Album{}, common.NewPostgresNotAvailableError("pg not available while creating album")
		}
		return album, common.NewInternalError(result.Error, "failed to create album")
	}

	return album, nil
}

func (a *AlbumPostgresRepo) Delete(ctx context.Context, id string) error {
	if !a.circuitBreaker.IsAvailable() {
		return common.NewPostgresNotAvailableError("pg not available while removing album")
	}

	if result := a.db.WithContext(ctx).Delete(&models.Album{}, id); result.Error != nil {
		if a.checkNetworkError(result.Error) {
			return common.NewPostgresNotAvailableError("pg not available while removing album")
		}
		return common.NewInternalError(result.Error, fmt.Sprintf("failed to delete album with id '%s'", id))
	}
	return nil
}

func (a *AlbumPostgresRepo) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return album, common.NewPostgresNotAvailableError("pg not available while updating album")
	}

	var ca albumJoinRow
	tx := a.db.WithContext(ctx).Table("album").Where("id = ?", album.ID).First(&ca)
	if tx.Error != nil {
		if a.checkNetworkError(tx.Error) {
			return album, common.NewPostgresNotAvailableError("pg not available while updating album")
		}
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return album, common.NewEntityNotFound(fmt.Sprintf("album '%s' not found", album.ID))
		}
		return album, common.NewInternalError(tx.Error, fmt.Sprintf("failed to get album with id '%s'", album.ID))
	}

	newAlbum := entity.Album{
		Name:        album.Name,
		CreatedAt:   album.CreatedAt,
		Description: album.Description,
		Location:    album.Location,
		Owner:       album.Owner,
		Bucket:      album.Bucket,
		Thumbnail:   album.Thumbnail,
	}

	m := toModel(newAlbum)
	m.ID = album.ID

	result := a.db.WithContext(ctx).Save(&m)
	if result.Error != nil {
		if a.checkNetworkError(tx.Error) {
			return album, common.NewPostgresNotAvailableError("pg not available while updating album")
		}
		return album, common.NewInternalError(result.Error, fmt.Sprintf("failed to update album with id '%s'", album.ID))
	}

	return album, nil
}

// Get returns all the albums sorted by id.
func (a *AlbumPostgresRepo) Get(ctx context.Context) ([]entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while retrieving albums")
	}

	var albums albumJoinRows
	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permission_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery)

	tx.Find(&albums)
	if tx.Error != nil {
		if a.checkNetworkError(tx.Error) {
			return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while retrieving albums")
		}
		return []entity.Album{}, common.NewInternalError(tx.Error, "failed to fetch albums")
	}

	if len(albums) == 0 {
		return []entity.Album{}, common.NewEntityNotFound("albums not found")
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByID return the album if any with id id.
func (a *AlbumPostgresRepo) GetByID(ctx context.Context, id string) (entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return entity.Album{}, common.NewPostgresNotAvailableError("pg not available while retrieving album by id")
	}

	var albums albumJoinRows
	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permission_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album.id = ?", id).
		Find(&albums)

	if tx.Error != nil {
		if !a.circuitBreaker.IsAvailable() {
			return entity.Album{}, common.NewPostgresNotAvailableError("pg not available while retrieving album by id")
		}
		return entity.Album{}, common.NewInternalError(tx.Error, fmt.Sprintf("failed to fetch album by id '%s'", id))
	}

	if len(albums) == 0 {
		return entity.Album{}, common.NewEntityNotFound(fmt.Sprintf("failed to fetch album with id '%s'", id))
	}

	entities := albums.Merge()

	return entities[0], nil
}

// GetByOwnerID return all albums of an user.
func (a *AlbumPostgresRepo) GetByOwner(ctx context.Context, owner string) ([]entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by owner")
	}

	var albums albumJoinRows
	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permission_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album.owner_id = ?", owner)

	tx.Find(&albums)
	if tx.Error != nil {
		if !a.circuitBreaker.IsAvailable() {
			return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by owner")
		}
		return []entity.Album{}, common.NewInternalError(tx.Error, fmt.Sprintf("failed to get albums by owner '%s'", owner))
	}

	if len(albums) == 0 {
		return []entity.Album{}, common.NewEntityNotFound(fmt.Sprintf("albums not found with owner '%s'", owner))
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByUserID returns a list of albums for which the user has at one permission set.
func (a *AlbumPostgresRepo) GetByUser(ctx context.Context, username string) ([]entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by user")
	}

	var albums albumJoinRows
	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permission_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_permissions.owner_kind = ?", "user").
		Where("album_permissions.owner_id = ?", username)

	tx.Find(&albums)
	if tx.Error != nil {
		if !a.circuitBreaker.IsAvailable() {
			return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by user")
		}
		return []entity.Album{}, common.NewInternalError(tx.Error, fmt.Sprintf("failed to get albums by user '%s'", username))
	}

	if len(albums) == 0 {
		return []entity.Album{}, common.NewEntityNotFound(fmt.Sprintf("albums not found with user '%s'", username))
	}

	entities := albums.Merge()

	return entities, nil
}

// GetAlbumsByGroup returns a list of albums for which the group has at one permission set.
func (a *AlbumPostgresRepo) GetByGroupName(ctx context.Context, groupName string) ([]entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by group")
	}

	var albums albumJoinRows
	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permission_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_permissions.owner_kind = ?", "group").
		Where("album_permissions.owner_id = ?", groupName)

	tx.Find(&albums)
	if tx.Error != nil {
		if !a.circuitBreaker.IsAvailable() {
			return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by group")
		}
		return []entity.Album{}, common.NewInternalError(tx.Error, fmt.Sprintf("failed to get albums by group '%s'", groupName))
	}

	if len(albums) == 0 {
		return []entity.Album{}, common.NewEntityNotFound(fmt.Sprintf("albums not found by group '%s'", groupName))
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByGroups returns a list of albums with at least one persmission for at least on group in the list.
func (a *AlbumPostgresRepo) GetByGroups(ctx context.Context, groupNames []string) ([]entity.Album, error) {
	if !a.circuitBreaker.IsAvailable() {
		return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by groups")
	}

	var albums albumJoinRows

	if len(groupNames) == 0 {
		return []entity.Album{}, nil
	}

	var groups strings.Builder
	for idx, g := range groupNames {
		groups.WriteString(fmt.Sprintf("'%s'", g))

		if idx < len(groupNames)-1 {
			groups.WriteString(",")
		}
	}

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permissions_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_permissions.owner_kind = ?", "group").
		Where(fmt.Sprintf("album_permissions.owner_id = ANY(ARRAY[%s])", groups.String()))

	tx.Find(&albums)
	if tx.Error != nil {
		if !a.circuitBreaker.IsAvailable() {
			return []entity.Album{}, common.NewPostgresNotAvailableError("pg not available while fetching albums by groups")
		}
		return []entity.Album{}, common.NewInternalError(tx.Error, fmt.Sprintf("failed to get albums by groups [%+v]", groupNames))
	}

	if len(albums) == 0 {
		return []entity.Album{}, common.NewEntityNotFound(fmt.Sprintf("albums not found by groups [%+v]", groupNames))
	}

	entities := albums.Merge()

	return entities, nil
}

func (a *AlbumPostgresRepo) SetPermissions(ctx context.Context, albumId string, permissions []entity.AlbumPermission) error {
	if !a.circuitBreaker.IsAvailable() {
		return common.NewPostgresNotAvailableError("pg not available while setting permissions")
	}

	mapper := func(perms []entity.Permission) models.PermissionIDs {
		mperms := make(models.PermissionIDs, 0, len(permissions))
		for _, p := range perms {
			mperms = append(mperms, models.PermissionID(p.String()))
		}
		return mperms
	}

	tx := a.db.WithContext(ctx).Begin()
	for _, permission := range permissions {
		tx.Create(&models.AlbumPermissions{
			AlbumID:     albumId,
			OwnerID:     permission.OwnerID,
			OwnerKind:   permission.OwnerKind,
			Permissions: mapper(permission.Permissions),
		})
		if tx.Error != nil {
			if !a.circuitBreaker.IsAvailable() {
				return common.NewPostgresNotAvailableError("pg not available while setting permissions")
			}
			return common.NewInternalError(tx.Error, fmt.Sprintf("failed to set permissions for album '%s' for owner '%s'", albumId, permission.OwnerID))
		}
	}

	if err := tx.Commit().Error; err != nil {
		if !a.circuitBreaker.IsAvailable() {
			return common.NewPostgresNotAvailableError("pg not available while setting permissions")
		}
		return common.NewInternalError(tx.Error, fmt.Sprintf("failed to set permissions for album '%s'", albumId))
	}

	return nil
}

func (a *AlbumPostgresRepo) RemovePermissions(ctx context.Context, albumId string) error {
	if !a.circuitBreaker.IsAvailable() {
		return common.NewPostgresNotAvailableError("pg not available while removing permissions")
	}
	tx := a.db.WithContext(ctx).Where("album_id = ?", albumId).Delete(&models.AlbumPermissions{})
	if tx.Error != nil {
		if !a.circuitBreaker.IsAvailable() {
			return common.NewPostgresNotAvailableError("pg not available while removing permissions")
		}
		return common.NewInternalError(tx.Error, fmt.Sprintf("failed to remove permissions of album '%s'", albumId))
	}
	return nil
}

func (a *AlbumPostgresRepo) checkNetworkError(err error) (isOpen bool) {
	isOpen = a.circuitBreaker.BreakOnNetworkError(err)
	if isOpen {
		zap.S().Warn("circuit breaker is now open")
	}
	return
}
