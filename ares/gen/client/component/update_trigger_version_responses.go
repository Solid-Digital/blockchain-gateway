// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// UpdateTriggerVersionReader is a Reader for the UpdateTriggerVersion structure.
type UpdateTriggerVersionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateTriggerVersionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateTriggerVersionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewUpdateTriggerVersionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateTriggerVersionOK creates a UpdateTriggerVersionOK with default headers values
func NewUpdateTriggerVersionOK() *UpdateTriggerVersionOK {
	return &UpdateTriggerVersionOK{}
}

/*UpdateTriggerVersionOK handles this case with default header values.

Status 200
*/
type UpdateTriggerVersionOK struct {
	Payload *dto.GetComponentVersionResponse
}

func (o *UpdateTriggerVersionOK) Error() string {
	return fmt.Sprintf("[PATCH /orgs/{orgName}/triggers/{name}/versions/{version}][%d] updateTriggerVersionOK  %+v", 200, o.Payload)
}

func (o *UpdateTriggerVersionOK) GetPayload() *dto.GetComponentVersionResponse {
	return o.Payload
}

func (o *UpdateTriggerVersionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(dto.GetComponentVersionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateTriggerVersionInternalServerError creates a UpdateTriggerVersionInternalServerError with default headers values
func NewUpdateTriggerVersionInternalServerError() *UpdateTriggerVersionInternalServerError {
	return &UpdateTriggerVersionInternalServerError{}
}

/*UpdateTriggerVersionInternalServerError handles this case with default header values.

Internal server error
*/
type UpdateTriggerVersionInternalServerError struct {
	Payload interface{}
}

func (o *UpdateTriggerVersionInternalServerError) Error() string {
	return fmt.Sprintf("[PATCH /orgs/{orgName}/triggers/{name}/versions/{version}][%d] updateTriggerVersionInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateTriggerVersionInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *UpdateTriggerVersionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
