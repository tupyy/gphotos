package index

import (
	"github.com/tupyy/gophoto/internal/domain/entity"
	userFilters "github.com/tupyy/gophoto/internal/domain/filters/user"
)

// generateFilters generates 3 filters: notUserNameFilter, FilterByRole and FilterByCanShare.
func generateFilters(currentUser entity.User) (userFilters.Filters, error) {
	filters := make(map[string]userFilters.Filter)

	// get other users with can_share true except the current user
	usernameFilter, err := userFilters.GenerateFilterFuncs(userFilters.NotFilterByUsername, []string{currentUser.Username})
	if err != nil {
		return nil, err
	}

	filters[currentUser.ID] = usernameFilter

	// if admin do not filter users. get them all
	if currentUser.Role == entity.RoleAdmin {
		return filters, nil
	}

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
