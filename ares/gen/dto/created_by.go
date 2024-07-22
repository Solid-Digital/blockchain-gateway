// Code generated by go-swagger; DO NOT EDIT.

package dto

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// CreatedBy created by
// swagger:model CreatedBy
type CreatedBy struct {

	// id
	ID int64 `json:"id,omitempty"`

	// full name
	FullName string `json:"fullName,omitempty"`
}

// Validate validates this created by
func (m *CreatedBy) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreatedBy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreatedBy) UnmarshalBinary(b []byte) error {
	var res CreatedBy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
