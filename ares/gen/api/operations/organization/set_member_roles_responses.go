// Code generated by go-swagger; DO NOT EDIT.

package organization

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// SetMemberRolesOKCode is the HTTP code returned for type SetMemberRolesOK
const SetMemberRolesOKCode int = 200

/*SetMemberRolesOK Status 200

swagger:response setMemberRolesOK
*/
type SetMemberRolesOK struct {
}

// NewSetMemberRolesOK creates SetMemberRolesOK with default headers values
func NewSetMemberRolesOK() *SetMemberRolesOK {

	return &SetMemberRolesOK{}
}

// WriteResponse to the client
func (o *SetMemberRolesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// SetMemberRolesInternalServerErrorCode is the HTTP code returned for type SetMemberRolesInternalServerError
const SetMemberRolesInternalServerErrorCode int = 500

/*SetMemberRolesInternalServerError Internal server error

swagger:response setMemberRolesInternalServerError
*/
type SetMemberRolesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewSetMemberRolesInternalServerError creates SetMemberRolesInternalServerError with default headers values
func NewSetMemberRolesInternalServerError() *SetMemberRolesInternalServerError {

	return &SetMemberRolesInternalServerError{}
}

// WithPayload adds the payload to the set member roles internal server error response
func (o *SetMemberRolesInternalServerError) WithPayload(payload interface{}) *SetMemberRolesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set member roles internal server error response
func (o *SetMemberRolesInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetMemberRolesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
