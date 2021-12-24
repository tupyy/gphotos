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

func (t *TagRepo) Get(ctx context.Context, userID string) ([]entity.Tag, error) {
	tagsModels := []models.Tag{}

	if err := t.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tagsModels).Error; err != nil {
		return []entity.Tag{}, err
	}

	tags := make([]entity.Tag, 0, len(tagsModels))

	for _, m := range tagsModels {
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

func (t *TagRepo) GetByName(ctx context.Context, userID, name string) (entity.Tag, error) {
	tagModel := models.Tag{}

	if err := t.db.WithContext(ctx).Where("user_id = ?", userID).Where("name = ?", name).First(tagModel).Error; err != nil {
		return entity.Tag{}, err
	}

	return entity.Tag{
		ID:     tagModel.ID,
		Name:   tagModel.Name,
		Color:  tagModel.Color,
		UserID: tagModel.UserID,
	}, nil
}

func (t *TagRepo) Associate(ctx context.Context, albumID, tagID int32) error {
	pgModel := models.AlbumsTags{
		AlbumID: albumID,
		TagID:   tagID,
	}

	return t.db.WithContext(ctx).Create(&pgModel).Error
}
