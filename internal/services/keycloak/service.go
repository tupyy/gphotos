package keycloak

import (
	"context"
	"fmt"

	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

type Service struct {
	repos domain.Repositories
}

func New(repos domain.Repositories) *Service {
	return &Service{repos}
}

func (u *Service) Get(ctx context.Context, id string) (entity.User, error) {
	keycloakRepo := u.repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	logger := logutil.GetLogger(ctx)

	user, err := keycloakRepo.GetUserByID(ctx, id)
	if err != nil {
		logger.WithError(err).WithField("user id", id).Error("failed to get user")

		return entity.User{}, fmt.Errorf("[%w] failed to get user '%s'", err, id)
	}

	// TODO get user's groups

	return user, nil
}

func (u *Service) GetUsers(ctx context.Context) ([]entity.User, error) {
	keycloakRepo := u.repos[domain.KeycloakRepoName].(domain.KeycloakRepo)

	logger := logutil.GetLogger(ctx)

	users, err := keycloakRepo.GetUsers(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("failed to get users")

		return []entity.User{}, err
	}

	// TODO get users groups

	return users, nil
}
