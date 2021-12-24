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

func (s *Service) Create(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	id, err := s.repo.Create(ctx, tag)
	if err != nil {
		logutil.GetLogger(ctx).WithField("tag", tag.String()).Error(err)

		return entity.Tag{}, fmt.Errorf("failed to create tag: %+v", err)
	}

	tag.ID = id

	return tag, nil
}

func (s *Service) CreateAndAssociate(ctx context.Context, tag entity.Tag, albumID int32) error {
	id, err := s.repo.Create(ctx, tag)
	if err != nil {
		logutil.GetLogger(ctx).WithField("tag", tag.String()).Error(err)

		return fmt.Errorf("failed to create tag: %+v", err)
	}

	if err := s.repo.Associate(ctx, albumID, id); err != nil {
		logutil.GetLogger(ctx).WithFields(logrus.Fields{
			"tag":      tag.String(),
			"album_id": albumID,
		}).WithError(err).Error("failed to associate album with tag")

		return fmt.Errorf("failed to create tag: %+v", err)
	}

	return nil
}
