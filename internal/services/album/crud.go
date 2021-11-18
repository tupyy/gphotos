package album

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services"
	"github.com/tupyy/gophoto/internal/services/media"
)

type Service struct {
	albumRepo    domain.Album
	mediaService *media.Service
}

func New(albumRepo domain.Album, media *media.Service) *Service {
	return &Service{albumRepo, media}
}

func (s *Service) Create(ctx context.Context, newAlbum entity.Album) (int32, error) {
	// generate bucket name
	bucketID := strings.ReplaceAll(uuid.New().String(), "-", "")
	n := strings.ReplaceAll(strings.ToLower(newAlbum.Name), " ", "-")
	newAlbum.Bucket = fmt.Sprintf("%s-%s", n, bucketID[:8])

	// create the bucket
	if err := s.mediaService.CreateBucket(ctx, newAlbum.Bucket); err != nil {
		return 0, fmt.Errorf("%w '%s': %v", services.ErrCreateBucket, newAlbum.Name, err)
	}

	albumID, err := s.albumRepo.Create(ctx, newAlbum)
	if err != nil {
		return 0, fmt.Errorf("%w '%s': %v", services.ErrCreateAlbum, newAlbum.Name, err)
	}

	return albumID, nil
}

func (s *Service) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	err := s.albumRepo.Update(ctx, album)
	if err != nil {
		return album, fmt.Errorf("%w '%d': %v", services.ErrUpdateAlbum, album.ID, err)
	}

	return album, nil
}

func (s *Service) Delete(ctx context.Context, album entity.Album) error {
	err := s.mediaService.DeleteBucket(ctx, album.Bucket)
	if err != nil {
		return fmt.Errorf("%w '%s': %v", services.ErrDeleteBucket, album.Bucket, err)
	}

	err = s.albumRepo.Delete(ctx, album.ID)
	if err != nil {
		return fmt.Errorf("%w '%d': %v", services.ErrDeleteAlbum, album.ID, err)
	}

	return nil
}
