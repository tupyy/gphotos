package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services"
	"github.com/tupyy/gophoto/internal/services/image"
)

type MediaType int

const (
	Photo MediaType = iota
	Video
)

type Service struct {
	repos domain.Repositories
}

func New(repos domain.Repositories) *Service {
	return &Service{repos}
}

func (s *Service) List(ctx context.Context, album entity.Album) ([]entity.Media, error) {
	minioRepo := s.repos[domain.MinioRepoName].(domain.Store)

	media, err := minioRepo.ListBucket(ctx, album.Bucket)
	if err != nil {
		return []entity.Media{}, fmt.Errorf("%w album '%d': %v", services.ErrListBucket, album.ID, err)
	}

	return media, nil
}

func (s *Service) GetPhoto(ctx context.Context, bucket, filename string) (io.Reader, error) {
	minioRepo := s.repos[domain.MinioRepoName].(domain.Store)

	r, err := minioRepo.GetFile(ctx, bucket, filename)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Service) SaveMedia(ctx context.Context, bucket, filename string, r io.Reader, size int64, mediaType MediaType) error {
	minioRepo := s.repos[domain.MinioRepoName].(domain.Store)

	switch mediaType {
	case Photo:
		return processPhoto(ctx, minioRepo, bucket, filename, r, size)
	case Video:
		return fmt.Errorf("not implementated")
	default:
		return fmt.Errorf("media type not supported")
	}

}

func processPhoto(ctx context.Context, repo domain.Store, bucket, filename string, r io.Reader, size int64) error {
	err := repo.PutFile(ctx, conf.GetMinioTemporaryBucket(), filename, size, r)
	if err != nil {
		return fmt.Errorf("failed to copy file to temporary bucket: %v", err)
	}

	// do image processing
	var imgBuffer bytes.Buffer
	var imgThumbnailBuffer bytes.Buffer
	if err := image.Process(r, &imgBuffer, &imgThumbnailBuffer); err != nil {
		return fmt.Errorf("failed to process image: %v", err)
	}

	basename := strings.Split(filename, ".")[0]

	if err := repo.PutFile(ctx, bucket, fmt.Sprintf("%s.jpg", basename), int64(imgBuffer.Len()), &imgBuffer); err != nil {
		return fmt.Errorf("failed to copy processed image to bucket '%s': %v", bucket, err)
	}

	if err := repo.PutFile(ctx, bucket, fmt.Sprintf("%s_thumbnail.jpg", basename), int64(imgThumbnailBuffer.Len()), &imgThumbnailBuffer); err != nil {
		return fmt.Errorf("failed to copy thumbnail image to bucket '%s': %v", bucket, err)
	}

	return nil
}
