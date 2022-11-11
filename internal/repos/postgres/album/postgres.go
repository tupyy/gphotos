package album

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/entity"
	repo "github.com/tupyy/gophoto/internal/repos"
	"github.com/tupyy/gophoto/internal/repos/models"
	"github.com/tupyy/gophoto/internal/utils/logutil"
	"github.com/tupyy/gophoto/internal/utils/pgclient"
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

func (a *AlbumPostgresRepo) Create(ctx context.Context, album entity.Album) (entity.Album, error) {
	logger := logutil.GetDefaultLogger()

	tx := a.db.WithContext(ctx).Begin()

	m := toModel(album)
	m.ID = xid.New().String()
	album.ID = m.ID

	result := tx.Create(&m)
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot create album: %v", album)

		return album, fmt.Errorf("%w cannot create album %+v", repo.ErrCreateAlbum, result.Error)
	}

	if err := tx.Commit().Error; err != nil {
		logger.WithError(result.Error).Warnf("error commit album: %v", album)

		return album, fmt.Errorf("%w cannot create album %+v", repo.ErrCreateAlbum, result.Error)
	}

	return album, nil
}

func (a *AlbumPostgresRepo) Delete(ctx context.Context, id string) error {
	if res := a.db.WithContext(ctx).Delete(&models.Album{}, id); res.Error != nil {
		return fmt.Errorf("%w %+v", repo.ErrDeleteAlbum, res.Error)
	}

	return nil
}

func (a *AlbumPostgresRepo) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	var ca albumJoinRow

	logger := logutil.GetDefaultLogger()

	tx := a.db.WithContext(ctx).Table("album").Where("id = ?", album.ID).First(&ca)
	if tx.Error != nil {
		return album, fmt.Errorf("%w %v album_id=%s", repo.ErrAlbumNotFound, tx.Error, album.ID)
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

	tx = a.db.WithContext(ctx).Begin()

	m := toModel(newAlbum)
	m.ID = album.ID

	result := tx.Save(&m)
	if result.Error != nil {
		logger.WithError(result.Error).Warnf("cannot update album: %v", album)

		return album, fmt.Errorf("%w %+v", repo.ErrUpdateAlbum, result.Error)
	}

	if err := tx.Commit().Error; err != nil {
		logger.WithError(result.Error).WithFields(logrus.Fields{
			"new album": fmt.Sprintf("%+v", album),
			"old album": fmt.Sprintf("%+v", ca),
		}).Warnf("error commit album: %v", album)

		return album, fmt.Errorf("%w %+v album_id: %s", repo.ErrUpdateAlbum, result.Error, album.ID)
	}

	return album, nil
}

// Get returns all the albums sorted by id.
func (a *AlbumPostgresRepo) Get(ctx context.Context) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").Table("tag").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permission_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery)

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
func (a *AlbumPostgresRepo) GetByID(ctx context.Context, id string) (entity.Album, error) {
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
		return entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return entity.Album{}, fmt.Errorf("album not found")
	}

	entities := albums.Merge()

	return entities[0], nil
}

// GetByOwnerID return all albums of an user.
func (a *AlbumPostgresRepo) GetByOwner(ctx context.Context, owner string) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permissions_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album.owner_id = ?", owner)

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}

// GetByUserID returns a list of albums for which the user has at one permission set.
func (a *AlbumPostgresRepo) GetByUser(ctx context.Context, username string) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permissions_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_permissions.owner_kind = ?", "user").
		Where("album_permissions.owner_id = ?", username)

	tx.Find(&albums)
	if tx.Error != nil {
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}

// GetAlbumsByGroup returns a list of albums for which the group has at one permission set.
func (a *AlbumPostgresRepo) GetByGroupName(ctx context.Context, groupName string) ([]entity.Album, error) {
	var albums albumJoinRows

	tagSubQuery := a.db.WithContext(ctx).Table("tag").
		Select("id, albums_tags.album_id, name, color").
		Joins("JOIN albums_tags ON (albums_tags.tag_id = tag.id)")

	tx := a.db.WithContext(ctx).Table("album").
		Select(`album.*, tags.id as tag_id, tags.name as tag_name,tags.color as tag_color, album_permissions.permissions as permissions, album_permissions.owner_id as permissions_owner_id,
				album_permissions.owner_kind as permission_owner_kind`).
		Joins("LEFT JOIN album_permissions ON (album.id = album_permissions.album_id)").
		Joins("LEFT JOIN (?) as tags ON (tags.album_id = album.id)", tagSubQuery).
		Where("album_permissions.owner_kind = ?", "group").
		Where("album_permissions.owner_id = ?", groupName)

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
func (a *AlbumPostgresRepo) GetByGroups(ctx context.Context, groupNames []string) ([]entity.Album, error) {
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
		return []entity.Album{}, fmt.Errorf("%w internal error: %v", repo.ErrInternalError, tx.Error)
	}

	if len(albums) == 0 {
		logutil.GetDefaultLogger().WithField("group names", fmt.Sprintf("%+v", groupNames)).Warn("no album found by group name")

		return []entity.Album{}, nil
	}

	entities := albums.Merge()

	return entities, nil
}
