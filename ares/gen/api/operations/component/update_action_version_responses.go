// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// UpdateActionVersionOKCode is the HTTP code returned for type UpdateActionVersionOK
const UpdateActionVersionOKCode int = 200

/*UpdateActionVersionOK Status 200

swagger:response updateActionVersionOK
*/
type UpdateActionVersionOK struct {

	/*
	  In: Body
	*/
	Payload *dto.GetComponentVersionResponse `json:"body,omitempty"`
}

// NewUpdateActionVersionOK creates UpdateActionVersionOK with default headers values
func NewUpdateActionVersionOK() *UpdateActionVersionOK {

	return &UpdateActionVersionOK{}
}

// WithPayload adds the payload to the update action version o k response
func (o *UpdateActionVersionOK) WithPayload(payload *dto.GetComponentVersionResponse) *UpdateActionVersionOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update action version o k response
func (o *UpdateActionVersionOK) SetPayload(payload *dto.GetComponentVersionResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateActionVersionOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateActionVersionInternalServerErrorCode is the HTTP code returned for type UpdateActionVersionInternalServerError
const UpdateActionVersionInternalServerErrorCode int = 500

/*UpdateActionVersionInternalServerError Internal server error

swagger:response updateActionVersionInternalServerError
*/
type UpdateActionVersionInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewUpdateActionVersionInternalServerError creates UpdateActionVersionInternalServerError with default headers values
func NewUpdateActionVersionInternalServerError() *UpdateActionVersionInternalServerError {

	return &UpdateActionVersionInternalServerError{}
}

// WithPayload adds the payload to the update action version internal server error response
func (o *UpdateActionVersionInternalServerError) WithPayload(payload interface{}) *UpdateActionVersionInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update action version internal server error response
func (o *UpdateActionVersionInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateActionVersionInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
