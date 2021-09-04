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
