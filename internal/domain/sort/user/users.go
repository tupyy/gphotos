package user

import (
	"fmt"
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
	SortByUsername SortName = iota
	SortByName
)

type Sorter interface {
	Sort(albums []entity.User)
}

type userSorter struct {
	album    []entity.User
	lessFunc func(u1, u2 entity.User) bool
}

func NewSorter(name SortName, order SortOrder) *userSorter {
	var lessFunc func(u1, u2 entity.User) bool

	switch name {
	case SortByUsername:
		lessFunc = func(u1, u2 entity.User) bool {
			if order == ReverseOrder {
				return u1.Username > u2.Username
			}

			return u1.Username < u2.Username
		}
	case SortByName:
		lessFunc = func(u1, u2 entity.User) bool {
			name1 := fmt.Sprintf("%s %s", u1.FirstName, u1.LastName)
			name2 := fmt.Sprintf("%s %s", u2.FirstName, u2.LastName)

			if order == ReverseOrder {
				return name1 > name2
			}

			return name2 < name1
		}
	default:
		// dont sort here
		lessFunc = func(u1, u2 entity.User) bool {
			return true
		}
	}

	return &userSorter{lessFunc: lessFunc}
}

func (as *userSorter) Sort(albums []entity.User) {
	as.album = albums
	sort.Sort(as)
}

func (as *userSorter) Len() int {
	return len(as.album)
}

func (as *userSorter) Swap(i, j int) {
	as.album[i], as.album[j] = as.album[j], as.album[i]
}

func (as *userSorter) Less(i, j int) bool {
	return as.lessFunc(as.album[i], as.album[j])
}
