package album

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	repo "github.com/tupyy/gophoto/internal/domain"
	albumFilters "github.com/tupyy/gophoto/internal/domain/filters/album"
	"github.com/tupyy/gophoto/internal/domain/models"
	"github.com/tupyy/gophoto/internal/entity"
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
	var ca albumJoinRow

	logger := logutil.GetDefaultLogger()

	if err := album.Validate(); err != nil {
		return fmt.Errorf("%w %+v album_id=%d", repo.ErrUpdateAlbum, err, album.ID)
	}

	tx := a.db.WithContext(ctx).Table("album").Where("id = ?", album.ID).First(&ca)
	if tx.Error != nil {
		return fmt.Errorf("%w %v album_id=%d", repo.ErrAlbumNotFound, tx.Error, album.ID)
	}

	newAlbum := entity.Album{
		Name:        album.Name,
		CreatedAt:   album.CreatedAt,
		Description: album.Description,
		Location:    album.Location,
		OwnerID:     album.OwnerID,
		Bucket:      album.Bucket,
		Thumbnail:   album.Thumbnail,
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
func (a *AlbumPostgresRepo) Get(ctx context.Context, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").Table("tag").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery)

	for _, f := range filters {
		tx = f(tx)
	}

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().Warn("no albums found")

		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByID return the album if any with id id.
func (a *AlbumPostgresRepo) GetByID(ctx context.Context, id int32) (entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album.id = ?", id).
		Find(&albums)

	if tx.Error != nil {
		return entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("album id", id).Warn("no album found by id")

		return entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities[0], nil
}

// GetByOwnerID return all albums of an user.
// It does not sort or filter the album here. The sorting and filter is done at cache level.
func (a *AlbumPostgresRepo) GetByOwnerID(ctx context.Context, ownerID string, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album.owner_id = ?", ownerID)

	for _, f := range filters {
		tx = f(tx)
	}

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("ownerID", ownerID).Warn("no album found by owner id")

		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByUserID returns a list of albums for which the user has at one permission set.
// It does not sort or filter the album here. The sorting and filter is done at cache level.
func (a *AlbumPostgresRepo) GetByUserID(ctx context.Context, userID string, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_user_permissions.user_id = ?", userID)

	for _, f := range filters {
		tx = f(tx)
	}

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("userID", userID).Warn("no album found by user id")

		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}

// GetAlbumsByGroup returns a list of albums for which the group has at one permission set.
// It does not sort or filter the album here. The sorting and filter is done at cache level.
func (a *AlbumPostgresRepo) GetByGroupName(ctx context.Context, groupName string, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_group_permissions.group_name= ?", groupName)

	for _, f := range filters {
		tx = f(tx)
	}

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("group", groupName).Warn("no album found by group name")

		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByGroups returns a list of albums with at least one persmission for at least on group in the list.
func (a *AlbumPostgresRepo) GetByGroups(ctx context.Context, groupNames []string, filters albumFilters.Filters) ([]entity.Album, error) {
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
		Select(`album.*, tags.name as tag_name,tags.color as tag_color, album_user_permissions.permissions as user_permissions, album_user_permissions.user_id as user_id,
				album_group_permissions.permissions as group_permissions, album_group_permissions.group_name as group_name`).
		Joins("LEFT JOIN album_user_permissions ON (album.id = album_user_permissions.album_id)").
		Joins("LEFT JOIN album_group_permissions ON (album.id = album_group_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where(fmt.Sprintf("album_group_permissions.group_name = ANY(ARRAY[%s])", groups.String()))

	for _, f := range filters {
		tx = f(tx)
	}

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("group names", fmt.Sprintf("%+v", groupNames)).Warn("no album found by group name")

		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}
