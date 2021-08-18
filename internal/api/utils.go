package api

import (
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

// substract remove all albums2 from albums1
func substract(albums1, albums2 []entity.Album) []entity.Album {
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
