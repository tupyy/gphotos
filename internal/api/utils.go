package api

import (
	"errors"

	"github.com/tupyy/gophoto/internal/domain/entity"
)

// joinsAlbums joins two list of albums
func merge(albums1, albums2 []entity.Album) []entity.Album {
	joinedAlbums := make([]entity.Album, 0, len(albums1)+len(albums2))

	joinedAlbums = append(joinedAlbums, albums1...)
	joinedAlbums = append(joinedAlbums, albums2...)

	return joinedAlbums
}

// join joins two list of albums but removes the duplicates.
func join(albums1, albums2 []entity.Album) []entity.Album {
	m1 := make(map[int32]entity.Album)
	m2 := make(map[int32]entity.Album)

	iLimit := len(albums1)
	if iLimit < len(albums2) {
		iLimit = len(albums2)
	}

	for i := 0; i < iLimit; i++ {
		if i < len(albums1) {
			m1[albums1[i].ID] = albums1[i]
		}

		if i < len(albums2) {
			m2[albums2[i].ID] = albums2[i]
		}
	}

	joinedAlbums := make([]entity.Album, 0, len(albums1)+len(albums2))
	for id, album := range m1 {
		joinedAlbums = append(joinedAlbums, album)

		if _, found := m2[id]; found {
			delete(m2, id)
		}
	}

	// get the remains from m2
	for _, a := range m2 {
		joinedAlbums = append(joinedAlbums, a)
	}

	return joinedAlbums
}

func substract(a, b interface{}) (interface{}, error) {
	a1, isAlbum1 := a.([]entity.Album)
	a2, isAlbum2 := b.([]entity.Album)

	if isAlbum1 && isAlbum2 {
		return substractAlbums(a1, a2), nil
	}

	u1, isUser1 := a.([]entity.User)
	u2, isUser2 := b.([]entity.User)

	if isUser1 && isUser2 {
		return substractUsers(u1, u2), nil
	}

	return nil, errors.New("wrong arguments type")
}

// substract remove all albums2 from albums1
func substractAlbums(albums1, albums2 []entity.Album) []entity.Album {
	m1 := make(map[int32]entity.Album)
	m2 := make(map[int32]entity.Album)

	iLimit := len(albums1)
	if iLimit < len(albums2) {
		iLimit = len(albums2)
	}

	for i := 0; i < iLimit; i++ {
		if i < len(albums1) {
			m1[albums1[i].ID] = albums1[i]
		}

		if i < len(albums2) {
			m2[albums2[i].ID] = albums2[i]
		}
	}

	res := make([]entity.Album, 0, len(albums1)+len(albums2))
	for id, album := range m1 {
		if _, found := m2[id]; !found {
			res = append(res, album)
		}
	}

	return res
}

// substract remove all users2 from users1
func substractUsers(users1, users2 []entity.User) []entity.User {
	m1 := make(map[string]entity.User)
	m2 := make(map[string]entity.User)

	iLimit := len(users1)
	if iLimit < len(users2) {
		iLimit = len(users2)
	}

	for i := 0; i < iLimit; i++ {
		if i < len(users1) {
			m1[users1[i].ID] = users1[i]
		}

		if i < len(users2) {
			m2[users2[i].ID] = users2[i]
		}
	}

	res := make([]entity.User, 0, len(users1)+len(users2))
	for id, user := range m1 {
		if _, found := m2[id]; !found {
			res = append(res, user)
		}
	}

	return res
}
