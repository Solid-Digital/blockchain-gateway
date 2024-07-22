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

// CreateTriggerVersionReader is a Reader for the CreateTriggerVersion structure.
type CreateTriggerVersionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateTriggerVersionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateTriggerVersionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewCreateTriggerVersionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCreateTriggerVersionOK creates a CreateTriggerVersionOK with default headers values
func NewCreateTriggerVersionOK() *CreateTriggerVersionOK {
	return &CreateTriggerVersionOK{}
}

/*CreateTriggerVersionOK handles this case with default header values.

Status 200
*/
type CreateTriggerVersionOK struct {
	Payload *dto.GetComponentVersionResponse
}

func (o *CreateTriggerVersionOK) Error() string {
	return fmt.Sprintf("[POST /orgs/{orgName}/triggers/{name}/versions][%d] createTriggerVersionOK  %+v", 200, o.Payload)
}

func (o *CreateTriggerVersionOK) GetPayload() *dto.GetComponentVersionResponse {
	return o.Payload
}

func (o *CreateTriggerVersionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(dto.GetComponentVersionResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateTriggerVersionInternalServerError creates a CreateTriggerVersionInternalServerError with default headers values
func NewCreateTriggerVersionInternalServerError() *CreateTriggerVersionInternalServerError {
	return &CreateTriggerVersionInternalServerError{}
}

/*CreateTriggerVersionInternalServerError handles this case with default header values.

Internal server error
*/
type CreateTriggerVersionInternalServerError struct {
	Payload interface{}
}

func (o *CreateTriggerVersionInternalServerError) Error() string {
	return fmt.Sprintf("[POST /orgs/{orgName}/triggers/{name}/versions][%d] createTriggerVersionInternalServerError  %+v", 500, o.Payload)
}

func (o *CreateTriggerVersionInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *CreateTriggerVersionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
