package album

import (
	"context"
	"fmt"
	"sort"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo/postgres"
	"github.com/tupyy/gophoto/models"
	"github.com/tupyy/gophoto/utils/logutil"
	"github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type AlbumPostgresRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func NewPostgresRepo(client pgclient.Client) (*AlbumPostgresRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &AlbumPostgresRepo{}, err
	}

	return &AlbumPostgresRepo{gormDB, client}, nil
}

func (a *AlbumPostgresRepo) Create(ctx context.Context, album entity.Album) (albumID int32, err error) {
	logger := logutil.GetDefaultLogger()

	if err := album.Validate(); err != nil {
		return -1, fmt.Errorf("%w cannot create album: %+v", postgres.ErrCreateAlbum, err)
	}

	tx := a.db.WithContext(ctx).Begin()

	m := toModel(album)

	result := tx.Create(&m)
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot create album: %v", album)

		return -1, fmt.Errorf("%w cannot create album %+v", postgres.ErrCreateAlbum, result.Error)
	}

	// create permissions entries
	if len(album.UserPermissions) > 0 {
		permModels := toUserPermissionsModels(m.ID, album.UserPermissions)

		if result := tx.CreateInBatches(permModels, len(permModels)); result.Error != nil {
			logger.WithError(result.Error).Warnf("cannot create album user permissions: %v", permModels)
			tx.Rollback()

			return -1, fmt.Errorf("%w cannot create user permissions %+v", postgres.ErrCreateAlbum, result.Error)
		}
	}

	if len(album.GroupPermissions) > 0 {
		permModels := toGroupPermissionsModels(m.ID, album.GroupPermissions)

		if result := tx.CreateInBatches(permModels, len(permModels)); result.Error != nil {
			logger.WithError(result.Error).Warnf("cannot create album group permissions: %v", permModels)
			tx.Rollback()

			return -1, fmt.Errorf("%w cannot create group permissions %+v", postgres.ErrCreateAlbum, result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.WithError(result.Error).Warnf("error commit album: %v", album)

		return -1, fmt.Errorf("%w cannot create album %+v", postgres.ErrCreateAlbum, result.Error)
	}

	return m.ID, nil
}

func (a *AlbumPostgresRepo) Delete(ctx context.Context, id int32) error {
	if res := a.db.WithContext(ctx).Delete(&models.Album{}, id); res.Error != nil {
		return fmt.Errorf("%w %+v", postgres.ErrDeleteAlbum, res.Error)
	}

	return nil
}

func (a *AlbumPostgresRepo) Update(ctx context.Context, album entity.Album) error {
	var oldAlbum entity.Album

	if err := album.Validate(); err != nil {
		return fmt.Errorf("%w cannot create album: %+v", postgres.ErrUpdateAlbum, err)
	}

	tx := a.db.WithContext(ctx).Where("id = ?", album.ID).First(&oldAlbum)
	if tx.Error != nil {
		return fmt.Errorf("%w cannot update album: %v", postgres.ErrAlbumNotFound, tx.Error)
	}

	// update all fields except the owner
	oldAlbum.Name = album.Name
	oldAlbum.CreatedAt = album.CreatedAt
	oldAlbum.Description = album.Description
	oldAlbum.Location = album.Location

	tx = a.db.WithContext(ctx).Save(&oldAlbum)
	if tx.Error != nil {
		return fmt.Errorf("%w cannot update album %v", postgres.ErrUpdateAlbum, tx.Error)
	}

	return nil
}

// Get returns all the albums sorted by id.
func (a *AlbumPostgresRepo) Get(ctx context.Context) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions,
				album_group_permissions.permissions as group_permissions, users.id as user_id,
				groups.id as group_id`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN users ON (album_user_permissions.user_id = users.id)").
		Joins("LEFT JOIN groups ON (album_group_permissions.group_id = groups.id)").
		Order("album.id").
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", postgres.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w no album found", postgres.ErrAlbumNotFound)
	}

	entities := albums.Merge()

	// sort by id
	albumSorter := entity.NewAlbumSorter(entities, func(a1, a2 entity.Album) bool { return a1.ID < a2.ID })
	sort.Sort(albumSorter)

	return entities, nil
}

// GetByID return the album if any with id id.
func (a *AlbumPostgresRepo) GetByID(ctx context.Context, id int32) (entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions,
				album_group_permissions.permissions as group_permissions, users.id as user_id,
				groups.id as group_id`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN users ON (album_user_permissions.user_id = users.id)").
		Joins("LEFT JOIN groups ON (album_group_permissions.group_id = groups.id)").
		Where("album.id = ?", id).
		Find(&albums)
	if tx.Error != nil {
		return entity.Album{}, fmt.Errorf("%w internal error: %v", postgres.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return entity.Album{}, fmt.Errorf("%w no album found with id %d", postgres.ErrAlbumNotFound, id)
	}

	entities := albums.Merge()

	return entities[0], nil
}

func (a *AlbumPostgresRepo) GetAlbumsByOwnerID(ctx context.Context, ownerID int32) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions,
				album_group_permissions.permissions as group_permissions, users.id as user_id,
				groups.id as group_id`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN users ON (album_user_permissions.user_id = users.id)").
		Joins("LEFT JOIN groups ON (album_group_permissions.group_id = groups.id)").
		Where("album.owner_id = ?", ownerID).
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", postgres.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w no album found with id %d", postgres.ErrAlbumNotFound, ownerID)
	}

	entities := albums.Merge()

	return entities, nil
}

// GetAlbumsByUser returns a list of albums for which the user has at one permission set.
func (a *AlbumPostgresRepo) GetAlbumsByUserID(ctx context.Context, userID int32) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions,
				album_group_permissions.permissions as group_permissions, users.id as user_id,
				groups.id as group_id`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN users ON (album_user_permissions.user_id = users.id)").
		Joins("LEFT JOIN groups ON (album_group_permissions.group_id = groups.id)").
		Where("album_user_permissions.user_id = ?", userID).
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", postgres.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w no album found with id %d", postgres.ErrAlbumNotFound, userID)
	}

	entities := albums.Merge()

	return entities, nil
}

// GetAlbumsByGroup returns a list of albums for which the group has at one permission set.
func (a *AlbumPostgresRepo) GetAlbumsByGroupID(ctx context.Context, groupID int32) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions,
				album_group_permissions.permissions as group_permissions, users.id as user_id,
				groups.id as group_id`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN users ON (album_user_permissions.user_id = users.id)").
		Joins("LEFT JOIN groups ON (album_group_permissions.group_id = groups.id)").
		Where("album_group_permissions.group_id = ?", groupID).
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", postgres.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w no album found with id %d", postgres.ErrAlbumNotFound, groupID)
	}

	entities := albums.Merge()

	return entities, nil
}
