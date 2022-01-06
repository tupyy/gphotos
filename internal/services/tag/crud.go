package tag

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

type Service struct {
	repo domain.Tag
}

func New(r domain.Tag) *Service {
	return &Service{repo: r}
}

func (s *Service) Get(ctx context.Context, userID string) ([]entity.Tag, error) {
	return s.repo.GetByUser(ctx, userID)
}

func (s *Service) Create(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	id, err := s.repo.Create(ctx, tag)
	if err != nil {
		return entity.Tag{}, fmt.Errorf("create tag: %+v", err)
	}

	tag.ID = id

	return tag, nil
}

func (s *Service) CreateAndAssociate(ctx context.Context, tag entity.Tag, albumID int32) error {
	id, err := s.repo.Create(ctx, tag)
	if err != nil {
		return fmt.Errorf("create tag: %+v", err)
	}

	if err := s.repo.Associate(ctx, albumID, id); err != nil {
		return fmt.Errorf("failed to create tag: %+v", err)
	}

	return nil
}

func (s *Service) Dissociate(ctx context.Context, tag entity.Tag, albumID int32) error {
	if err := s.repo.Dissociate(ctx, albumID, tag.ID); err != nil {
		return fmt.Errorf("dissociate tag from album: %+v", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, tag entity.Tag) error {
	// dissociate from all albums
	for _, albumID := range tag.Albums {
		if err := s.Dissociate(ctx, tag, albumID); err != nil {
			logutil.GetLogger(ctx).WithFields(logrus.Fields{
				"tag":      tag.String(),
				"album_id": albumID,
			}).WithError(err).Error("dissociate tag from album")
		}
	}

	return s.repo.Delete(ctx, tag.ID)
}
