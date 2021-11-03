package models

import (
	"database/sql/driver"

	"github.com/lib/pq"
)

// Role is one of (admin, editor, user)
type Role string

// PermissionID is one of
/*
   'album.read',
   'album.write',
   'album.edit',
   'album.delete'
*/
type PermissionID string

// PermissionIDs is an array of PermissionID
// We need to define explicitly this type to overload the method 'Scan' on it.
// Which is called by GORM.
type PermissionIDs []PermissionID

func (p *PermissionIDs) Scan(src interface{}) error {
	arr := &pq.StringArray{}

	// Use pq.StringArray.Scan()
	if err := arr.Scan(src); err != nil {
		return err
	}

	// Convert pq.StringArray to PermissionIDs
	res := make(PermissionIDs, len(*arr))
	for i, v := range *arr {
		res[i] = PermissionID(v)
	}

	*p = res
	return nil
}

func (p PermissionIDs) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}

	// Convert PermissionIDs to pq.StringArray
	arr := make(pq.StringArray, len(p))
	for i, v := range p {
		arr[i] = string(v)
	}

	// Use pq.StringArray.Value()
	return arr.Value()
}
