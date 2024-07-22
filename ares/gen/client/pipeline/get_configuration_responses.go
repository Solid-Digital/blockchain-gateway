// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// GetConfigurationReader is a Reader for the GetConfiguration structure.
type GetConfigurationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetConfigurationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetConfigurationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetConfigurationInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetConfigurationOK creates a GetConfigurationOK with default headers values
func NewGetConfigurationOK() *GetConfigurationOK {
	return &GetConfigurationOK{}
}

/*GetConfigurationOK handles this case with default header values.

Status 200
*/
type GetConfigurationOK struct {
	Payload *dto.GetConfigurationResponse
}

func (o *GetConfigurationOK) Error() string {
	return fmt.Sprintf("[GET /orgs/{orgName}/pipelines/{pipelineName}/configurations/{revision}][%d] getConfigurationOK  %+v", 200, o.Payload)
}

func (o *GetConfigurationOK) GetPayload() *dto.GetConfigurationResponse {
	return o.Payload
}

func (o *GetConfigurationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(dto.GetConfigurationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetConfigurationInternalServerError creates a GetConfigurationInternalServerError with default headers values
func NewGetConfigurationInternalServerError() *GetConfigurationInternalServerError {
	return &GetConfigurationInternalServerError{}
}

/*GetConfigurationInternalServerError handles this case with default header values.

Internal server error
*/
type GetConfigurationInternalServerError struct {
	Payload interface{}
}

func (o *GetConfigurationInternalServerError) Error() string {
	return fmt.Sprintf("[GET /orgs/{orgName}/pipelines/{pipelineName}/configurations/{revision}][%d] getConfigurationInternalServerError  %+v", 500, o.Payload)
}

func (o *GetConfigurationInternalServerError) GetPayload() interface{} {
	return o.Payload
}

func (o *GetConfigurationInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
