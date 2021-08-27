package filters

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Filter int

const (
	// FilterAfterDate returns true if album createdAt is filter value
	FilterAfterDate = iota
	// FilterBeforeDate returns true if album createdAt is before filter value
	FilterBeforeDate
	// FilterByOwnerID returns true if albums ownerid equal filter value
	FilterByOwnerID
	// NotFilterByOwnerID returns true if albums owner is not the filter value
	NotFilterByOwnerID
)

type AlbumFilter func(tx *gorm.DB) *gorm.DB

func GenerateAlbumFilterFuncs(filter Filter, filterValues interface{}) (AlbumFilter, error) {
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
