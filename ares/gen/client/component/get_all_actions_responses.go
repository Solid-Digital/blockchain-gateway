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

// GetAllActionsReader is a Reader for the GetAllActions structure.
type GetAllActionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAllActionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAllActionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetAllActionsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAllActionsOK creates a GetAllActionsOK with default headers values
func NewGetAllActionsOK() *GetAllActionsOK {
	return &GetAllActionsOK{}
}

/*GetAllActionsOK handles this case with default header values.

Status 200
*/
type GetAllActionsOK struct {
	Payload []*dto.GetComponentResponse
}

func (o *GetAllActionsOK) Error() string {
	return fmt.Sprintf("[GET /orgs/{orgName}/actions][%d] getAllActionsOK  %+v", 200, o.Payload)
}

func (o *GetAllActionsOK) GetPayload() []*dto.GetComponentResponse {
	return o.Payload
}

func (o *GetAllActionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAllActionsInternalServerError creates a GetAllActionsInternalServerError with default headers values
func NewGetAllActionsInternalServerError() *GetAllActionsInternalServerError {
	return &GetAllActionsInternalServerError{}
}

/*GetAllActionsInternalServerError handles this case with default header values.

Internal server error
*/
type GetAllActionsInternalServerError struct {
	Payload interface{}
}

func (o *GetAllActionsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /orgs/{orgName}/actions][%d] getAllActionsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetAllActionsInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *GetAllActionsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
