package conf

import "github.com/spf13/viper"

const (
	// params for minio
	minioUrl             = "MINIO_SERVER_URL"
	minioUser            = "MINIO_ACCESS_ID"
	minioPwd             = "MINIO_ACCESS_KEY"
	minioTemporaryBucket = "MINIO_TEMP_BUCKET"
)

type MinioConfig struct {
	Url      string
	User     string
	Password string
}

func GetMinioConfig() MinioConfig {
	return MinioConfig{
		Url:      viper.GetString(minioUrl),
		User:     viper.GetString(minioUser),
		Password: viper.GetString(minioPwd),
	}
}

func GetMinioTemporaryBucket() string {
	return viper.GetString(minioTemporaryBucket)
}
