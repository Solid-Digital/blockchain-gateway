// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// LogoutReader is a Reader for the Logout structure.
type LogoutReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *LogoutReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewLogoutOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewLogoutInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewLogoutOK creates a LogoutOK with default headers values
func NewLogoutOK() *LogoutOK {
	return &LogoutOK{}
}

/*LogoutOK handles this case with default header values.

Status 200
*/
type LogoutOK struct {
}

func (o *LogoutOK) Error() string {
	return fmt.Sprintf("[GET /auth/logout][%d] logoutOK ", 200)
}

func (o *LogoutOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewLogoutInternalServerError creates a LogoutInternalServerError with default headers values
func NewLogoutInternalServerError() *LogoutInternalServerError {
	return &LogoutInternalServerError{}
}

/*LogoutInternalServerError handles this case with default header values.

Internal server error
*/
type LogoutInternalServerError struct {
	Payload interface{}
}

func (o *LogoutInternalServerError) Error() string {
	return fmt.Sprintf("[GET /auth/logout][%d] logoutInternalServerError  %+v", 500, o.Payload)
}

func (o *LogoutInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *LogoutInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
