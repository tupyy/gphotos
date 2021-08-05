package album

import (
	"context"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repo"
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
		return -1, fmt.Errorf("%w cannot create album: %+v", repo.ErrCreateAlbum, err)
	}

	tx := a.db.WithContext(ctx).Begin()

	m := toModel(album)

	result := tx.Create(&m)
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot create album: %v", album)

		return -1, fmt.Errorf("%w cannot create album %+v", repo.ErrCreateAlbum, result.Error)
	}

	// create permissions entries
	if len(album.UserPermissions) > 0 {
		permModels := toUserPermissionsModels(m.ID, album.UserPermissions)

		if result := tx.CreateInBatches(permModels, len(permModels)); result.Error != nil {
			logger.WithError(result.Error).Warnf("cannot create album user permissions: %v", permModels)
			tx.Rollback()

			return -1, fmt.Errorf("%w cannot create user permissions %+v", repo.ErrCreateAlbum, result.Error)
		}
	}

	if len(album.GroupPermissions) > 0 {
		permModels := toGroupPermissionsModels(m.ID, album.GroupPermissions)

		if result := tx.CreateInBatches(permModels, len(permModels)); result.Error != nil {
			logger.WithError(result.Error).Warnf("cannot create album group permissions: %v", permModels)
			tx.Rollback()

			return -1, fmt.Errorf("%w cannot create group permissions %+v", repo.ErrCreateAlbum, result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.WithError(result.Error).Warnf("error commit album: %v", album)

		return -1, fmt.Errorf("%w cannot create album %+v", repo.ErrCreateAlbum, result.Error)
	}

	return m.ID, nil
}

func (a *AlbumPostgresRepo) Delete(ctx context.Context, id int32) error {
	if res := a.db.WithContext(ctx).Delete(&models.Album{}, id); res.Error != nil {
		return fmt.Errorf("%w %+v", repo.ErrDeleteAlbum, res.Error)
	}

	return nil
}

func (a *AlbumPostgresRepo) Update(ctx context.Context, album entity.Album) error {
	var ca customAlbum

	logger := logutil.GetDefaultLogger()

	if err := album.Validate(); err != nil {
		return fmt.Errorf("%w %+v album_id=%d", repo.ErrUpdateAlbum, err, album.ID)
	}

	tx := a.db.WithContext(ctx).Table("album").Where("id = ?", album.ID).First(&ca)
	if tx.Error != nil {
		return fmt.Errorf("%w %v album_id=%d", repo.ErrAlbumNotFound, tx.Error, album.ID)
	}

	// update all fields except the owner
	newAlbum := entity.Album{
		Name:        album.Name,
		CreatedAt:   album.CreatedAt,
		Description: album.Description,
		Location:    album.Location,
		OwnerID:     album.OwnerID,
	}

	tx = a.db.WithContext(ctx).Begin()

	m := toModel(newAlbum)
	m.ID = album.ID

	result := tx.Save(&m)
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot update album: %v", album)

		return fmt.Errorf("%w %+v", repo.ErrUpdateAlbum, result.Error)
	}

	// update user permissions
	result = tx.Where("album_id = ?", album.ID).Delete(models.AlbumUserPermissions{})
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot delete user permissions while updating album: %v", album)
		tx.Rollback()

		return fmt.Errorf("%w %+v album_id: %d", repo.ErrUpdateAlbum, result.Error, album.ID)
	}

	if len(album.UserPermissions) != 0 {
		permModels := toUserPermissionsModels(m.ID, album.UserPermissions)

		if result := tx.CreateInBatches(permModels, len(permModels)); result.Error != nil {
			logger.WithError(result.Error).Warnf("cannot create album user permissions: %v", permModels)
			tx.Rollback()

			return fmt.Errorf("%w %+v album_id: %d", repo.ErrUpdateAlbum, result.Error, album.ID)
		}
	}

	result = tx.Where("album_id = ?", album.ID).Delete(models.AlbumGroupPermissions{})
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot delete group permissions while updating album: %v", album)
		tx.Rollback()

		return fmt.Errorf("%w %+v album_id: %d", repo.ErrUpdateAlbum, result.Error, album.ID)
	}

	if len(album.GroupPermissions) != 0 {
		permModels := toGroupPermissionsModels(m.ID, album.GroupPermissions)

		if result := tx.CreateInBatches(permModels, len(permModels)); result.Error != nil {
			logger.WithError(result.Error).Warnf("cannot create album group permissions: %v", permModels)
			tx.Rollback()

			return fmt.Errorf("%w %+v album_id: %d", repo.ErrUpdateAlbum, result.Error, album.ID)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.WithError(result.Error).WithFields(logrus.Fields{
			"new album": fmt.Sprintf("%+v", album),
			"old album": fmt.Sprintf("%+v", ca),
		}).Warnf("error commit album: %v", album)

		return fmt.Errorf("%w %+v album_id: %d", repo.ErrUpdateAlbum, result.Error, album.ID)
	}

	return nil
}

// Get returns all the albums sorted by id.
func (a *AlbumPostgresRepo) Get(ctx context.Context) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Order("album.id").
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w no album found", repo.ErrAlbumNotFound)
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
		Select(`album.*, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Where("album.id = ?", id).
		Find(&albums)
	if tx.Error != nil {
		return entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return entity.Album{}, fmt.Errorf("%w no album found with id %d", repo.ErrAlbumNotFound, id)
	}

	entities := albums.Merge()

	return entities[0], nil
}

func (a *AlbumPostgresRepo) GetByOwnerID(ctx context.Context, ownerID string) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Where("album.owner_id = ?", ownerID).
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w ownerr id %s", repo.ErrAlbumNotFound, ownerID)
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByUserID returns a list of albums for which the user has at one permission set.
func (a *AlbumPostgresRepo) GetByUserID(ctx context.Context, userID string) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Where("album_user_permissions.user_id = ?", userID).
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("user_id", userID).Warn("no album found")

		return []entity.Album{}, repo.ErrAlbumNotFound
	}

	entities := albums.Merge()

	return entities, nil
}

// GetAlbumsByGroup returns a list of albums for which the group has at one permission set.
func (a *AlbumPostgresRepo) GetByGroupID(ctx context.Context, groupName string) ([]entity.Album, error) {
	var albums customAlbums

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Where("album_group_permissions.group_id = ?", groupName).
		Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, fmt.Errorf("%w no album found with id %s", repo.ErrAlbumNotFound, groupName)
	}

	entities := albums.Merge()

	return entities, nil
}
