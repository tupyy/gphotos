package album

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type FilterName int

const (
	// FilterAfterDate filter albums created after filter value.
	FilterAfterDate FilterName = iota
	// FilterBeforeDate filter albums created before filter value.
	FilterBeforeDate
	// FilterByOwnerID filter the album with owner in filter values.
	FilterByOwnerID
	// NotFilterByOwnerID filter out the albums whose owner is not in filter values.
	NotFilterByOwnerID
)

type Filter func(tx *gorm.DB) *gorm.DB

func GenerateFilterFuncs(filter FilterName, filterValues interface{}) (Filter, error) {
	switch filter {
	case FilterByOwnerID:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting list of strings", filter)
		}
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Where("album.owner_id IN ?", v)
		}, nil
	case NotFilterByOwnerID:
		v, ok := filterValues.([]string)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting list of strings", filter)
		}
		return func(tx *gorm.DB) *gorm.DB {
			return tx.Not("album.owner_id IN ?", v)
		}, nil
	case FilterBeforeDate:
		v, ok := filterValues.(time.Time)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting time.Time", filter)
		}
		return func(tx *gorm.DB) *gorm.DB {
			midnight := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
			return tx.Where("album.create_at < ?", midnight)
		}, nil
	case FilterAfterDate:
		v, ok := filterValues.(time.Time)
		if !ok {
			return nil, errors.Errorf("%v invalid value. expecting time.Time", filter)
		}
		return func(tx *gorm.DB) *gorm.DB {
			midnight := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
			return tx.Where("album.created_at > ?", midnight)
		}, nil
	}

	return nil, errors.Errorf("%v is invalid filter", filter)
}
