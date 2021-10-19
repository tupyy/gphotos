package album

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

type Service struct {
	repos domain.Repositories
}

func New(repos domain.Repositories) *Service {
	return &Service{repos}
}

// TODO remove the logs for returned error and create typed errors
func (s *Service) Create(ctx context.Context, newAlbum entity.Album) (int32, error) {
	minioRepo := s.repos[domain.MinioRepoName].(domain.Store)
	albumRepo := s.repos[domain.AlbumRepoName].(domain.Album)

	logger := logutil.GetLogger(ctx)

	// generate bucket name
	bucketID := strings.ReplaceAll(uuid.New().String(), "-", "")
	n := strings.ReplaceAll(strings.ToLower(newAlbum.Name), " ", "-")
	newAlbum.Bucket = fmt.Sprintf("%s-%s", n, bucketID[:8])

	// create the bucket
	if err := minioRepo.CreateBucket(ctx, newAlbum.Bucket); err != nil {
		logger.WithError(err).Error("failed to create bucket")

		return 0, fmt.Errorf("failed to create album '%s': %v", newAlbum.Name, err)
	}

	albumID, err := albumRepo.Create(ctx, newAlbum)
	if err != nil {
		return 0, fmt.Errorf("failed to create album '%s': %v", newAlbum.Name, err)
	}

	return albumID, nil
}

func (s *Service) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	albumRepo := s.repos[domain.AlbumRepoName].(domain.Album)

	logger := logutil.GetLogger(ctx)

	err := albumRepo.Update(ctx, album)
	if err != nil {
		logger.WithError(err).WithField("album id", album.ID).Error("failed to update album")

		return album, fmt.Errorf("failed to update album '%d': %v", album.ID, err)
	}

	return album, nil
}

func (s *Service) Delete(ctx context.Context, album entity.Album) error {
	minioRepo := s.repos[domain.MinioRepoName].(domain.Store)
	albumRepo := s.repos[domain.AlbumRepoName].(domain.Album)

	logger := logutil.GetLogger(ctx)

	err := minioRepo.DeleteBucket(ctx, album.Bucket)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"bucket":   album.Bucket,
			"album id": album.ID,
		}).WithError(err).Error("failed to remove album's bucket")

		return fmt.Errorf("failed to remove album's bucket '%s': %v", album.Bucket, err)
	}

	err = albumRepo.Delete(ctx, album.ID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"bucket":   album.Bucket,
			"album id": album.ID,
		}).WithError(err).Error("failed to remove album")

		return fmt.Errorf("failed to remove album '%d': %v", album.ID, err)
	}

	return nil
}
