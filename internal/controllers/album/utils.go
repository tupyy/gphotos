package album

import (
	"fmt"

	"github.com/tupyy/gophoto/internal/controllers/common"
	"github.com/tupyy/gophoto/internal/domain/entity"
	domainFilters "github.com/tupyy/gophoto/internal/domain/filters"
	"github.com/tupyy/gophoto/utils/logutil"
)

// generateFilters generates 3 filters: notUserNameFilter, FilterByRole and FilterByCanShare.
func generateFilters(currentUser entity.User) ([]domainFilters.UserFilter, error) {
	userFilters := make([]domainFilters.UserFilter, 0, 3)

	// get other users with can_share true except the current user
	usernameFilter, err := domainFilters.GenerateUserFilterFuncs(domainFilters.NotFilterByUsername, []string{currentUser.Username})
	if err != nil {
		return []domainFilters.UserFilter{}, err
	}

	userFilters = append(userFilters, usernameFilter)

	// only can share users
	canShareFilter, err := domainFilters.GenerateUserFilterFuncs(domainFilters.FilterByCanShare, []string{})
	if err != nil {
		return []domainFilters.UserFilter{}, err
	}

	userFilters = append(userFilters, canShareFilter)

	// remove admins
	notAdminFilter, err := domainFilters.GenerateUserFilterFuncs(domainFilters.FilterByRole, []entity.Role{entity.RoleUser, entity.RoleEditor})
	if err != nil {
		return []domainFilters.UserFilter{}, err
	}

	userFilters = append(userFilters, notAdminFilter)

	return userFilters, nil
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
