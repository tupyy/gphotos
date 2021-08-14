package api

import (
	"github.com/tupyy/gophoto/internal/domain/entity"
)

// joinsAlbums joins two list of albums
func joinAlbums(albums1, albums2 []entity.Album) []entity.Album {
	joinedAlbums := make([]entity.Album, 0, len(albums1)+len(albums2))

	joinedAlbums = append(joinedAlbums, albums1...)
	joinedAlbums = append(joinedAlbums, albums2...)

	return joinedAlbums
}
