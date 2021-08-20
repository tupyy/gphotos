package api

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// simpleAlbum is a simplified struct of an album used in rendering templates.
// The ID of the simple album is encrypted.
type simpleAlbum struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Day         string `json:"day"`
	Month       string `json:"month"`
	Year        string `json:"year"`
	Description string `json:"description"`
	Location    string `json:"location"`
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
		Day:         fmt.Sprintf("%d", a.CreatedAt.Day()),
		Month:       a.CreatedAt.Month().String()[:3],
		Year:        fmt.Sprintf("%d", a.CreatedAt.Year()),
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

// SerializedUser is a simplified version of user to be used in templates.
// The username is encrypted.
type SerializedUser struct {
	EncryptedID string
	Username    string
	Name        string
	Role        entity.Role
	CanShare    bool
}

func NewSerializedUser(u entity.User) (SerializedUser, error) {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	encryptedUsername, err := gen.EncryptData(u.Username)
	if err != nil {
		return SerializedUser{}, err
	}

	encryptedID, err := gen.EncryptData(u.Username)
	if err != nil {
		return SerializedUser{}, err
	}

	return SerializedUser{
		EncryptedID: encryptedID,
		Username:    encryptedUsername,
		Name:        fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		Role:        u.Role,
		CanShare:    u.CanShare,
	}, nil
}
