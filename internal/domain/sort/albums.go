package sort

import (
	"sort"

	"github.com/tupyy/gophoto/internal/domain/entity"
)

type SortOrder int

const (
	NormalOrder SortOrder = iota
	ReverseOrder
)

type AlbumSorter interface {
	Sort(albums []entity.Album)
}

type AlbumLessFunc func(a1, a2 entity.Album) bool

type albumSorter struct {
	album    []entity.Album
	lessFunc AlbumLessFunc
}

// NewAlbumSorterById returns a sorter by IDs.
func NewAlbumSorterById(order SortOrder) *albumSorter {
	lessFunc := func(a1, a2 entity.Album) bool {
		if order == ReverseOrder {
			return a1.ID > a2.ID
		}

		return a1.ID < a2.ID
	}

	return NewAlbumSorter(lessFunc)
}

// NewAlbumSorterByName returns a sorter by name.
func NewAlbumSorterByName(order SortOrder) *albumSorter {
	nameLessFunc := func(a1, a2 entity.Album) bool {
		if order == ReverseOrder {
			return a1.Name > a2.Name
		}

		return a1.Name < a2.Name
	}

	return NewAlbumSorter(nameLessFunc)
}

func NewAlbumSorterByDate(order SortOrder) *albumSorter {
	dateLessFunc := func(a1, a2 entity.Album) bool {
		if order == ReverseOrder {
			return a1.CreatedAt.After(a2.CreatedAt)
		}

		return a1.CreatedAt.Before(a2.CreatedAt)
	}

	return NewAlbumSorter(dateLessFunc)
}

// NewAlbumSorter returns a custom sorter. The user must provide a lessFunc.
func NewAlbumSorter(lessFunc AlbumLessFunc) *albumSorter {
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
