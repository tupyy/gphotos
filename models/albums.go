// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Albums albums
//
// swagger:model Albums
type Albums struct {

	// albums
	Albums []*UserAlbum `json:"albums"`

	// number of the page
	Page int32 `json:"page,omitempty"`

	// number of album in a page
	Size int32 `json:"size,omitempty"`

	// total
	Total int32 `json:"total,omitempty"`
}

// Validate validates this albums
func (m *Albums) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlbums(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Albums) validateAlbums(formats strfmt.Registry) error {
	if swag.IsZero(m.Albums) { // not required
		return nil
	}

	for i := 0; i < len(m.Albums); i++ {
		if swag.IsZero(m.Albums[i]) { // not required
			continue
		}

		if m.Albums[i] != nil {
			if err := m.Albums[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("albums" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("albums" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this albums based on the context it is used
func (m *Albums) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAlbums(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Albums) contextValidateAlbums(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Albums); i++ {

		if m.Albums[i] != nil {
			if err := m.Albums[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("albums" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("albums" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Albums) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Albums) UnmarshalBinary(b []byte) error {
	var res Albums
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
