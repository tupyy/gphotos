package sort

import (
	"sort"

	"github.com/tupyy/gophoto/internal/domain/entity"
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
func NewAlbumSorterById(albums []entity.Album, reverse bool) *albumSorter {
	lessFunc := func(a1, a2 entity.Album) bool {
		if reverse {
			return a1.ID > a2.ID
		}

		return a1.ID < a2.ID
	}

	return NewAlbumSorter(albums, lessFunc)
}

// NewAlbumSorterByName returns a sorter by name.
func NewAlbumSorterByName(albums []entity.Album, reverse bool) *albumSorter {
	nameLessFunc := func(a1, a2 entity.Album) bool {
		if reverse {
			return a1.Name > a2.Name
		}

		return a1.Name < a2.Name
	}

	return NewAlbumSorter(albums, nameLessFunc)
}

func NewAlbumSorterByDate(albums []entity.Album, reverse bool) *albumSorter {
	dateLessFunc := func(a1, a2 entity.Album) bool {
		if reverse {
			return a1.CreatedAt.After(a2.CreatedAt)
		}

		return a1.CreatedAt.Before(a2.CreatedAt)
	}

	return NewAlbumSorter(albums, dateLessFunc)
}

// NewAlbumSorter returns a custom sorter. The user must provide a lessFunc.
func NewAlbumSorter(albums []entity.Album, lessFunc AlbumLessFunc) *albumSorter {
	return &albumSorter{albums, lessFunc}
}

func (as *albumSorter) Sort(albums []entity.Album) {
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
