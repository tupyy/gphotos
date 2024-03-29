package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	miniotags "github.com/minio/minio-go/v7/pkg/tags"
	"github.com/tupyy/gophoto/internal/entity"
	"go.uber.org/zap"
)

// this format depends on exif extract library
const (
	dateFormat       = "2006:01:02 15:04:05"
	photoContentType = "application/jpg"
	dateKey          = "X-Amz-Meta-Date"
)

type MinioRepo struct {
	client *minio.Client
}

func New(client *minio.Client) *MinioRepo {
	return &MinioRepo{client}
}

func (m *MinioRepo) CreateBucket(ctx context.Context, bucket string, tags map[string]string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w failed to create bucket %s on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	if exists {
		return fmt.Errorf("bucket already exists on endpoint %s", m.client.EndpointURL())
	}

	err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("%w failed to create bucket %s on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	return m.SetBucketTagging(ctx, bucket, tags)
}

func (m *MinioRepo) DeleteBucket(ctx context.Context, bucket string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w failed to delete bucket %s on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	if !exists {
		return nil
	}

	// remove all objects
	objectsCh := make(chan minio.ObjectInfo)
	doneCh := make(chan interface{}, 1)
	errCh := make(chan error)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
			if object.Err != nil {
				zap.S().Errorw("failed to list bucket", "bucket", bucket, "error", object.Err)

				errCh <- object.Err
				break
			}
			objectsCh <- object

			select {
			case <-doneCh:
				break
			default:
			}
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for {
		rerr, more := <-m.client.RemoveObjects(ctx, bucket, objectsCh, opts)
		if !more {
			break
		}

		if rerr.Err != nil {
			doneCh <- struct{}{}
			return rerr.Err
		}

		select {
		case e := <-errCh:
			return e
		default:
		}
	}

	err = m.client.RemoveBucket(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to delete bucket %s on endpoint %s: %+v", bucket, m.client.EndpointURL(), err)
	}

	return nil
}

func (m *MinioRepo) PutFile(ctx context.Context, bucket, filename string, size int64, r io.Reader, metadata map[string]string) error {
	if len(bucket) == 0 || len(filename) == 0 {
		return errors.New("failed to upload file to minio. bucket or filename missing.")
	}

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to upload file %s to bucket %s on endpoint %s: %+v", filename, bucket, m.client.EndpointURL(), err)
	}

	if !exists {
		return fmt.Errorf("failed to upload file %s to bucket %s on endpoint %s: %+v", filename, bucket, m.client.EndpointURL(), err)
	}

	_, err = m.client.PutObject(ctx, bucket, filename, r, size, minio.PutObjectOptions{ContentType: photoContentType, UserMetadata: metadata})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to bucket %s on endpoint %s: %+v", filename, bucket, m.client.EndpointURL(), err)
	}

	return nil
}

func (m *MinioRepo) GetFile(ctx context.Context, bucket, filename string) (io.ReadSeeker, map[string]string, error) {
	if len(bucket) == 0 || len(filename) == 0 {
		return nil, nil, errors.New("failed to get file. bucket or filename missing.")
	}

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, nil, fmt.Errorf("%w internal error on endpoint %s", err, m.client.EndpointURL())
	}

	if !exists {
		return nil, nil, fmt.Errorf("%w bucket %s does not exists on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	r, err := m.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("%w failed to read file '%s/%s'", err, bucket, filename)
	}

	objectInfo, err := r.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("%w failed to stat file '%s/%s'", err, bucket, filename)
	}

	return r, objectInfo.UserMetadata, nil
}

func (m *MinioRepo) DeleteFile(ctx context.Context, bucket, filename string) error {
	if len(bucket) == 0 || len(filename) == 0 {
		return errors.New("failed to get file. bucket or filename missing.")
	}

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w internal error on endpoint %s", err, m.client.EndpointURL())
	}

	if !exists {
		return fmt.Errorf("%w bucket %s does not exists on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	err = m.client.RemoveObject(ctx, bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%w failed to remove file '%s/%s'", err, bucket, filename)
	}

	return nil
}

func (m *MinioRepo) ListBucket(ctx context.Context, bucket string) ([]entity.Media, error) {
	medias := make([]entity.Media, 0, 100)

	if len(bucket) == 0 {
		return medias, errors.New("bucket missing")
	}

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return medias, fmt.Errorf("%w internal error on endpoint %s", err, m.client.EndpointURL())
	}

	if !exists {
		return medias, fmt.Errorf("%w bucket %s does not exists on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	objectCh := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive:    true,
		WithMetadata: true,
	})

	mediaMap := make(map[string]entity.Media)
	thumbnailMap := make(map[string]string)

	for object := range objectCh {
		if object.Err != nil {
			return medias, fmt.Errorf("[%w] failed to list bucket '%s'", object.Err, bucket)
		}

		if isThumbnail(object) {
			thumbnailMap[filename(object.Key)] = object.Key
		} else {
			mediaMap[filename(object.Key)] = toEntity(object, bucket)
		}
	}

	for k, v := range mediaMap {
		if vv, found := thumbnailMap[k]; found {
			v.Thumbnail = vv
		}

		medias = append(medias, v)
	}

	return medias, nil
}

func (m *MinioRepo) SetBucketTagging(ctx context.Context, bucket string, tags map[string]string) error {
	// Create tags from a map.
	bucketTags, err := miniotags.NewTags(tags, false)
	if err != nil {
		return err
	}

	err = m.client.SetBucketTagging(ctx, bucket, bucketTags)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioRepo) GetBucketTagging(ctx context.Context, bucket string) (map[string]string, error) {
	// get present tags
	bucketTags, err := m.client.GetBucketTagging(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get present tags from bucket '%s': %w", bucket, err)
	}
	return bucketTags.ToMap(), nil
}

func toEntity(o minio.ObjectInfo, bucket string) entity.Media {
	e := entity.Media{
		Filename: o.Key,
		Bucket:   bucket,
		Metadata: o.UserMetadata,
	}

	if createTime, found := o.UserMetadata[dateKey]; found {
		t, err := time.Parse(dateFormat, createTime)
		if err != nil {
			zap.S().Errorw("failed to parse time from metadata", "error", err, "time", createTime)
		} else {
			e.CreateDate = t
		}
	}

	if strings.Index(o.Key, "jpg") > 0 {
		e.MediaType = entity.Photo
	} else {
		e.MediaType = entity.Unknown
	}

	return e
}

func filename(objFilename string) string {
	hasFolder := strings.HasPrefix(objFilename, "thumbnail") || strings.HasPrefix(objFilename, "photos")

	filename := objFilename
	if hasFolder {
		parts := strings.Split(objFilename, "/")
		filename = parts[1]
	}

	return filename
}

func isThumbnail(o minio.ObjectInfo) bool {
	return strings.HasPrefix(o.Key, "thumbnail")
}
