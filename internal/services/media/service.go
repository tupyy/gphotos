package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/image"
	"go.uber.org/zap"
)

// Store describe photo store operations
type MinioRepository interface {
	// GetFile returns a reader to file.
	GetFile(ctx context.Context, bucket, filename string) (io.ReadSeeker, map[string]string, error)
	// PutFile save a file to a bucket.
	PutFile(ctx context.Context, bucket, filename string, size int64, r io.Reader, metadata map[string]string) error
	// ListFiles list the content of a bucket
	ListBucket(ctx context.Context, bucket string) ([]entity.Media, error)
	// DeleteFile deletes a file from a bucket.
	DeleteFile(ctx context.Context, bucket, filename string) error
	// CreateBucket create a bucket.
	CreateBucket(ctx context.Context, bucket string, tags map[string]string) error
	// DeleteBucket removes bucket.
	DeleteBucket(ctx context.Context, bucket string) error
	// set tags to bucket
	SetBucketTagging(ctx context.Context, bucket string, tags map[string]string) error
	// get bucket tags
	GetBucketTagging(ctx context.Context, buclet string) (map[string]string, error)
}

type MediaType int

const (
	Photo MediaType = iota
	Video
)

type Service struct {
	repo MinioRepository
}

func New(repo MinioRepository) *Service {
	return &Service{repo}
}

func (s *Service) CreateBucket(ctx context.Context, bucket string, tags map[string]string) error {
	return s.repo.CreateBucket(ctx, bucket, tags)
}

// DeleteBucket does not delete the bucket. Only set the tags delete_at
func (s *Service) DeleteBucket(ctx context.Context, bucket string) error {
	bucketTags, err := s.repo.GetBucketTagging(ctx, bucket)
	if err != nil {
		return err
	}
	bucketTags["album/deleted_at"] = time.Now().Format(time.RFC3339)
	return s.repo.SetBucketTagging(ctx, bucket, bucketTags)
}

func (s *Service) ListBucket(ctx context.Context, bucket string) ([]entity.Media, error) {
	media, err := s.repo.ListBucket(ctx, bucket)
	if err != nil {
		return []entity.Media{}, fmt.Errorf("failed to list bucket '%s': %v", bucket, err)
	}

	// if a media has no thumbnail, create it now
	for _, m := range media {
		if len(m.Thumbnail) == 0 {
			r, _, err := s.GetPhoto(ctx, m.Bucket, m.Filename)
			if err != nil {
				zap.S().Errorw("failed to get photo from repo", "error", err, "filename", m.Filename)
				continue
			}

			if err := createThumbnail(ctx, s.repo, m.Bucket, m.Filename, r); err != nil {
				zap.S().Errorw("failed to create thumbnail", "error", err, "filename", m.Filename)
				continue
			}

			m.Thumbnail = fmt.Sprintf("thumbnail/%s", m.Filename)
		}
	}

	ms := newSorter(media)
	sort.Sort(ms)

	return ms.medias, nil
}

func (s *Service) GetPhoto(ctx context.Context, bucket, filename string) (io.ReadSeeker, map[string]string, error) {
	r, metadata, err := s.repo.GetFile(ctx, bucket, filename)
	if err != nil {
		return nil, nil, err
	}

	return r, metadata, nil
}

func (s *Service) Save(ctx context.Context, bucket, filename string, r io.ReadSeeker, mediaType MediaType) error {
	switch mediaType {
	case Photo:
		if err := processPhoto(ctx, s.repo, bucket, filename, r); err != nil {
			return err
		}

		if err := createThumbnail(ctx, s.repo, bucket, filename, r); err != nil {
			return err
		}

		return nil
	case Video:
		return fmt.Errorf("not implementated")
	default:
		return fmt.Errorf("media type not supported")
	}
}

func (s *Service) Delete(ctx context.Context, bucket, filename string) error {
	if strings.Index(filename, "/") > 0 {
		parts := strings.Split(filename, "/")

		// delete thumbnail
		if err := s.repo.DeleteFile(ctx, bucket, fmt.Sprintf("thumbnail/%s", parts[len(parts)-1])); err != nil {
			return err
		}

	}

	return s.repo.DeleteFile(ctx, bucket, filename)
}

func processPhoto(ctx context.Context, repo MinioRepository, bucket, filename string, r io.ReadSeeker) error {
	var imgBuffer bytes.Buffer

	if err := image.Process(r, &imgBuffer); err != nil {
		return fmt.Errorf("failed to process image: %v", err)
	}

	basename := strings.Split(filename, ".")[0]

	_, _ = r.Seek(0, 0)

	metadata, err := image.Metadata(r)
	if err != nil {
		return err
	}

	if err := repo.PutFile(ctx, bucket, fmt.Sprintf("photos/%s.jpg", basename), int64(imgBuffer.Len()), &imgBuffer, metadata); err != nil {
		return fmt.Errorf("failed to copy processed image to bucket '%s': %v", bucket, err)
	}

	return nil
}

func createThumbnail(ctx context.Context, repo MinioRepository, bucket, filename string, r io.ReadSeeker) error {
	var imgThumbnailBuffer bytes.Buffer

	if err := image.CreateThumbnail(r, &imgThumbnailBuffer); err != nil {
		return fmt.Errorf("failed to create thumbnail for image: %v", err)
	}

	_, basename := path.Split(filename)

	emptyMetadata := make(map[string]string)

	if err := repo.PutFile(ctx, bucket, fmt.Sprintf("thumbnail/%s", basename), int64(imgThumbnailBuffer.Len()), &imgThumbnailBuffer, emptyMetadata); err != nil {
		return fmt.Errorf("failed to copy thumbnail image to bucket '%s': %v", bucket, err)
	}

	return nil
}
