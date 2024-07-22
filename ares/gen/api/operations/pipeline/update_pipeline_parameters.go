// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// NewUpdatePipelineParams creates a new UpdatePipelineParams object
// no default values defined in spec.
func NewUpdatePipelineParams() UpdatePipelineParams {

	return UpdatePipelineParams{}
}

// UpdatePipelineParams contains all the bound params for the update pipeline operation
// typically these are obtained from a http.Request
//
// swagger:parameters UpdatePipeline
type UpdatePipelineParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Body *dto.UpdatePipelineRequest
	/*
	  Required: true
	  In: path
	*/
	OrgName string
	/*
	  Required: true
	  In: path
	*/
	PipelineName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdatePipelineParams() beforehand.
func (o *UpdatePipelineParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body dto.UpdatePipelineRequest
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body"))
	}
	rOrgName, rhkOrgName, _ := route.Params.GetOK("orgName")
	if err := o.bindOrgName(rOrgName, rhkOrgName, route.Formats); err != nil {
		res = append(res, err)
	}

	rPipelineName, rhkPipelineName, _ := route.Params.GetOK("pipelineName")
	if err := o.bindPipelineName(rPipelineName, rhkPipelineName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindOrgName binds and validates parameter OrgName from path.
func (o *UpdatePipelineParams) bindOrgName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.OrgName = raw

	return nil
}

// bindPipelineName binds and validates parameter PipelineName from path.
func (o *UpdatePipelineParams) bindPipelineName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.PipelineName = raw

	return nil
}
