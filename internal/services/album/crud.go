package album

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services"
	"github.com/tupyy/gophoto/internal/services/media"
)

// AlbumRepository is the interface to be implemented by transport layer.
type AlbumRepository interface {
	// Create creates an album.
	Create(ctx context.Context, album entity.Album) (albumID int32, err error)
	// Update an album.
	Update(ctx context.Context, album entity.Album) error
	// Delete removes an album from postgres.
	Delete(ctx context.Context, id int32) error
	// Get return all the albums.
	Get(ctx context.Context) ([]entity.Album, error)
	// GetByID return an album by id.
	GetByID(ctx context.Context, id int32) (entity.Album, error)
	// GetByOwner return all albums of a user for which he is the owner.
	GetByOwnerID(ctx context.Context, ownerID string) ([]entity.Album, error)
	// GetByUserID returns a list of albums for which the user has at least one permission set.
	GetByUserID(ctx context.Context, userID string) ([]entity.Album, error)
	// GetByGroup returns a list of albums for which the group has at least one permission.
	GetByGroupName(ctx context.Context, groupName string) ([]entity.Album, error)
	// GetByGroups returns a list of albums with at least one persmission for at least on group in the list.
	GetByGroups(ctx context.Context, groups []string) ([]entity.Album, error)
}

const forbittenChar = "_$"

type Service struct {
	albumRepo    AlbumRepository
	mediaService *media.Service
}

func New(albumRepo AlbumRepository, media *media.Service) *Service {
	return &Service{albumRepo, media}
}

func (s *Service) Create(ctx context.Context, newAlbum entity.Album) (int32, error) {
	// generate bucket name
	bucketID := strings.ReplaceAll(uuid.New().String(), "-", "")
	n := strings.ReplaceAll(strings.ToLower(newAlbum.Name), " ", "-")

	for _, s := range forbittenChar {
		n = strings.ReplaceAll(n, string(s), "-")
	}

	newAlbum.Bucket = fmt.Sprintf("%s-%s", n, bucketID[:8])

	// create the bucket
	if err := s.mediaService.CreateBucket(ctx, newAlbum.Bucket); err != nil {
		return 0, fmt.Errorf("%w '%s': %v", services.ErrCreateBucket, newAlbum.Name, err)
	}

	albumID, err := s.albumRepo.Create(ctx, newAlbum)
	if err != nil {
		return 0, fmt.Errorf("%w '%s': %v", services.ErrCreateAlbum, newAlbum.Name, err)
	}

	return albumID, nil
}

func (s *Service) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	err := s.albumRepo.Update(ctx, album)
	if err != nil {
		return album, fmt.Errorf("%w '%d': %v", services.ErrUpdateAlbum, album.ID, err)
	}

	return album, nil
}

func (s *Service) Delete(ctx context.Context, album entity.Album) error {
	err := s.mediaService.DeleteBucket(ctx, album.Bucket)
	if err != nil {
		return fmt.Errorf("%w '%s': %v", services.ErrDeleteBucket, album.Bucket, err)
	}

	err = s.albumRepo.Delete(ctx, album.ID)
	if err != nil {
		return fmt.Errorf("%w '%d': %v", services.ErrDeleteAlbum, album.ID, err)
	}

	return nil
}
