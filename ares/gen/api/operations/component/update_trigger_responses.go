// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// UpdateTriggerOKCode is the HTTP code returned for type UpdateTriggerOK
const UpdateTriggerOKCode int = 200

/*UpdateTriggerOK Status 200

swagger:response updateTriggerOK
*/
type UpdateTriggerOK struct {

	/*
	  In: Body
	*/
	Payload *dto.GetComponentResponse `json:"body,omitempty"`
}

// NewUpdateTriggerOK creates UpdateTriggerOK with default headers values
func NewUpdateTriggerOK() *UpdateTriggerOK {

	return &UpdateTriggerOK{}
}

// WithPayload adds the payload to the update trigger o k response
func (o *UpdateTriggerOK) WithPayload(payload *dto.GetComponentResponse) *UpdateTriggerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger o k response
func (o *UpdateTriggerOK) SetPayload(payload *dto.GetComponentResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerConflictCode is the HTTP code returned for type UpdateTriggerConflict
const UpdateTriggerConflictCode int = 409

/*UpdateTriggerConflict Conflict

swagger:response updateTriggerConflict
*/
type UpdateTriggerConflict struct {

	/*
	  In: Body
	*/
	Payload *dto.ErrorResponse `json:"body,omitempty"`
}

// NewUpdateTriggerConflict creates UpdateTriggerConflict with default headers values
func NewUpdateTriggerConflict() *UpdateTriggerConflict {

	return &UpdateTriggerConflict{}
}

// WithPayload adds the payload to the update trigger conflict response
func (o *UpdateTriggerConflict) WithPayload(payload *dto.ErrorResponse) *UpdateTriggerConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger conflict response
func (o *UpdateTriggerConflict) SetPayload(payload *dto.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerInternalServerErrorCode is the HTTP code returned for type UpdateTriggerInternalServerError
const UpdateTriggerInternalServerErrorCode int = 500

/*UpdateTriggerInternalServerError Internal server error

swagger:response updateTriggerInternalServerError
*/
type UpdateTriggerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *dto.ErrorResponse `json:"body,omitempty"`
}

// NewUpdateTriggerInternalServerError creates UpdateTriggerInternalServerError with default headers values
func NewUpdateTriggerInternalServerError() *UpdateTriggerInternalServerError {

	return &UpdateTriggerInternalServerError{}
}

// WithPayload adds the payload to the update trigger internal server error response
func (o *UpdateTriggerInternalServerError) WithPayload(payload *dto.ErrorResponse) *UpdateTriggerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger internal server error response
func (o *UpdateTriggerInternalServerError) SetPayload(payload *dto.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
