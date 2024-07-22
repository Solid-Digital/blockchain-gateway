// Code generated by go-swagger; DO NOT EDIT.

package dto

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// LogLine log line
// swagger:model LogLine
type LogLine struct {

	// caller
	Caller string `json:"caller,omitempty"`

	// function
	Function string `json:"function,omitempty"`

	// instance ID
	InstanceID string `json:"instanceID,omitempty"`

	// level
	Level string `json:"level,omitempty"`

	// text
	Text string `json:"text,omitempty"`

	// time
	Time string `json:"time,omitempty"`

	// timestamp
	Timestamp int64 `json:"timestamp,omitempty"`
}

// Validate validates this log line
func (m *LogLine) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LogLine) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LogLine) UnmarshalBinary(b []byte) error {
	var res LogLine
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
