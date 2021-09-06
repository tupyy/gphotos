package album

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// serializedAlbum is a simplified struct of an album used in rendering templates.
// The ID of the simple album is encrypted.
type serializedAlbum struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

func serializeAlbum(a entity.Album, owner entity.User) (serializedAlbum, error) {
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
		return serializedAlbum{}, err
	}

	logutil.GetDefaultLogger().WithFields(logrus.Fields{
		"id":           a.ID,
		"encrypted_id": encryptedID,
	}).Trace("encrypt album id")

	return serializedAlbum{
		ID:          encryptedID,
		Name:        a.Name,
		Date:        a.CreatedAt.Format("2 January 2006"),
		Location:    a.Location,
		Description: a.Description,
		Owner:       ownerName,
	}, nil
}

// serializeAlbums returns a new list of simpleAlbums.
func serializeAlbums(albums []entity.Album, users []entity.User) []serializedAlbum {
	sAlbums := make([]serializedAlbum, 0, len(albums))

	// put users into a map
	usersMap := make(map[string]entity.User)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	for _, pa := range albums {
		if user, found := usersMap[pa.OwnerID]; found {
			sAlbum, err := serializeAlbum(pa, user)
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

// serializedUser is a simplified version of user to be used in templates.
// The username is encrypted.
type serializedUser struct {
	EncryptedID string
	Username    string
	Name        string
	Role        entity.Role
	CanShare    bool
}

func serializeUser(u entity.User) (serializedUser, error) {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	encryptedUsername, err := gen.EncryptData(u.Username)
	if err != nil {
		return serializedUser{}, err
	}

	encryptedID, err := gen.EncryptData(u.Username)
	if err != nil {
		return serializedUser{}, err
	}

	return serializedUser{
		EncryptedID: encryptedID,
		Username:    encryptedUsername,
		Name:        fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		Role:        u.Role,
		CanShare:    u.CanShare,
	}, nil
}
