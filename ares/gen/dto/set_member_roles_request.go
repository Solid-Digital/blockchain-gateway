// Code generated by go-swagger; DO NOT EDIT.

package dto

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// SetMemberRolesRequest set member roles request
// swagger:model SetMemberRolesRequest
type SetMemberRolesRequest struct {

	// Roles is map[role]bool
	Roles map[string]bool `json:"roles,omitempty"`
}

// Validate validates this set member roles request
func (m *SetMemberRolesRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SetMemberRolesRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SetMemberRolesRequest) UnmarshalBinary(b []byte) error {
	var res SetMemberRolesRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
