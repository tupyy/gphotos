package keycloak

import (
	keycloak "github.com/Nerzal/gocloak/v8"
	"github.com/tupyy/gophoto/internal/entity"
)

func mapper(u keycloak.User) entity.User {
	user := entity.User{
		Username: *u.Username,
		ID:       *u.ID,
	}

	if u.FirstName != nil {
		user.FirstName = *u.FirstName
	}

	if u.LastName != nil {
		user.LastName = *u.LastName
	}

	if u.Attributes != nil {
		m := *u.Attributes
		if attrs, found := m["can_share"]; found {
			for _, attr := range attrs {
				switch attr {
				case "true":
					user.CanShare = true
				default:
					user.CanShare = false
				}
			}
		}
	}

	return user
}
