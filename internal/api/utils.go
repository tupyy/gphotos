package api

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

func getUsers(ctx context.Context, k domain.KeycloakRepo) (map[string]entity.User, error) {
	users, err := k.GetUsers(ctx)
	if err != nil {
		return nil, err

	}

	mappedUsers := make(map[string]entity.User)
	for _, u := range users {
		mappedUsers[u.ID] = u
	}

	return mappedUsers, nil
}

// substractAlbums returns albums1 - albums2
func substractAlbums(albums1, albums2 []entity.Album) []entity.Album {
	album1Map := make(map[int32]entity.Album)
	album2Map := make(map[int32]entity.Album)

	var iLimit int
	if len(albums1) > len(albums2) {
		iLimit = len(albums1)
	} else {
		iLimit = len(albums2)
	}

	for i := 0; i < iLimit; i++ {
		if i < len(albums1) {
			album1Map[albums1[i].ID] = albums1[i]
		}

		if i < len(albums2) {
			album2Map[albums2[i].ID] = albums2[i]
		}
	}

	distinctAlbums := make([]entity.Album, 0, len(albums1))
	for k, v := range album1Map {
		if _, found := album2Map[k]; !found {
			distinctAlbums = append(distinctAlbums, v)
		}
	}

	return distinctAlbums
}

// addGroups add distinct albums from albums2 to album1
func addGroups(albums1, albums2 []entity.Album) []entity.Album {
	album1Map := make(map[int32]entity.Album)
	album2Map := make(map[int32]entity.Album)

	var iLimit int
	if len(albums1) > len(albums2) {
		iLimit = len(albums1)
	} else {
		iLimit = len(albums2)
	}

	for i := 0; i < iLimit; i++ {
		if i < len(albums1) {
			album1Map[albums1[i].ID] = albums1[i]
		}

		if i < len(albums2) {
			album2Map[albums2[i].ID] = albums2[i]
		}
	}

	albumSum := make([]entity.Album, 0, len(albums1)+len(albums2))
	for k, v := range album1Map {
		if _, found := album1Map[k]; found {
			albumSum = append(albumSum, v)
			delete(album2Map, k)
		}
	}

	for _, v := range album2Map {
		albumSum = append(albumSum, v)
	}

	return albumSum
}

func mapUsersToAlbums(albums []entity.Album, users map[string]entity.User) []tAlbum {
	tAlbums := make([]tAlbum, 0, len(albums))
	for _, pa := range albums {
		if user, found := users[pa.OwnerID]; found {
			talbum, err := mapTAbum(pa, user)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("album", fmt.Sprintf("%+v", pa)).Error("cannot map shared album")

				continue
			}

			tAlbums = append(tAlbums, talbum)
		} else {
			logutil.GetDefaultLogger().WithField("album", fmt.Sprintf("%+v", pa.String())).Warn("owner don't exists anymore")
		}
	}

	return tAlbums
}

func mapTAbum(a entity.Album, user entity.User) (tAlbum, error) {
	var owner string

	if len(user.FirstName) == 0 && len(user.LastName) == 0 {
		owner = user.Username
	} else {
		owner = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	// encrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", a.ID))
	if err != nil {
		return tAlbum{}, err
	}

	logutil.GetDefaultLogger().WithFields(logrus.Fields{
		"id":           a.ID,
		"encrypted_id": encryptedID,
	}).Trace("encrypt album id")

	return tAlbum{
		ID:          encryptedID,
		Name:        a.Name,
		Date:        a.CreatedAt,
		Location:    a.Location,
		Description: a.Description,
		Owner:       owner,
	}, nil
}

type tAlbum struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Owner       string    `json:"owner"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}
