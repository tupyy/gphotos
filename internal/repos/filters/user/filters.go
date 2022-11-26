package user

import (
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/entity"
	"github.com/tupyy/gophoto/internal/repos/utils"
)

type FilterName int

const (
	// FilterByUsername returns true if album's name is in filterValues
	FilterByUsername FilterName = iota
	// FilterByRole returns true if the user has the filter role
	FilterByRole
	// NotFilterByUsername negates the FilterByUsername
	NotFilterByUsername
	// FilterByCanShare returns true if users has can_share set to true.
	FilterByCanShare
)

type Filter func(user entity.User) bool

// Filters defines a collection of filters. The key is the id of the filter which depends on the value of the filter.
// The id is used to compute cache keys in order to cache query results based on filters' values.
type Filters []Filter

func GenerateFilterFuncs(filter FilterName, filterValues interface{}) (Filter, error) {
	switch filter {
	case FilterByUsername:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting []string", filter)
		}

		return func(user entity.User) bool {
			return utils.StringMatchRegexSlice(user.Username, v)
		}, nil
	case FilterByRole:
		v, ok := filterValues.([]entity.Role)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting []entity.Role", filter)
		}

		return func(user entity.User) bool {
			for _, r := range v {
				if user.Role == r {
					return true
				}
			}

			return false
		}, nil
	case NotFilterByUsername:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting []string", filter)
		}

		return func(user entity.User) bool {
			return !utils.StringMatchRegexSlice(user.Username, v)
		}, nil
	case FilterByCanShare:
		v, ok := filterValues.(bool)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting bool", filter)
		}
		return func(user entity.User) bool {
			return user.CanShare == v
		}, nil
	}

	return nil, errors.Errorf("%v is invalid filter", filter)
}
