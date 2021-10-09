package media

import (
	"context"
	"fmt"
	"io"

	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
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
		return []entity.Media{}, fmt.Errorf("%w failed to read bucket for album '%d'", err, album.ID)
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
