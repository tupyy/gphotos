package minio

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
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
