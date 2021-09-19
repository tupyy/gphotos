package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/tupyy/gophoto/internal/domain/entity"
)

type MinioRepo struct {
	client *minio.Client
}

func New(client *minio.Client) *MinioRepo {
	return &MinioRepo{client}
}

func (m *MinioRepo) CreateBucket(ctx context.Context, bucket string) error {
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

	return nil
}

// TODO delete the content before removing the bucket
func (m *MinioRepo) DeleteBucket(ctx context.Context, bucket string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w failed to delete bucket %s on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	if !exists {
		return nil
	}

	err = m.client.RemoveBucket(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w failed to delete bucket %s on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	return nil
}

func (m *MinioRepo) PutFile(ctx context.Context, bucket, filename string, size int64, r io.Reader) error {
	if len(bucket) == 0 || len(filename) == 0 {
		return errors.New("failed to upload file to minio. bucket or filename missing.")
	}

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("%w failed to upload file %s to bucket %s on endpoint %s", err, filename, bucket, m.client.EndpointURL())
	}

	if !exists {
		return fmt.Errorf("%w failed to upload file %s to bucket %s on endpoint %s. bucket does not exists", err, filename, bucket, m.client.EndpointURL())
	}

	_, err = m.client.PutObject(ctx, bucket, filename, r, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("%w failed to upload file %s to bucket %s on endpoint %s", err, filename, bucket, m.client.EndpointURL())
	}

	return nil
}

func (m *MinioRepo) GetFile(ctx context.Context, bucket, filename string) (io.Reader, error) {
	if len(bucket) == 0 || len(filename) == 0 {
		return nil, errors.New("failed to get file. bucket or filename missing.")
	}

	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("%w internal error on endpoint %s", err, m.client.EndpointURL())
	}

	if !exists {
		return nil, fmt.Errorf("%w bucket %s does not exists on endpoint %s", err, bucket, m.client.EndpointURL())
	}

	r, err := m.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w failed to read file '%s/%s'", err, bucket, filename)
	}

	return r, nil
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
		Recursive: false,
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

func toEntity(o minio.ObjectInfo, bucket string) entity.Media {
	e := entity.Media{
		Filename: o.Key,
		Bucket:   bucket,
	}

	if strings.Index(o.Key, "jpg") > 0 {
		e.MediaType = entity.Photo
	} else {
		e.MediaType = entity.Unknown
	}

	return e
}

func filename(objFilename string) string {
	thmbIdx := strings.Index(objFilename, "thumbnail")
	if thmbIdx > 0 {
		return objFilename[:thmbIdx-1]
	}

	if strings.Index(objFilename, ".") > 0 {
		parts := strings.Split(objFilename, ".")
		return parts[0]
	}

	return objFilename
}

func isThumbnail(o minio.ObjectInfo) bool {
	return strings.Index(o.Key, "thumbnail") > 0
}
