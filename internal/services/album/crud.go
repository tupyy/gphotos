package album

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services"
	"github.com/tupyy/gophoto/internal/services/media"
)

// AlbumRepository is the interface to be implemented by transport layer.
type AlbumRepository interface {
	// Create creates an album.
	Create(ctx context.Context, album entity.Album) (entity.Album, error)
	// Update an album.
	Update(ctx context.Context, album entity.Album) (entity.Album, error)
	// Delete removes an album from postgres.
	Delete(ctx context.Context, id string) error
	// Get return all the albums.
	Get(ctx context.Context) ([]entity.Album, error)
	// Set permissions for the album
	SetPermissions(ctx context.Context, albumId string, permissions []entity.AlbumPermission) error
	// remove permissions of ownerID for the album
	RemovePermissions(ctx context.Context, albumId string) error
	// GetByID return an album by id.
	GetByID(ctx context.Context, id string) (entity.Album, error)
	// GetByOwner return all albums of a user for which he is the owner.
	GetByOwner(ctx context.Context, owner string) ([]entity.Album, error)
	// GetByUserID returns a list of albums for which the user has at least one permission set.
	GetByUser(ctx context.Context, userName string) ([]entity.Album, error)
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

func (s *Service) Create(ctx context.Context, newAlbum entity.Album) (entity.Album, error) {
	// generate bucket name
	bucketID := strings.ReplaceAll(uuid.New().String(), "-", "")
	n := strings.ReplaceAll(strings.ToLower(newAlbum.Name), " ", "-")

	for _, s := range forbittenChar {
		n = strings.ReplaceAll(n, string(s), "-")
	}

	newAlbum.Bucket = fmt.Sprintf("%s-%s", n, bucketID[:8])

	tags := map[string]string{
		"album/name":     newAlbum.Name,
		"album/date":     newAlbum.CreatedAt.Format(time.RFC3339),
		"owner/username": newAlbum.Owner,
	}
	// create the bucket
	if err := s.mediaService.CreateBucket(ctx, newAlbum.Bucket, tags); err != nil {
		return entity.Album{}, fmt.Errorf("%w '%s': %v", services.ErrCreateBucket, newAlbum.Name, err)
	}

	album, err := s.albumRepo.Create(ctx, newAlbum)
	if err != nil {
		return entity.Album{}, fmt.Errorf("%w '%s': %v", services.ErrCreateAlbum, newAlbum.Name, err)
	}

	return album, nil
}

func (s *Service) Update(ctx context.Context, album entity.Album) (entity.Album, error) {
	album, err := s.albumRepo.Update(ctx, album)
	if err != nil {
		return album, fmt.Errorf("%w '%s': %v", services.ErrUpdateAlbum, album.ID, err)
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
		return fmt.Errorf("%w '%s': %v", services.ErrDeleteAlbum, album.ID, err)
	}

	return nil
}

func (s *Service) SetPermissions(ctx context.Context, album entity.Album, permissions []entity.AlbumPermission) error {
	// remove old permissions
	err := s.albumRepo.RemovePermissions(ctx, album.ID)
	if err != nil {
		return err
	}
	if err := s.albumRepo.SetPermissions(ctx, album.ID, permissions); err != nil {
		return err
	}
	return nil
}
