package album

import (
	"fmt"

	"github.com/tupyy/gophoto/internal/domain/entity"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
	"github.com/tupyy/gophoto/internal/handlers/common"
	"github.com/tupyy/gophoto/utils/logutil"
)

// generateFilters generates 3 filters: notUserNameFilter, FilterByRole and FilterByCanShare.
func generateFilters(currentUser entity.User) (userFilters.Filters, error) {
	filters := make(userFilters.Filters)

	// get other users with can_share true except the current user
	usernameFilter, err := userFilters.GenerateFilterFuncs(userFilters.NotFilterByUsername, []string{currentUser.Username})
	if err != nil {
		return nil, err
	}

	filters[currentUser.ID] = usernameFilter

	// only can share users
	canShareFilter, err := userFilters.GenerateFilterFuncs(userFilters.FilterByCanShare, []string{})
	if err != nil {
		return nil, err
	}

	filters["canshare"] = canShareFilter

	// remove admins
	notAdminFilter, err := userFilters.GenerateFilterFuncs(userFilters.FilterByRole, []entity.Role{entity.RoleUser, entity.RoleEditor})
	if err != nil {
		return nil, err
	}

	filters["admin"] = notAdminFilter
	return filters, nil
}

// serialize serialized a list of users
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
