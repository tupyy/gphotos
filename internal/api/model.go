package api

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// simpleAlbum is a simplified struct of an album used in rendering templates.
// The ID of the simple album is encrypted.
type simpleAlbum struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Owner       string    `json:"owner"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}

func newSimpleAlbum(a entity.Album, owner entity.User) (simpleAlbum, error) {
	var ownerName string

	if len(owner.FirstName) == 0 && len(owner.LastName) == 0 {
		ownerName = owner.Username
	} else {
		ownerName = fmt.Sprintf("%s %s", owner.FirstName, owner.LastName)
	}

	// encrypt album id
	gen := encryption.NewGenerator(conf.GetEncryptionKey())
	encryptedID, err := gen.EncryptData(fmt.Sprintf("%d", a.ID))
	if err != nil {
		return simpleAlbum{}, err
	}

	logutil.GetDefaultLogger().WithFields(logrus.Fields{
		"id":           a.ID,
		"encrypted_id": encryptedID,
	}).Trace("encrypt album id")

	return simpleAlbum{
		ID:          encryptedID,
		Name:        a.Name,
		Date:        a.CreatedAt,
		Location:    a.Location,
		Description: a.Description,
		Owner:       ownerName,
	}, nil
}

// newSimpleAlbums returns a new list of simpleAlbums.
func newSimpleAlbums(albums []entity.Album, users []entity.User) []simpleAlbum {
	sAlbums := make([]simpleAlbum, 0, len(albums))

	// put users into a map
	usersMap := make(map[string]entity.User)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	for _, pa := range albums {
		if user, found := usersMap[pa.OwnerID]; found {
			sAlbum, err := newSimpleAlbum(pa, user)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("album", fmt.Sprintf("%+v", pa)).Error("cannot create simple album")

				continue
			}

			sAlbums = append(sAlbums, sAlbum)
		} else {
			logutil.GetDefaultLogger().WithField("album", fmt.Sprintf("%+v", pa.String())).Warn("owner don't exists anymore")
		}
	}

	return sAlbums
}
