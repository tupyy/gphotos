package filters

import (
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/utils"
)

const (
	// FilterByUsername returns true if album's name is in filterValues
	FilterByUsername Filter = iota
	// FilterByRole returns true if the user has the filter role
	FilterByRole
	// NotFilterByUsername negates the FilterByUsername
	NotFilterByUsername
	// FilterByCanShare returns true if users has can_share set to true.
	FilterByCanShare
	FilterByName
)

type UserFilter func(user entity.User) bool

func GenerateUserFilterFuncs(filter Filter, filterValues interface{}) (UserFilter, error) {
	switch filter {
	case FilterByName:
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
		return func(user entity.User) bool {
			return user.CanShare == true
		}, nil
	}

	return nil, errors.Errorf("%v is invalid filter", filter)
}
