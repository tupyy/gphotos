package users

import (
	"context"

	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repos/filters/user"
	userFilters "github.com/tupyy/gophoto/internal/repos/filters/user"
)

type KeycloakRepository interface {
	// Get returns all the users.
	GetUsers(ctx context.Context, filters userFilters.Filters) ([]entity.User, error)
	// GetByID return the user by id.
	GetUserByID(ctx context.Context, id string) (entity.User, error)
	// GetGroups returns all groups.
	GetGroups(ctx context.Context) ([]entity.Group, error)
}

// Postgres repo to handler relationships between users
type UserRespository interface {
	GetRelatedUsers(ctx context.Context, user entity.User) (ids []string, err error)
}

type Service struct {
	keycloakRepo KeycloakRepository
	userRepo     UserRespository
}

type Query struct {
	predicates   []Predicate
	keycloakRepo KeycloakRepository
	userRepo     UserRespository
}

func New(keycloakRepo KeycloakRepository, userRepo UserRespository) *Service {
	return &Service{
		keycloakRepo: keycloakRepo,
		userRepo:     userRepo,
	}
}

func (s *Service) Query() *Query {
	return &Query{
		predicates:   []Predicate{},
		keycloakRepo: s.keycloakRepo,
		userRepo:     s.userRepo,
	}
}

func (q *Query) Where(p Predicate) *Query {
	q.predicates = append(q.predicates, p)

	return q
}

func (q *Query) All(ctx context.Context) ([]entity.User, error) {
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

func (q *Query) AllRelatedUsers(ctx context.Context, u entity.User) ([]entity.User, error) {
	// get the ids of related users
	ids, err := q.userRepo.GetRelatedUsers(ctx, u)
	if err != nil {
		return []entity.User{}, err
	}

	filters := make([]user.Filter, 0, len(q.predicates))
	for _, p := range q.predicates {
		filters = append(filters, p())
	}

	// get all the users from keycloak
	users, err := q.keycloakRepo.GetUsers(ctx, filters)
	if err != nil {
		return []entity.User{}, err
	}

	relatedUsers := make([]entity.User, 0, len(ids))

	// remove users which are not relevant for albums found.
	addedUsers := make(map[string]interface{})
	for _, id := range ids {
		for _, u := range users {
			_, alreadyAdded := addedUsers[u.ID]

			if u.ID == id && !alreadyAdded {
				relatedUsers = append(relatedUsers, u)
				addedUsers[u.ID] = true
			}
		}
	}

	return relatedUsers, err
}

func (q *Query) First(ctx context.Context, id string) (entity.User, error) {
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
