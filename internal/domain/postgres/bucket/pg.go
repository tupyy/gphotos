package bucket

import (
	"context"
	"fmt"

	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/models"
	"github.com/tupyy/gophoto/utils/pgclient"
	"gorm.io/gorm"
)

type BucketPostgresRepo struct {
	db     *gorm.DB
	client pgclient.Client
}

func NewPostgresRepo(client pgclient.Client) (*BucketPostgresRepo, error) {
	config := gorm.Config{
		SkipDefaultTransaction: true, // No need transaction for those use cases.
	}

	gormDB, err := client.Open(config)
	if err != nil {
		return &BucketPostgresRepo{}, err
	}

	return &BucketPostgresRepo{gormDB, client}, nil
}

func (b *BucketPostgresRepo) Get(ctx context.Context, albumID int32) (entity.Bucket, error) {
	var bucket models.Bucket

	if err := b.db.WithContext(ctx).Where("album_id = ?", albumID).First(&bucket).Error; err != nil {
		return entity.Bucket{}, fmt.Errorf("%w failed to get bucket for album %d", err, albumID)
	}

	e := entity.Bucket{
		AlbumID: albumID,
		Urn:     bucket.Urn,
	}

	return e, nil
}

func (b *BucketPostgresRepo) Create(ctx context.Context, bucket entity.Bucket) error {
	m := models.Bucket{
		Urn:     bucket.Urn,
		AlbumID: bucket.AlbumID,
	}

	if err := b.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("%w failed to create bucket %v", err, bucket)
	}

	return nil
}

func (b *BucketPostgresRepo) Delete(ctx context.Context, bucket entity.Bucket) error {
	return nil
}
