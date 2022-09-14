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

// Tags tags
//
// swagger:model Tags
type Tags struct {
	ObjectReference

	// number of the page
	Page int32 `json:"page,omitempty"`

	// number of album in a page
	Size int32 `json:"size,omitempty"`

	// tags
	Tags []*Tag `json:"tags"`

	// total
	Total int32 `json:"total,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *Tags) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 ObjectReference
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.ObjectReference = aO0

	// AO1
	var dataAO1 struct {
		Page int32 `json:"page,omitempty"`

		Size int32 `json:"size,omitempty"`

		Tags []*Tag `json:"tags"`

		Total int32 `json:"total,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Page = dataAO1.Page

	m.Size = dataAO1.Size

	m.Tags = dataAO1.Tags

	m.Total = dataAO1.Total

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m Tags) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.ObjectReference)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Page int32 `json:"page,omitempty"`

		Size int32 `json:"size,omitempty"`

		Tags []*Tag `json:"tags"`

		Total int32 `json:"total,omitempty"`
	}

	dataAO1.Page = m.Page

	dataAO1.Size = m.Size

	dataAO1.Tags = m.Tags

	dataAO1.Total = m.Total

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this tags
func (m *Tags) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with ObjectReference
	if err := m.ObjectReference.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Tags) validateTags(formats strfmt.Registry) error {

	if swag.IsZero(m.Tags) { // not required
		return nil
	}

	for i := 0; i < len(m.Tags); i++ {
		if swag.IsZero(m.Tags[i]) { // not required
			continue
		}

		if m.Tags[i] != nil {
			if err := m.Tags[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tags" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this tags based on the context it is used
func (m *Tags) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with ObjectReference
	if err := m.ObjectReference.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTags(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Tags) contextValidateTags(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Tags); i++ {

		if m.Tags[i] != nil {
			if err := m.Tags[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tags" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Tags) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Tags) UnmarshalBinary(b []byte) error {
	var res Tags
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
