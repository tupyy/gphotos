package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services/image"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

// Store describe photo store operations
type MediaRepository interface {
	// GetFile returns a reader to file.
	GetFile(ctx context.Context, bucket, filename string) (io.ReadSeeker, map[string]string, error)
	// PutFile save a file to a bucket.
	PutFile(ctx context.Context, bucket, filename string, size int64, r io.Reader, metadata map[string]string) error
	// ListFiles list the content of a bucket
	ListBucket(ctx context.Context, bucket string) ([]entity.Media, error)
	// DeleteFile deletes a file from a bucket.
	DeleteFile(ctx context.Context, bucket, filename string) error
	// CreateBucket create a bucket.
	CreateBucket(ctx context.Context, bucket string) error
	// DeleteBucket removes bucket.
	DeleteBucket(ctx context.Context, bucket string) error
}

type MediaType int

const (
	Photo MediaType = iota
	Video
)

type Service struct {
	repo MediaRepository
}

func New(repo MediaRepository) *Service {
	return &Service{repo}
}

func (s *Service) CreateBucket(ctx context.Context, bucket string) error {
	return s.repo.CreateBucket(ctx, bucket)
}

// TODO move logic from repo to here
func (s *Service) DeleteBucket(ctx context.Context, bucket string) error {
	return s.repo.DeleteBucket(ctx, bucket)
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
				logutil.GetLogger(ctx).WithError(err).WithField("filename", m.Filename).Error("failed to get photo from repo")

				continue
			}

			if err := createThumbnail(ctx, s.repo, m.Bucket, m.Filename, r); err != nil {
				logutil.GetLogger(ctx).WithError(err).WithField("filename", m.Filename).Error("failed to get create thumbnail")

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

func processPhoto(ctx context.Context, repo MediaRepository, bucket, filename string, r io.ReadSeeker) error {
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

func createThumbnail(ctx context.Context, repo MediaRepository, bucket, filename string, r io.ReadSeeker) error {
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
