package album

import (
	"context"
	"fmt"

	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/filters/album"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/services"
	"github.com/tupyy/gophoto/internal/services/media"
)

type Query struct {
	limit  *int
	offset *int
	// get personal albums.
	personalAlbums bool
	// get shared albums.
	sharedAlbums bool
	// list of predicates for the query
	predicates []Predicate
	// album repo
	albumRepo domain.Album
	// media service
	mediaService *media.Service
	//album sorter
	sorter *albumSorter
}

func (s *Service) Query() *Query {
	return &Query{
		predicates:   []Predicate{},
		albumRepo:    s.albumRepo,
		mediaService: s.mediaService,
	}
}

func (q *Query) Where(p Predicate) *Query {
	q.predicates = append(q.predicates, p)

	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = &limit

	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = &offset

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

func (q *Query) All(ctx context.Context, user entity.User) ([]entity.Album, error) {
	// generate filters from predicates
	filters := make([]album.Filter, 0, len(q.predicates))
	for _, p := range q.predicates {
		filters = append(filters, p())
	}

	albums := make(map[int32]entity.Album)

	if q.personalAlbums {
		// fetch personal albums
		pa, err := q.albumRepo.GetByOwnerID(ctx, user.ID, filters)
		if err != nil {
			return []entity.Album{}, fmt.Errorf("%w personal album: %v", services.ErrGetAlbums, err)
		}

		for _, a := range pa {
			albums[a.ID] = a
		}
	}

	if q.sharedAlbums {
		// if the user is an admin, get all albums regardless of permissions
		if user.Role == entity.RoleAdmin {
			sa, err := q.albumRepo.Get(ctx, filters)
			if err != nil {
				return []entity.Album{}, fmt.Errorf("%w all albums: %v", services.ErrGetAlbums, err)
			}

			for _, a := range sa {
				albums[a.ID] = a
			}
		} else if user.CanShare {
			sharedAlbums, err := q.albumRepo.GetByUserID(ctx, user.ID, filters)
			if err != nil {
				return []entity.Album{}, fmt.Errorf("%w shared albums: %v", services.ErrGetAlbums, err)
			}

			// get albums shared by the user's groups but filter out the ones owns by the user
			groupSharedAlbum, err := q.albumRepo.GetByGroups(ctx, groupsToList(user.Groups), filters)
			if err != nil {
				return []entity.Album{}, fmt.Errorf("%w shared albums by group: %v", services.ErrGetAlbums, err)
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
	ret := make([]entity.Album, 0, len(albums))
	for _, v := range albums {
		ret = append(ret, v)
	}

	if q.sorter != nil {
		q.sorter.Sort(ret)
	}

	return ret, nil
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

func groupsToList(groups []entity.Group) []string {
	l := make([]string, 0, len(groups))

	for _, g := range groups {
		l = append(l, g.Name)
	}

	return l
}
