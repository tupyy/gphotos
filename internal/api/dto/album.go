package dto

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// Album is a simplified struct of an album used in rendering templates.
// The ID of the simple album is encrypted.
type Album struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Owner       string  `json:"owner"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Thumbnail   string  `json:"thumbnail"`
	Photos      []Media `json:"photos"`
	Videos      []Media `json:"videos"`
}

type Media struct {
	Filename  string
	Thumbnail string
}

func NewAlbumDTO(a entity.Album, owner entity.User) (Album, error) {
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
		return Album{}, err
	}

	// encrypt media filenames
	encryptedPhotos := make([]Media, 0, len(a.Photos))
	encryptedVideos := make([]Media, 0, len(a.Videos))
	for i := 0; i < len(a.Photos)+len(a.Videos); i++ {
		more := false
		if i < len(a.Photos) {
			more = true
			filename, thumbnail := a.Photos[i].Filename, a.Photos[i].Thumbnail

			encryptedFilename, err := gen.EncryptData(filename)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("photo filename", filename).Error("failed to encrypt data")

				continue
			}

			encryptedThumbnail, err := gen.EncryptData(thumbnail)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("photo thumbnail", thumbnail).Error("failed to encrypt data")

				continue
			}

			encryptedPhotos = append(encryptedPhotos, Media{encryptedFilename, encryptedThumbnail})
		}

		if i < len(a.Videos) {
			more = true
			filename, thumbnail := a.Videos[i].Filename, a.Videos[i].Thumbnail

			encryptedFilename, err := gen.EncryptData(filename)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("video filename", filename).Error("failed to encrypt data")

				continue
			}

			encryptedThumbnail, err := gen.EncryptData(thumbnail)
			if err != nil {
				logutil.GetDefaultLogger().WithError(err).WithField("video thumbnail", thumbnail).Error("failed to encrypt data")

				continue
			}

			encryptedVideos = append(encryptedVideos, Media{encryptedFilename, encryptedThumbnail})
		}

		if !more {
			break
		}
	}

	logutil.GetDefaultLogger().WithFields(logrus.Fields{
		"id":           a.ID,
		"encrypted_id": encryptedID,
	}).Trace("encrypt album id")

	thumbnail := "/static/img/image_not_available.png"
	if len(a.Thumbnail) > 0 {
		thumbnail = fmt.Sprintf("/api/albums/%s/album/%s/media", encryptedID, a.Thumbnail)
	}

	return Album{
		ID:          encryptedID,
		Name:        a.Name,
		Date:        a.CreatedAt.Format("2 January 2006"),
		Location:    a.Location,
		Description: a.Description,
		Owner:       ownerName,
		Photos:      encryptedPhotos,
		Videos:      encryptedVideos,
		Thumbnail:   thumbnail,
	}, nil
}

// NewAlbumDTOs returns a new list of simpleAlbums.
func NewAlbumDTOs(albums []entity.Album, users []entity.User) []Album {
	sAlbums := make([]Album, 0, len(albums))

	// put users into a map
	usersMap := make(map[string]entity.User)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	for _, pa := range albums {
		if user, found := usersMap[pa.OwnerID]; found {
			sAlbum, err := NewAlbumDTO(pa, user)
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
