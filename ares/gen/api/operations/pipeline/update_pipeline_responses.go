// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// UpdatePipelineOKCode is the HTTP code returned for type UpdatePipelineOK
const UpdatePipelineOKCode int = 200

/*UpdatePipelineOK Status 200

swagger:response updatePipelineOK
*/
type UpdatePipelineOK struct {

	/*
	  In: Body
	*/
	Payload *dto.GetPipelineResponse `json:"body,omitempty"`
}

// NewUpdatePipelineOK creates UpdatePipelineOK with default headers values
func NewUpdatePipelineOK() *UpdatePipelineOK {

	return &UpdatePipelineOK{}
}

// WithPayload adds the payload to the update pipeline o k response
func (o *UpdatePipelineOK) WithPayload(payload *dto.GetPipelineResponse) *UpdatePipelineOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update pipeline o k response
func (o *UpdatePipelineOK) SetPayload(payload *dto.GetPipelineResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePipelineOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdatePipelineInternalServerErrorCode is the HTTP code returned for type UpdatePipelineInternalServerError
const UpdatePipelineInternalServerErrorCode int = 500

/*UpdatePipelineInternalServerError Internal server error

swagger:response updatePipelineInternalServerError
*/
type UpdatePipelineInternalServerError struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewUpdatePipelineInternalServerError creates UpdatePipelineInternalServerError with default headers values
func NewUpdatePipelineInternalServerError() *UpdatePipelineInternalServerError {

	return &UpdatePipelineInternalServerError{}
}

// WithPayload adds the payload to the update pipeline internal server error response
func (o *UpdatePipelineInternalServerError) WithPayload(payload interface{}) *UpdatePipelineInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update pipeline internal server error response
func (o *UpdatePipelineInternalServerError) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePipelineInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
