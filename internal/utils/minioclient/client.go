package minioclient

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tupyy/gophoto/internal/conf"
)

func New(c conf.MinioConfig) (*minio.Client, error) {
	// Initialize minio client object.
	return minio.New(c.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessID, c.AccessSecretKey, ""),
		Secure: false,
	})
}
