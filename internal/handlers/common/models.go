package common

import (
	"fmt"

	"github.com/tupyy/gophoto/internal/conf"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/utils/encryption"
)

// SerializedUser is a simplified version of user to be used in templates.
// The username is encrypted.
type SerializedUser struct {
	ID          string
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

	encryptedID, err := gen.EncryptData(u.ID)
	if err != nil {
		return SerializedUser{}, err
	}

	return SerializedUser{
		ID:          u.ID,
		EncryptedID: encryptedID,
		Username:    encryptedUsername,
		Name:        fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		Role:        u.Role,
		CanShare:    u.CanShare,
	}, nil
}
