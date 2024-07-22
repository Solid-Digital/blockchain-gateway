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

// CreateBaseReader is a Reader for the CreateBase structure.
type CreateBaseReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateBaseReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateBaseOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewCreateBaseInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCreateBaseOK creates a CreateBaseOK with default headers values
func NewCreateBaseOK() *CreateBaseOK {
	return &CreateBaseOK{}
}

/*CreateBaseOK handles this case with default header values.

Status 200
*/
type CreateBaseOK struct {
	Payload *dto.GetComponentResponse
}

func (o *CreateBaseOK) Error() string {
	return fmt.Sprintf("[POST /orgs/{orgName}/bases][%d] createBaseOK  %+v", 200, o.Payload)
}

func (o *CreateBaseOK) GetPayload() *dto.GetComponentResponse {
	return o.Payload
}

func (o *CreateBaseOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(dto.GetComponentResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateBaseInternalServerError creates a CreateBaseInternalServerError with default headers values
func NewCreateBaseInternalServerError() *CreateBaseInternalServerError {
	return &CreateBaseInternalServerError{}
}

/*CreateBaseInternalServerError handles this case with default header values.

Internal server error
*/
type CreateBaseInternalServerError struct {
	Payload interface{}
}

func (o *CreateBaseInternalServerError) Error() string {
	return fmt.Sprintf("[POST /orgs/{orgName}/bases][%d] createBaseInternalServerError  %+v", 500, o.Payload)
}

func (o *CreateBaseInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *CreateBaseInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
