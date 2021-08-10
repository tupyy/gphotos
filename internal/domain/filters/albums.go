package filters

import (
	"github.com/pkg/errors"
	"github.com/tupyy/gophoto/internal/domain/entity"
	"github.com/tupyy/gophoto/internal/domain/utils"
)

type AlbumFilter func(album entity.Album) bool

func GenerateAlbumFilterFuncs(filter string, filterValues []string) (AlbumFilter, error) {
	switch filter {
	case "name":
		return func(album entity.Album) bool {
			return utils.StringMatchRegexSlice(album.Name, filterValues)
		}, nil
	}

	return nil, errors.Errorf("%s is invalid filter", filter)
}
