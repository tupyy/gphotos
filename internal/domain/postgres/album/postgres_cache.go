package album

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/tupyy/gophoto/internal/domain"
	albumFilters "github.com/tupyy/gophoto/internal/domain/filters/album"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/logutil"
)

const (
	allAlbumsKey = "allAlbumsKey"
)

// albumCacheRepo implements decorator pattern to provide cache at repo level.
type albumCacheRepo struct {
	repo  domain.Album
	cache *gocache.Cache
}

func NewCacheRepo(r domain.Album, ttl time.Duration, cleanInterval time.Duration) domain.Album {
	return &albumCacheRepo{
		repo:  r,
		cache: gocache.New(ttl, cleanInterval),
	}
}

func (r albumCacheRepo) Create(ctx context.Context, album entity.Album) (albumID int32, err error) {
	id, err := r.repo.Create(ctx, album)
	if err != nil {
		return -1, err
	}

	// clear cache
	r.cache.Flush()
	logutil.GetDefaultLogger().Debug("cache flushed after create")

	return id, nil
}

func (r albumCacheRepo) Update(ctx context.Context, album entity.Album) error {
	err := r.repo.Update(ctx, album)
	if err != nil {
		return err
	}

	// clear cache
	r.cache.Flush()
	logutil.GetDefaultLogger().Debug("cache flushed after update")

	return nil
}

func (r albumCacheRepo) Delete(ctx context.Context, id int32) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// clear cache
	r.cache.Flush()
	logutil.GetDefaultLogger().Debug("cache flushed after delete")

	return nil
}

func (r albumCacheRepo) Get(ctx context.Context, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums []entity.Album

	// generate a cache key depending on filters
	cacheKey := generateCacheKey(allAlbumsKey, filters)

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		albums, err = r.repo.Get(ctx, filters)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(allAlbumsKey, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("albums cached")
	} else {
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("served album from cache")
		albums, _ = items.([]entity.Album)
	}

	return albums, nil

}

func (r albumCacheRepo) GetByID(ctx context.Context, id int32) (entity.Album, error) {
	item, found := r.cache.Get(string(id))
	if !found {
		album, err := r.repo.GetByID(ctx, id)
		if err != nil {
			return entity.Album{}, err
		}

		// cache album
		r.cache.Set(string(id), album, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("id", id).Debug("album cached")

		return album, nil
	}

	logutil.GetDefaultLogger().WithField("id", id).Debug("album served from cached")

	return item.(entity.Album), nil
}

func (r albumCacheRepo) GetByOwnerID(ctx context.Context, ownerID string, filters albumFilters.Filters) ([]entity.Album, error) {
	cacheKey := generateCacheKey(fmt.Sprintf("owner%s", ownerID), filters)

	var albums []entity.Album

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		albums, err = r.repo.GetByOwnerID(ctx, ownerID, filters)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(cacheKey, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("albums cached")
	} else {
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("served album from cache")
		albums, _ = items.([]entity.Album)
	}

	return albums, nil
}

func (r albumCacheRepo) GetByUserID(ctx context.Context, userID string, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums []entity.Album

	cacheKey := generateCacheKey(userID, filters)

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		albums, err = r.repo.GetByUserID(ctx, userID, filters)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(userID, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("albums cached")
	} else {
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("served album from cache")
		albums, _ = items.([]entity.Album)
	}

	return albums, nil
}

func (r albumCacheRepo) GetByGroupName(ctx context.Context, groupName string, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums []entity.Album

	cacheKey := generateCacheKey(groupName, filters)

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		albums, err = r.repo.GetByGroupName(ctx, groupName, filters)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(groupName, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("albums cached")
	} else {
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("served album from cache")
		albums, _ = items.([]entity.Album)
	}

	return albums, nil
}

func (r albumCacheRepo) GetByGroups(ctx context.Context, groupNames []string, filters albumFilters.Filters) ([]entity.Album, error) {
	var albums []entity.Album

	cacheKey := generateCacheKey(strings.Join(groupNames, "#"), filters)

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		albums, err = r.repo.GetByGroups(ctx, groupNames, filters)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(cacheKey, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("albums cached")
	} else {
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).WithField("cache key", cacheKey).Debug("served albums from cache")
		albums, _ = items.([]entity.Album)
	}

	return albums, nil
}

func generateCacheKey(initialKey string, filters albumFilters.Filters) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s", initialKey)
	for k := range filters {
		fmt.Fprintf(&sb, "%d", k)
	}

	h := base64.StdEncoding.EncodeToString([]byte(sb.String()))

	return string(h)
}
