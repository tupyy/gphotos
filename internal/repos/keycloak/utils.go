package keycloak

import (
	"github.com/tupyy/gophoto/internal/entity"
	userFilters "github.com/tupyy/gophoto/internal/repos/filters/user"
)

func ptrBool(b bool) *bool {
	return &b
}

func filter(filters userFilters.Filters, users []entity.User) []entity.User {
	filteredUsers := make([]entity.User, 0, len(users))
	for _, u := range users {
		pass := true
		for _, filter := range filters {
			if !filter(u) {
				pass = false
				break
			}
		}

		if pass {
			filteredUsers = append(filteredUsers, u)
		}
	}

	return filteredUsers
}
