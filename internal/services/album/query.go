package album

import (
	"context"
	"errors"
	"fmt"

	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services"
	"github.com/tupyy/gophoto/internal/services/media"
	"github.com/tupyy/gophoto/utils/logutil"
)

var (
	// AlbumSearchError means something went wrong when searching though albums
	AlbumSearchError = errors.New("album search error")
)

type Filter interface {
	Resolve(album entity.Album) (bool, error)
}

type Query struct {
	limit  int
	offset int
	// get personal albums.
	personalAlbums bool
	// get shared albums.
	sharedAlbums bool
	// filter
	filter Filter
	// album repo
	albumRepo domain.Album
	// media service
	mediaService *media.Service
	//album sorter
	sorter *albumSorter
}

func (s *Service) Query() *Query {
	return &Query{
		albumRepo:    s.albumRepo,
		mediaService: s.mediaService,
	}
}

func (q *Query) Filter(filter Filter) *Query {
	q.filter = filter

	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit

	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = offset

	return q
}

func (q *Query) OwnAlbums(b bool) *Query {
	q.personalAlbums = b

	return q
}

func (q *Query) Sort(name SortType, order SortOrder) *Query {
	as := newSorter(name, order)
	q.sorter = as

	return q
}

func (q *Query) SharedAlbums(b bool) *Query {
	q.sharedAlbums = b

	return q
}

// All returns a list of albums sliced if offset & limit are set and the total number of albums.
func (q *Query) All(ctx context.Context, user entity.User) ([]entity.Album, int, error) {
	albums := make(map[int32]entity.Album)

	if q.personalAlbums {
		// fetch personal albums
		pa, err := q.albumRepo.GetByOwnerID(ctx, user.ID)
		if err != nil {
			return []entity.Album{}, 0, fmt.Errorf("%w personal album: %v", services.ErrGetAlbums, err)
		}

		for _, a := range pa {
			albums[a.ID] = a
		}
	}

	if q.sharedAlbums {
		// if the user is an admin, get all albums regardless of permissions
		if user.Role == entity.RoleAdmin {
			sa, err := q.albumRepo.Get(ctx)
			if err != nil {
				return []entity.Album{}, 0, fmt.Errorf("%w all albums: %v", services.ErrGetAlbums, err)
			}

			for _, a := range sa {
				albums[a.ID] = a
			}
		} else if user.CanShare {
			sharedAlbums, err := q.albumRepo.GetByUserID(ctx, user.ID)
			if err != nil {
				return []entity.Album{}, 0, fmt.Errorf("%w shared albums: %v", services.ErrGetAlbums, err)
			}

			// get albums shared by the user's groups but filter out the ones owns by the user
			groupSharedAlbum, err := q.albumRepo.GetByGroups(ctx, groupsToList(user.Groups))
			if err != nil {
				return []entity.Album{}, 0, fmt.Errorf("%w shared albums by group: %v", services.ErrGetAlbums, err)
			}

			for i := 0; i < len(sharedAlbums)+len(groupSharedAlbum); i++ {
				found := false
				if i < len(sharedAlbums) {
					albums[sharedAlbums[i].ID] = sharedAlbums[i]
					found = true
				}

				if i < len(groupSharedAlbum) {
					albums[groupSharedAlbum[i].ID] = groupSharedAlbum[i]
					found = true
				}

				if !found {
					break
				}
			}
		}
	}

	// put all the albums into a list and return them
	albs := make([]entity.Album, 0, len(albums))
	for _, a := range albums {
		if q.filterEngine != nil {
			resolved, err := q.filterEngine.Resolve(a)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("album id", a.ID).Error("failed to resolve album")

				continue
			}

			if resolved {
				albs = append(albs, a)
			}
		} else {
			albs = append(albs, a)
		}
	}

	if q.sorter != nil {
		q.sorter.Sort(albs)
	}

	pages := q.paginate(albs)

	return pages, len(albs), nil
}

func (q *Query) First(ctx context.Context, id int32) (entity.Album, error) {
	album, err := q.albumRepo.GetByID(ctx, id)
	if err != nil {
		return entity.Album{}, fmt.Errorf("failed to get album '%d': %v", id, err)
	}

	medias, err := q.mediaService.ListBucket(ctx, album.Bucket)
	if err != nil {
		return entity.Album{}, fmt.Errorf("%w album id '%d': %v", services.ErrListBucket, id, err)
	}

	photos := make([]entity.Media, 0, len(medias))
	videos := make([]entity.Media, 0, len(medias))

	for _, m := range medias {
		switch m.MediaType {
		case entity.Photo:
			photos = append(photos, m)
		case entity.Video:
			videos = append(videos, m)
		}
	}

	album.Photos = photos
	album.Videos = videos

	return album, nil
}

func (q *Query) paginate(albums []entity.Album) []entity.Album {
	// pagination
	var page []entity.Album

	if q.offset > 0 && q.limit > 0 {
		if q.offset >= len(albums) {
			return []entity.Album{}
		}

		limit := q.limit
		if q.offset+limit >= len(albums) {
			limit = len(albums) - q.offset
		}

		page = append(page, albums[q.offset:q.offset+limit]...)

		return page
	}

	if q.offset > 0 && q.limit == 0 {
		if q.offset > len(albums) {
			return []entity.Album{}
		}

		page = append(page, albums[q.offset:]...)

		return page
	}

	if q.offset == 0 && q.limit > 0 {
		if q.limit > len(albums) {
			return albums
		}

		page = append(page, albums[:q.limit]...)

		return page
	}

	return albums
}

func groupsToList(groups []entity.Group) []string {
	l := make([]string, 0, len(groups))

	for _, g := range groups {
		l = append(l, g.Name)
	}

	return l
}
