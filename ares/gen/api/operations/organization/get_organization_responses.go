// Code generated by go-swagger; DO NOT EDIT.

package organization

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// GetOrganizationOKCode is the HTTP code returned for type GetOrganizationOK
const GetOrganizationOKCode int = 200

/*GetOrganizationOK Status 200

swagger:response getOrganizationOK
*/
type GetOrganizationOK struct {

	/*
	  In: Body
	*/
	Payload *dto.GetOrganizationResponse `json:"body,omitempty"`
}

// NewGetOrganizationOK creates GetOrganizationOK with default headers values
func NewGetOrganizationOK() *GetOrganizationOK {

	return &GetOrganizationOK{}
}

// WithPayload adds the payload to the get organization o k response
func (o *GetOrganizationOK) WithPayload(payload *dto.GetOrganizationResponse) *GetOrganizationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get organization o k response
func (o *GetOrganizationOK) SetPayload(payload *dto.GetOrganizationResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOrganizationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetOrganizationInternalServerErrorCode is the HTTP code returned for type GetOrganizationInternalServerError
const GetOrganizationInternalServerErrorCode int = 500

/*GetOrganizationInternalServerError Internal server error

swagger:response getOrganizationInternalServerError
*/
type GetOrganizationInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetOrganizationInternalServerError creates GetOrganizationInternalServerError with default headers values
func NewGetOrganizationInternalServerError() *GetOrganizationInternalServerError {

	return &GetOrganizationInternalServerError{}
}

// WithPayload adds the payload to the get organization internal server error response
func (o *GetOrganizationInternalServerError) WithPayload(payload interface{}) *GetOrganizationInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get organization internal server error response
func (o *GetOrganizationInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOrganizationInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
