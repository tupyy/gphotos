package keycloak

import (
	"context"

	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/filters/user"
)

type Service struct {
	repos domain.Repositories
}

type Query struct {
	predicates   []Predicate
	keycloakRepo domain.KeycloakRepo
}

func New(repos domain.Repositories) *Service {
	return &Service{repos}
}

func (s *Service) Query() *Query {
	return &Query{
		predicates:   []Predicate{},
		keycloakRepo: s.repos[domain.KeycloakRepoName].(domain.KeycloakRepo),
	}
}

func (q *Query) Where(p Predicate) *Query {
	q.predicates = append(q.predicates, p)

	return q
}

func (q *Query) AllUsers(ctx context.Context) ([]entity.User, error) {
	filters := make([]user.Filter, 0, len(q.predicates))
	for _, p := range q.predicates {
		filters = append(filters, p())
	}

	users, err := q.keycloakRepo.GetUsers(ctx, filters)
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil

}

func (q *Query) FirstUser(ctx context.Context, id string) (entity.User, error) {
	user, err := q.keycloakRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (q *Query) AllGroups(ctx context.Context) ([]entity.Group, error) {
	groups, err := q.keycloakRepo.GetGroups(ctx)
	if err != nil {
		return []entity.Group{}, err
	}

	return groups, nil
}
