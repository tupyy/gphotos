package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

type Album struct {
	repos domain.Repositories
}

func NewAlbumService(repos domain.Repositories) *Album {
	return &Album{repos}
}

func (a *Album) Create(ctx context.Context, newAlbum entity.Album) (int32, error) {
	minioRepo := a.repos[domain.MinioRepoName].(domain.Store)
	albumRepo := a.repos[domain.AlbumRepoName].(domain.Album)

	logger := logutil.GetLogger(ctx)

	// generate bucket name
	bucketID := strings.ReplaceAll(uuid.New().String(), "-", "")

	// create the bucket
	if err := minioRepo.CreateBucket(ctx, bucketID); err != nil {
		logger.WithError(err).Error("failed to create bucket")

		return 0, fmt.Errorf("[%w] failed to create album '%s'", err, newAlbum.Name)
	}

	newAlbum.Bucket = bucketID

	albumID, err := albumRepo.Create(ctx, newAlbum)
	if err != nil {
		return 0, fmt.Errorf("[%w] failed to create album '%s'", err, newAlbum.Name)
	}

	return albumID, nil
}

func (a *Album) Get(ctx context.Context, id int32) (entity.Album, error) {
	albumRepo := a.repos[domain.AlbumRepoName].(domain.Album)
	minioRepo := a.repos[domain.MinioRepoName].(domain.Store)

	logger := logutil.GetLogger(ctx)

	album, err := albumRepo.GetByID(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("album id", id).Error("failed to get album")

		return entity.Album{}, fmt.Errorf("[%w] failed to get album '%d'", err, id)
	}

	// encrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", album.ID))
	if err != nil {
		logger.WithError(err).WithField("album id", id).Error("encrypt album id")

		return entity.Album{}, fmt.Errorf("[%w] failed to encrypt album id '%d'", err, id)
	}

	album.EncryptedID = encryptedID

	// replace id with this one

	medias, err := minioRepo.ListBucket(ctx, album.Bucket)
	if err != nil {
		logger.WithField("album id", album.ID).WithError(err).Error("failed to list media for album")

		return entity.Album{}, fmt.Errorf("[%w] failed to list media for album id '%d'", err, id)
	}

	// encrypt thumbnail filenames
	encryptedPhotos := make([]entity.Media, 0, len(medias))
	encryptedVideos := make([]entity.Media, 0, len(medias))
	for _, m := range medias {
		encryptedMedia, err := encryptMedia(m, gen)
		if err != nil {
			logger.WithError(err).WithField("media", fmt.Sprintf("%+v", m)).Error("failed to encrypted media")

			continue
		}

		switch m.MediaType {
		case entity.Photo:
			encryptedPhotos = append(encryptedPhotos, encryptedMedia)
		case entity.Video:
			encryptedVideos = append(encryptedVideos, encryptedMedia)
		}
	}

	album.Photos = encryptedPhotos
	album.Videos = encryptedVideos

	return album, nil
}

func (a *Album) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	albumRepo := a.repos[domain.AlbumRepoName].(domain.Album)

	logger := logutil.GetLogger(ctx)

	err := albumRepo.Update(ctx, album)
	if err != nil {
		logger.WithError(err).WithField("album id", album.ID).Error("failed to update album")

		return album, fmt.Errorf("[%w] failed to update album '%d'", err, album.ID)
	}

	return album, nil
}

func (a *Album) Delete(ctx context.Context, album entity.Album) error {
	minioRepo := a.repos[domain.MinioRepoName].(domain.Store)
	albumRepo := a.repos[domain.AlbumRepoName].(domain.Album)

	logger := logutil.GetLogger(ctx)

	err := minioRepo.DeleteBucket(ctx, album.Bucket)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"bucket":   album.Bucket,
			"album id": album.ID,
		}).WithError(err).Error("failed to remove album's bucket")

		return fmt.Errorf("[%w] failed to remove album's bucket '%s'", err, album.Bucket)
	}

	err = albumRepo.Delete(ctx, album.ID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"bucket":   album.Bucket,
			"album id": album.ID,
		}).WithError(err).Error("failed to remove album")

		return fmt.Errorf("[%w] failed to remove album '%d'", err, album.ID)
	}

	return nil
}

func encryptMedia(m entity.Media, gen *encryption.Generator) (entity.Media, error) {
	encryptedFilename, err := gen.EncryptData(m.Filename)
	if err != nil {
		return entity.Media{}, err
	}

	encryptedThumbnail, err := gen.EncryptData(m.Thumbnail)
	if err != nil {
		return entity.Media{}, err
	}

	encryptedBucket, err := gen.EncryptData(m.Bucket)
	if err != nil {
		return entity.Media{}, err
	}

	return entity.Media{
		MediaType: m.MediaType,
		Filename:  encryptedFilename,
		Bucket:    encryptedBucket,
		Thumbnail: encryptedThumbnail,
	}, nil
}
