package sort

import (
	"fmt"
	"sort"

	"github.com/tupyy/gophoto/internal/domain/entity"
)

type UserSorter interface {
	Sort(albums []entity.User)
}

type UserLessFunc func(u1, u2 entity.User) bool

type userSorter struct {
	album    []entity.User
	lessFunc UserLessFunc
}

// NewUserSorterById returns a sorter by IDs.
func NewUserSorterByUsername(order SortOrder) *userSorter {
	lessFunc := func(u1, u2 entity.User) bool {
		if order == ReverseOrder {
			return u1.Username > u2.Username
		}

		return u1.Username < u2.Username
	}

	return NewUserSorter(lessFunc)
}

// NewUserSorterByName returns a sorter by name.
func NewUserSorterByName(order SortOrder) *userSorter {
	nameLessFunc := func(u1, u2 entity.User) bool {
		name1 := fmt.Sprintf("%s %s", u1.FirstName, u1.LastName)
		name2 := fmt.Sprintf("%s %s", u2.FirstName, u2.LastName)

		if order == ReverseOrder {
			return name1 > name2
		}

		return name2 < name1
	}

	return NewUserSorter(nameLessFunc)
}

// NewUserSorter returns a custom sorter. The user must provide a lessFunc.
func NewUserSorter(lessFunc UserLessFunc) *userSorter {
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
