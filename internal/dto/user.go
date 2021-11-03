package dto

import (
	"fmt"

	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
	"github.com/tupyy/gophoto/utils/logutil"
)

// User is a simplified version of user to be used in templates.
// The username is encrypted.
type User struct {
	ID          string
	EncryptedID string
	Username    string
	Name        string
	Role        string
	CanShare    bool
}

func NewUserDTO(u entity.User) (User, error) {
	gen := encryption.NewGenerator(conf.GetEncryptionKey())

	encryptedUsername, err := gen.EncryptData(u.Username)
	if err != nil {
		return User{}, err
	}

	encryptedID, err := gen.EncryptData(u.ID)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:          u.ID,
		EncryptedID: encryptedID,
		Username:    encryptedUsername,
		Name:        fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		Role:        u.Role.String(),
		CanShare:    u.CanShare,
	}, nil
}

func NewUserDTOs(users []entity.User) []User {

	ret := make([]User, 0, len(users))
	for _, u := range users {
		userDTO, err := NewUserDTO(u)
		if err != nil {
			logutil.GetDefaultLogger().WithError(err).WithField("user", fmt.Sprintf("%+v", u)).Error("failed to serialize user")

			continue
		}

		ret = append(ret, userDTO)
	}

	return ret
}
