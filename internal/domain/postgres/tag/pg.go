package tag

import (
	"context"

	"github.com/tupyy/gophoto/internal/domain/models"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type TagRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func NewPostgresRepo(client pgclient.Client) (*TagRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &TagRepo{}, err
	}

	return &TagRepo{gormDB, client}, nil
}

func (t *TagRepo) Create(ctx context.Context, tag entity.Tag) (int32, error) {
	tagModel := models.Tag{
		UserID: tag.UserID,
		Name:   tag.Name,
		Color:  tag.Color,
	}

	if err := t.db.WithContext(ctx).Create(&tagModel).Error; err != nil {
		return 0, err
	}

	return tagModel.ID, nil
}

func (t *TagRepo) Update(ctx context.Context, tag entity.Tag) error {
	tagModel := models.Tag{
		ID:     tag.ID,
		UserID: tag.UserID,
		Name:   tag.Name,
		Color:  tag.Color,
	}

	if err := t.db.WithContext(ctx).Save(&tagModel).Error; err != nil {
		return err
	}

	return nil
}

func (t *TagRepo) Delete(ctx context.Context, id int32) error {
	return t.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Tag{}).Error
}

func (t *TagRepo) GetByUser(ctx context.Context, userID string) ([]entity.Tag, error) {
	pgModels := []struct {
		ID     int32   `gorm:"primary_key;column:id;type:INT4;"`
		Name   string  `gorm:"column:name;type:TEXT;"`
		Color  *string `gorm:"column:color;type:TEXT;"`
		UserID string  `gorm:"column:user_id;type:TEXT;"`
	}{}

	tx := t.db.WithContext(ctx).Table("tag").
		Select("id, name, user_id, color").
		Joins("LEFT JOIN albums_tags as at ON at.tag_id = tag.id").
		Group("id").
		Where("user_id = ?", userID).Find(&pgModels)

	if tx.Error != nil {
		return []entity.Tag{}, tx.Error
	}

	tags := make([]entity.Tag, 0, len(pgModels))

	for _, m := range pgModels {
		tag := entity.Tag{
			ID:     m.ID,
			Name:   m.Name,
			Color:  m.Color,
			UserID: m.UserID,
		}

		//if len(m.AlbumIDs) > 0 {
		//	copy(tag.Albums, m.AlbumIDs)
		//}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (t *TagRepo) GetByName(ctx context.Context, userID, name string) (entity.Tag, error) {
	pgModel := struct {
		ID     int32   `gorm:"primary_key;column:id;type:INT4;"`
		Name   string  `gorm:"column:name;type:TEXT;"`
		Color  *string `gorm:"column:color;type:TEXT;"`
		UserID string  `gorm:"column:user_id;type:TEXT;"`
	}{}

	tx := t.db.WithContext(ctx).Table("tag").
		Select("id, name, user_id, color").
		Joins("LEFT JOIN albums_tags as at ON at.tag_id = tag.id").
		Where("user_id = ?", userID).
		Where("name = ?", name).
		Group("id").
		First(&pgModel)

	if tx.Error != nil {
		return entity.Tag{}, tx.Error
	}

	tag := entity.Tag{
		ID:     pgModel.ID,
		Name:   pgModel.Name,
		Color:  pgModel.Color,
		UserID: pgModel.UserID,
	}

	// if len(pgModel.AlbumIDs) > 0 {
	// 	copy(tag.Albums, pgModel.AlbumIDs)
	// }

	return tag, nil
}

func (t *TagRepo) GetByID(ctx context.Context, userID string, id int32) (entity.Tag, error) {
	pgModel := struct {
		ID     int32   `gorm:"primary_key;column:id;type:INT4;"`
		Name   string  `gorm:"column:name;type:TEXT;"`
		Color  *string `gorm:"column:color;type:TEXT;"`
		UserID string  `gorm:"column:user_id;type:TEXT;"`
	}{}

	tx := t.db.WithContext(ctx).Table("tag").
		Select("id, name, user_id, color").
		Joins("LEFT JOIN albums_tags as at ON at.tag_id = tag.id").
		Where("id = ?", id).
		Where("user_id = ?", userID).
		First(&pgModel)

	if tx.Error != nil {
		return entity.Tag{}, tx.Error
	}

	tag := entity.Tag{
		ID:     pgModel.ID,
		Name:   pgModel.Name,
		Color:  pgModel.Color,
		UserID: pgModel.UserID,
	}

	// if len(pgModel.AlbumIDs) > 0 {
	// 	copy(tag.Albums, pgModel.AlbumIDs)
	// }

	return tag, nil
}

func (t *TagRepo) GetByAlbum(ctx context.Context, albumID int32) ([]entity.Tag, error) {
	pgModels := []struct {
		ID     int32   `gorm:"primary_key;column:id;type:INT4;"`
		Name   string  `gorm:"column:name;type:TEXT;"`
		Color  *string `gorm:"column:color;type:TEXT;"`
		UserID string  `gorm:"column:user_id;type:TEXT;"`
	}{}

	tx := t.db.WithContext(ctx).Table("tag").
		Select("id, name, user_id, color").
		Joins("LEFT JOIN albums_tags as at ON at.tag_id = tag.id").
		Where("album_id = ?", albumID).
		Find(&pgModels)

	if tx.Error != nil {
		return []entity.Tag{}, tx.Error
	}

	tags := make([]entity.Tag, 0, len(pgModels))

	for _, m := range pgModels {
		tag := entity.Tag{
			ID:     m.ID,
			Name:   m.Name,
			Color:  m.Color,
			UserID: m.UserID,
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (t *TagRepo) Associate(ctx context.Context, albumID, tagID int32) error {
	pgModel := models.AlbumsTags{
		AlbumID: albumID,
		TagID:   tagID,
	}

	return t.db.WithContext(ctx).Create(&pgModel).Error
}

func (t *TagRepo) Dissociate(ctx context.Context, albumID, tagID int32) error {
	return t.db.WithContext(ctx).Where("album_id = ?", albumID).Where("tag_id = ?", tagID).Delete(&models.AlbumsTags{}).Error
}
