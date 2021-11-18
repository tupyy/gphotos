package minio

import (
	"context"
	"io"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/entity"
)

type minioRepoCache struct {
	repo  domain.Store
	cache *gocache.Cache
}

func NewCacheRepo(r domain.Store, ttl time.Duration, cleanInterval time.Duration) domain.Store {
	return &minioRepoCache{
		repo:  r,
		cache: gocache.New(ttl, cleanInterval),
	}
}

func (r *minioRepoCache) GetFile(ctx context.Context, bucket, filename string) (io.ReadSeeker, error) {
	return r.repo.GetFile(ctx, bucket, filename)
}

func (r *minioRepoCache) PutFile(ctx context.Context, bucket, filename string, size int64, reader io.Reader) error {
	err := r.repo.PutFile(ctx, bucket, filename, size, reader)
	if err != nil {
		return err
	}

	// invalid the cache for this bucket
	r.cache.Delete(bucket)

	return nil
}

func (r *minioRepoCache) ListBucket(ctx context.Context, bucket string) ([]entity.Media, error) {
	items, found := r.cache.Get(bucket)
	if !found {
		items, err := r.repo.ListBucket(ctx, bucket)
		if err != nil {
			return []entity.Media{}, err
		}

		r.cache.Set(bucket, items, gocache.DefaultExpiration)

		return items, nil
	}

	return items.([]entity.Media), nil
}

func (r *minioRepoCache) DeleteFile(ctx context.Context, bucket, filename string) error {
	err := r.repo.DeleteFile(ctx, bucket, filename)
	if err != nil {
		return err
	}

	// invalid the cache for this bucket
	r.cache.Delete(bucket)

	return nil
}

func (r *minioRepoCache) CreateBucket(ctx context.Context, bucket string) error {
	return r.repo.CreateBucket(ctx, bucket)
}

func (r *minioRepoCache) DeleteBucket(ctx context.Context, bucket string) error {
	err := r.repo.DeleteBucket(ctx, bucket)
	if err != nil {
		return err
	}

	// invalid the cache for this bucket
	r.cache.Delete(bucket)

	return nil
}
