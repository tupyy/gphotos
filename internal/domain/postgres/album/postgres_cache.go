package album

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/filters"
	"github.com/tupyy/gophoto/internal/domain/sort"
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

func (r albumCacheRepo) Get(ctx context.Context, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	var albums []entity.Album

	items, found := r.cache.Get(allAlbumsKey)
	if !found {
		var err error

		albums, err = r.repo.Get(ctx, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(allAlbumsKey, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("albums cached")
	} else {
		albums, _ = items.([]entity.Album)
	}

	// sort
	if sorter != nil {
		sorter.Sort(albums)
	}

	//filter them
	if len(filters) > 0 {
		filteredAlbums := filterAlbums(filters, albums)
		logutil.GetDefaultLogger().WithFields(logrus.Fields{
			"count before filter": len(albums),
			"count after filter":  len(filteredAlbums),
		}).Debug("served album from cache")

		return filteredAlbums, nil
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

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

func (r albumCacheRepo) GetByOwnerID(ctx context.Context, ownerID string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	cacheKey := fmt.Sprintf("owner%s", ownerID)

	var albums []entity.Album

	items, found := r.cache.Get(cacheKey)
	if !found {
		var err error

		albums, err = r.repo.GetByOwnerID(ctx, ownerID, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(cacheKey, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("albums cached")
	} else {
		albums, _ = items.([]entity.Album)
	}

	// sort
	if sorter != nil {
		sorter.Sort(albums)
	}

	//filter them
	if len(filters) > 0 {
		filteredAlbums := filterAlbums(filters, albums)
		logutil.GetDefaultLogger().WithFields(logrus.Fields{
			"count before filter": len(albums),
			"count after filter":  len(filteredAlbums),
			"owner id":            ownerID,
		}).Debug("filtered albums")

		return filteredAlbums, nil
	}

	return albums, nil
}

func (r albumCacheRepo) GetByUserID(ctx context.Context, userID string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	var albums []entity.Album

	items, found := r.cache.Get(userID)
	if !found {
		var err error

		albums, err = r.repo.GetByUserID(ctx, userID, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(userID, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("albums cached")
	} else {
		albums, _ = items.([]entity.Album)
	}

	// sort
	if sorter != nil {
		sorter.Sort(albums)
	}

	//filter them
	if len(filters) > 0 {
		filteredAlbums := filterAlbums(filters, albums)
		logutil.GetDefaultLogger().WithFields(logrus.Fields{
			"count before filter": len(albums),
			"count after filter":  len(filteredAlbums),
			"user id":             userID,
		}).Debug("albums filtered")

		return filteredAlbums, nil
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

	return albums, nil
}

func (r albumCacheRepo) GetByGroupName(ctx context.Context, groupName string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	var albums []entity.Album

	items, found := r.cache.Get(groupName)
	if !found {
		var err error

		albums, err = r.repo.GetByGroupName(ctx, groupName, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(groupName, albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("albums cached")
	} else {
		albums, _ = items.([]entity.Album)
	}

	// sort
	if sorter != nil {
		sorter.Sort(albums)
	}

	//filter them
	if len(filters) > 0 {
		filteredAlbums := filterAlbums(filters, albums)
		logutil.GetDefaultLogger().WithFields(logrus.Fields{
			"count before filter": len(albums),
			"count after filter":  len(filteredAlbums),
			"group name":          groupName,
		}).Debug("albums filtered")

		return filteredAlbums, nil
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

	return albums, nil
}

func (r albumCacheRepo) GetByGroups(ctx context.Context, groupNames []string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	var albums []entity.Album

	items, found := r.cache.Get(strings.Join(groupNames, "#"))
	if !found {
		var err error

		albums, err = r.repo.GetByGroups(ctx, groupNames, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(strings.Join(groupNames, "#"), albums, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("albums cached")
	} else {
		albums, _ = items.([]entity.Album)
	}

	// sort
	if sorter != nil {
		sorter.Sort(albums)
	}

	//filter them
	if len(filters) > 0 {
		filteredAlbums := filterAlbums(filters, albums)
		logutil.GetDefaultLogger().WithFields(logrus.Fields{
			"count before filter": len(albums),
			"count after filter":  len(filteredAlbums),
			"group names":         groupNames,
		}).Debug("albums filtered")

		return filteredAlbums, nil
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served albums from cache")

	return albums, nil
}

func filterAlbums(filters []filters.AlbumFilter, albums []entity.Album) []entity.Album {
	filteredAlbums := make([]entity.Album, 0, len(albums))
	for _, a := range albums {
		pass := true
		for _, filter := range filters {
			if !filter(a) {
				logutil.GetDefaultLogger().WithFields(logrus.Fields{
					"filter":   reflect.TypeOf(filter).String(),
					"album_id": a.ID,
				}).Trace("filter did not passed")

				pass = false
				break
			}

			logutil.GetDefaultLogger().WithFields(logrus.Fields{
				"filter":   reflect.TypeOf(filter).Name(),
				"album_id": a.ID,
			}).Trace("filter passed")
		}

		if pass {
			filteredAlbums = append(filteredAlbums, a)
		}
	}

	return filteredAlbums
}
