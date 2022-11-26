package tag

import (
	"context"
	"fmt"

	"github.com/tupyy/gophoto/internal/entity"
)

type TagRepository interface {
	// Create -- create the tag.
	Create(ctx context.Context, tag entity.Tag) (string, error)
	// Update -- update the tag.
	Update(ctx context.Context, tag entity.Tag) error
	// Delete -- delete the tag. it does not cascade.
	Delete(ctx context.Context, id string) error
	// GetByUser -- fetch all user's tags
	GetByUser(ctx context.Context, userID string) ([]entity.Tag, error)
	// GetByName -- fetch the tag by name and user id.
	GetByName(ctx context.Context, userID, name string) (entity.Tag, error)
	// GetByID -- fetch the tag by id
	GetByID(ctx context.Context, userID string, id string) (entity.Tag, error)
	// GetByAlbum -- fetch all user's tag for the album
	GetByAlbum(ctx context.Context, albumID string) ([]entity.Tag, error)
	// AssociateTag -- associates a tag with an album.
	Associate(ctx context.Context, albumID, tagID string) error
	// Dissociate -- removes a tag from an album.
	Dissociate(ctx context.Context, albumID, tagID string) error
}

type Service struct {
	repo TagRepository
}

func New(r TagRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) Get(ctx context.Context, userID string) ([]entity.Tag, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *Service) GetByName(ctx context.Context, userID string, name string) (entity.Tag, error) {
	return s.repo.GetByName(ctx, userID, name)
}

func (s *Service) GetByID(ctx context.Context, userID string, tagID string) (entity.Tag, error) {
	return s.repo.GetByID(ctx, userID, tagID)
}

func (s *Service) GetByAlbum(ctx context.Context, albumID string) ([]entity.Tag, error) {
	return s.repo.GetByAlbum(ctx, albumID)
}

func (s *Service) Create(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	id, err := s.repo.Create(ctx, tag)
	if err != nil {
		return entity.Tag{}, fmt.Errorf("create tag: %+v", err)
	}

	tag.ID = id

	return tag, nil
}

func (s *Service) CreateAndAssociate(ctx context.Context, tag entity.Tag, albumID string) error {
	id, err := s.repo.Create(ctx, tag)
	if err != nil {
		return fmt.Errorf("create tag: %+v", err)
	}

	if err := s.repo.Associate(ctx, albumID, id); err != nil {
		return fmt.Errorf("failed to create tag: %+v", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, tag entity.Tag) error {
	return s.repo.Update(ctx, tag)
}

func (s *Service) Dissociate(ctx context.Context, tag entity.Tag, albumID string) error {
	if err := s.repo.Dissociate(ctx, albumID, tag.ID); err != nil {
		return fmt.Errorf("dissociate tag from album: %+v", err)
	}

	return nil
}

func (s *Service) Associate(ctx context.Context, tag entity.Tag, albumID string) error {
	if err := s.repo.Associate(ctx, albumID, tag.ID); err != nil {
		return fmt.Errorf("associate tag '%s' with album '%s': %+v", tag.ID, albumID, err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, tag entity.Tag) error {
	return s.repo.Delete(ctx, tag.ID)
}
