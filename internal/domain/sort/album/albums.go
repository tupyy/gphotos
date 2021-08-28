package album

import (
	"sort"

	"github.com/tupyy/gophoto/internal/domain/entity"
)

type SortOrder int

const (
	NormalOrder SortOrder = iota
	ReverseOrder
)

type SortName int

const (
	SortByName SortName = iota
	SortByID
	SortByOwner
	SortByDate
	SortByLocation
)

type NoSorter struct{}

func (n NoSorter) Sort(albums []entity.Album) {}

type Sorter interface {
	Sort(albums []entity.Album)
}

type albumSorter struct {
	album    []entity.Album
	lessFunc func(a1, a2 entity.Album) bool
}

// NewAlbumSorterById returns a sorter by IDs.
func NewSorter(name SortName, order SortOrder) *albumSorter {
	var lessFunc func(a1, a2 entity.Album) bool

	switch name {
	case SortByID:
		lessFunc = func(a1, a2 entity.Album) bool {
			if order == ReverseOrder {
				return a1.ID > a2.ID
			}

			return a1.ID < a2.ID
		}
	case SortByName:
		lessFunc = func(a1, a2 entity.Album) bool {
			if order == ReverseOrder {
				return a1.Name > a2.Name
			}

			return a1.Name < a2.Name
		}
	case SortByLocation:
		lessFunc = func(a1, a2 entity.Album) bool {
			if order == ReverseOrder {
				return a1.Location > a2.Location
			}

			return a1.Location < a2.Location
		}
	case SortByDate:
		lessFunc = func(a1, a2 entity.Album) bool {
			if order == ReverseOrder {
				return a1.CreatedAt.After(a2.CreatedAt)
			}
			return a1.CreatedAt.Before(a2.CreatedAt)
		}
	default:
		// dont sort here
		lessFunc = func(a1, a2 entity.Album) bool {
			return true
		}
	}

	return &albumSorter{lessFunc: lessFunc}
}

func (as *albumSorter) Sort(albums []entity.Album) {
	as.album = albums
	sort.Sort(as)
}

func (as *albumSorter) Len() int {
	return len(as.album)
}

func (as *albumSorter) Swap(i, j int) {
	as.album[i], as.album[j] = as.album[j], as.album[i]
}

func (as *albumSorter) Less(i, j int) bool {
	return as.lessFunc(as.album[i], as.album[j])
}
