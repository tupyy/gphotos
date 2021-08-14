package filters

import (
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/utils"
)

type Filter int

const (
	// FilterByName returns true if album's name is in filterValues
	FilterByName Filter = iota
	// FilterNotInList returns true if the album is not in filterValues
	FilterNotInList
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
	}

	return nil, errors.Errorf("%v is invalid filter", filter)
}
