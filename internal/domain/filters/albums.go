package filters

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/utils"
)

type Filter int

const (
	// FilterByName returns true if album's name is in filter values
	FilterByName Filter = iota
	// FilterNotInList returns true if the album is not in filter values
	FilterNotInList
	// FilterAfterDate returns true if album createdAt is filter value
	FilterAfterDate
	// FilterBeforeDate returns true if album createdAt is before filter value
	FilterBeforeDate
	// FilterByOwnerID returns true if albums ownerid equal filter value
	FilterByOwnerID
	// NotFilterByOwnerID returns true if albums owner is not the filter value
	NotFilterByOwnerID
)

type AlbumFilter func(album entity.Album) bool

func GenerateAlbumFilterFuncs(filter Filter, filterValues interface{}) (AlbumFilter, error) {
	switch filter {
	case FilterByName:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting []string", filter)
		}

		return func(album entity.Album) bool {
			return utils.StringMatchRegexSlice(album.Name, v)
		}, nil
	case FilterNotInList:
		v, ok := filterValues.([]entity.Album)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting list of albums.", filter)
		}

		return func(album entity.Album) bool {
			for _, a := range v {
				if a.ID == album.ID {
					return false
				}
			}

			return true
		}, nil
	case FilterByOwnerID:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting []string", filter)
		}
		return func(album entity.Album) bool {
			return utils.StringInSlice(album.OwnerID, v)
		}, nil
	case NotFilterByOwnerID:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid values. expecting []string", filter)
		}
		return func(album entity.Album) bool {
			return !utils.StringInSlice(album.OwnerID, v)
		}, nil
	case FilterBeforeDate:
		v, ok := filterValues.(time.Time)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting time.Time", filter)
		}
		return func(album entity.Album) bool {
			return album.CreatedAt.Before(v)
		}, nil
	case FilterAfterDate:
		v, ok := filterValues.(time.Time)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting time.Time", filter)
		}
		return func(album entity.Album) bool {
			return album.CreatedAt.After(v)
		}, nil
	}

	return nil, errors.Errorf("%v is invalid filter", filter)
}
