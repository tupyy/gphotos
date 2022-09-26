package keycloak

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/tupyy/gophoto/internal/domain"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/utils/logutil"
)

const (
	allUsersKey = "allUsersKey"
	allGroupKey = "allGroupKey"
)

// keycloakCacheRepo implements decorator pattern to provide cache at repo level.
type keycloakCacheRepo struct {
	repo  domain.KeycloakRepo
	cache *gocache.Cache
}

func NewCacheRepo(r domain.KeycloakRepo, ttl time.Duration, cleanInterval time.Duration) domain.KeycloakRepo {
	return &keycloakCacheRepo{
		repo:  r,
		cache: gocache.New(ttl, cleanInterval),
	}
}

func (r keycloakCacheRepo) GetUsers(ctx context.Context, filters userFilters.Filters) ([]entity.User, error) {
	var users []entity.User

	cacheKey := generateCacheKey(allUsersKey, filters)

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		users, err = r.repo.GetUsers(ctx, filters)
		if err != nil {
			return []entity.User{}, err
		}

		if len(users) == 0 {
			return users, nil
		}

		// set cache
		r.cache.Set(cacheKey, users, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count users", len(users)).WithField("cache key", cacheKey).Debug("users cached")
	} else {
		logutil.GetDefaultLogger().WithField("count users", len(users)).WithField("cache key", cacheKey).Debug("served users from cache")
		users, _ = items.([]entity.User)
	}

	return users, nil
}

func (r keycloakCacheRepo) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	item, found := r.cache.Get(id)
	if !found {
		user, err := r.repo.GetUserByID(ctx, id)
		if err != nil {
			return entity.User{}, err
		}

		// cache album
		r.cache.Set(string(id), user, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("id", id).Debug("user cached")

		return user, nil
	}

	logutil.GetDefaultLogger().WithField("id", id).Debug("user served from cached")

	return item.(entity.User), nil
}

func (r keycloakCacheRepo) GetGroups(ctx context.Context) ([]entity.Group, error) {
	item, found := r.cache.Get(allGroupKey)
	if !found {
		groups, err := r.repo.GetGroups(ctx)
		if err != nil {
			return []entity.Group{}, err
		}

		// cache album
		r.cache.Set(allGroupKey, groups, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().Debug("groups cached")

		return groups, nil
	}

	logutil.GetDefaultLogger().Debug("groups served from cached")

	return item.([]entity.Group), nil
}

func generateCacheKey(initialKey string, filters userFilters.Filters) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s", initialKey)
	for k := range filters {
		fmt.Fprintf(&sb, "%d", k)
	}

	h := base64.StdEncoding.EncodeToString([]byte(sb.String()))

	return string(h)
}
