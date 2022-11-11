package users

import (
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repos/filters/user"
	userFilters "github.com/tupyy/gophoto/internal/repos/filters/user"
)

type Predicate func() userFilters.Filter

func Username(username string) Predicate {
	return func() user.Filter {
		usernameFilter, _ := userFilters.GenerateFilterFuncs(userFilters.FilterByUsername, []string{username})

		return usernameFilter
	}
}

func NotUsername(username string) Predicate {
	return func() user.Filter {
		usernameFilter, _ := userFilters.GenerateFilterFuncs(userFilters.NotFilterByUsername, []string{username})

		return usernameFilter
	}
}

func CanShare(canShare bool) Predicate {
	return func() user.Filter {
		canShareFilter, _ := userFilters.GenerateFilterFuncs(userFilters.FilterByCanShare, canShare)

		return canShareFilter
	}
}

func Roles(roles []entity.Role) Predicate {
	return func() user.Filter {
		roleFilter, _ := userFilters.GenerateFilterFuncs(userFilters.FilterByRole, roles)
		return roleFilter
	}
}
