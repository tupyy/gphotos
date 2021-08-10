package album

import (
	"context"
	"time"

	gocache "github.com/patrickmn/go-cache"
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

	// save to cache
	r.cache.Set(string(id), album, gocache.DefaultExpiration)
	logutil.GetDefaultLogger().WithField("id", id).Debug("album cached")

	return id, nil
}

func (r albumCacheRepo) Update(ctx context.Context, album entity.Album) error {
	err := r.repo.Update(ctx, album)
	if err != nil {
		return err
	}

	// save to cache
	r.cache.Set(string(album.ID), album, gocache.DefaultExpiration)
	logutil.GetDefaultLogger().WithField("id", album.ID).Debug("album cached")

	return nil
}

func (r albumCacheRepo) Delete(ctx context.Context, id int32) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// remove it from cache
	r.cache.Delete(string(id))
	logutil.GetDefaultLogger().WithField("id", id).Debug("delete album from cache")

	return nil
}

func (r albumCacheRepo) Get(ctx context.Context, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	items, found := r.cache.Get(allAlbumsKey)
	if !found {
		ent, err := r.repo.Get(ctx, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(allAlbumsKey, ent, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(ent)).Debug("albums cached")

		// sort
		sorter.Sort(ent)

		//filter them
		filteredAlbums := make([]entity.Album, 0, len(ent))
		for _, filter := range filters {
			for _, e := range ent {
				if filter(e) {
					filteredAlbums = append(filteredAlbums, e)
				}
			}
		}

		return filteredAlbums, nil
	}

	albums, _ := items.([]entity.Album)

	// sort
	sorter.Sort(albums)

	//filter them
	filteredAlbums := make([]entity.Album, 0, len(albums))
	for _, filter := range filters {
		for _, e := range albums {
			if filter(e) {
				filteredAlbums = append(filteredAlbums, e)
			}
		}
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

	return filteredAlbums, nil

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
	items, found := r.cache.Get(ownerID)
	if !found {
		ent, err := r.repo.GetByOwnerID(ctx, ownerID, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(ownerID, ent, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(ent)).Debug("albums cached")

		// sort
		sorter.Sort(ent)

		//filter them
		filteredAlbums := make([]entity.Album, 0, len(ent))
		for _, filter := range filters {
			for _, e := range ent {
				if filter(e) {
					filteredAlbums = append(filteredAlbums, e)
				}
			}
		}

		return filteredAlbums, nil
	}

	albums, _ := items.([]entity.Album)

	// sort
	sorter.Sort(albums)

	//filter them
	filteredAlbums := make([]entity.Album, 0, len(albums))
	for _, filter := range filters {
		for _, e := range albums {
			if filter(e) {
				filteredAlbums = append(filteredAlbums, e)
			}
		}
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

	return filteredAlbums, nil
}

func (r albumCacheRepo) GetByUserID(ctx context.Context, userID string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	items, found := r.cache.Get(userID)
	if !found {
		ent, err := r.repo.GetByUserID(ctx, userID, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(userID, ent, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(ent)).Debug("albums cached")

		// sort
		sorter.Sort(ent)

		//filter them
		filteredAlbums := make([]entity.Album, 0, len(ent))
		for _, filter := range filters {
			for _, e := range ent {
				if filter(e) {
					filteredAlbums = append(filteredAlbums, e)
				}
			}
		}

		return filteredAlbums, nil
	}

	albums, _ := items.([]entity.Album)

	// sort
	sorter.Sort(albums)

	//filter them
	filteredAlbums := make([]entity.Album, 0, len(albums))
	for _, filter := range filters {
		for _, e := range albums {
			if filter(e) {
				filteredAlbums = append(filteredAlbums, e)
			}
		}
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

	return filteredAlbums, nil
}

func (r albumCacheRepo) GetByGroupName(ctx context.Context, groupName string, sorter sort.AlbumSorter, filters ...filters.AlbumFilter) ([]entity.Album, error) {
	items, found := r.cache.Get(groupName)
	if !found {
		ent, err := r.repo.GetByGroupName(ctx, groupName, sorter, filters...)
		if err != nil {
			return []entity.Album{}, err
		}

		// set cache
		r.cache.Set(groupName, ent, gocache.DefaultExpiration)
		logutil.GetDefaultLogger().WithField("count albums", len(ent)).Debug("albums cached")

		// sort
		sorter.Sort(ent)

		//filter them
		filteredAlbums := make([]entity.Album, 0, len(ent))
		for _, filter := range filters {
			for _, e := range ent {
				if filter(e) {
					filteredAlbums = append(filteredAlbums, e)
				}
			}
		}

		return filteredAlbums, nil
	}

	albums, _ := items.([]entity.Album)

	// sort
	sorter.Sort(albums)

	//filter them
	filteredAlbums := make([]entity.Album, 0, len(albums))
	for _, filter := range filters {
		for _, e := range albums {
			if filter(e) {
				filteredAlbums = append(filteredAlbums, e)
			}
		}
	}

	logutil.GetDefaultLogger().WithField("count albums", len(albums)).Debug("served album from cache")

	return filteredAlbums, nil
}
