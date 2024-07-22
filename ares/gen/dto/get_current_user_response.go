// Code generated by go-swagger; DO NOT EDIT.

package dto

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// GetCurrentUserResponse get current user response
// swagger:model GetCurrentUserResponse
type GetCurrentUserResponse struct {

	// id
	ID int64 `json:"id,omitempty"`

	// full name
	FullName string `json:"fullName,omitempty"`

	// email
	Email string `json:"email,omitempty"`

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`

	// default organization
	DefaultOrganization string `json:"defaultOrganization,omitempty"`

	// Roles is map[org][role]bool
	Roles map[string]map[string]bool `json:"roles,omitempty"`
}

// Validate validates this get current user response
func (m *GetCurrentUserResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GetCurrentUserResponse) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("createdAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GetCurrentUserResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetCurrentUserResponse) UnmarshalBinary(b []byte) error {
	var res GetCurrentUserResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
