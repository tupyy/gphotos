package index

import (
	"fmt"

	"github.com/tupyy/gophoto/internal/domain/entity"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
	"github.com/tupyy/gophoto/internal/handlers/common"
	"github.com/tupyy/gophoto/utils/logutil"
)

func serialize(users []entity.User) []common.SerializedUser {
	serializedUsers := make([]common.SerializedUser, 0, len(users))

	for _, u := range users {
		s, err := common.NewSerializedUser(u)
		if err != nil {
			logutil.GetDefaultLogger().WithError(err).WithField("user", fmt.Sprintf("%+v", u)).Error("serialize user")

			continue
		}

		serializedUsers = append(serializedUsers, s)
	}

	return serializedUsers
}

// generateFilters generates 3 filters: notUserNameFilter, FilterByRole and FilterByCanShare.
func generateFilters(currentUser entity.User) ([]userFilters.Filter, error) {
	filters := make([]userFilters.Filter, 0, 3)

	// get other users with can_share true except the current user
	usernameFilter, err := userFilters.GenerateFilterFuncs(userFilters.NotFilterByUsername, []string{currentUser.Username})
	if err != nil {
		return []userFilters.Filter{}, err
	}

	filters = append(filters, usernameFilter)

	// only can share users
	canShareFilter, err := userFilters.GenerateFilterFuncs(userFilters.FilterByCanShare, []string{})
	if err != nil {
		return []userFilters.Filter{}, err
	}

	filters = append(filters, canShareFilter)

	// remove admins
	notAdminFilter, err := userFilters.GenerateFilterFuncs(userFilters.FilterByRole, []entity.Role{entity.RoleUser, entity.RoleEditor})
	if err != nil {
		return []userFilters.Filter{}, err
	}

	filters = append(filters, notAdminFilter)

	return filters, nil
}
