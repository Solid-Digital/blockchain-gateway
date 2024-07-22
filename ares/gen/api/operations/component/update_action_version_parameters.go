// Code generated by go-swagger; DO NOT EDIT.

package component

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

// NewUpdateActionVersionParams creates a new UpdateActionVersionParams object
// no default values defined in spec.
func NewUpdateActionVersionParams() UpdateActionVersionParams {

	return UpdateActionVersionParams{}
}

// UpdateActionVersionParams contains all the bound params for the update action version operation
// typically these are obtained from a http.Request
//
// swagger:parameters UpdateActionVersion
type UpdateActionVersionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: body
	*/
	Body *dto.UpdateComponentVersionRequest
	/*
	  Required: true
	  In: path
	*/
	Name string
	/*
	  Required: true
	  In: path
	*/
	OrgName string
	/*
	  Required: true
	  In: path
	*/
	Version string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdateActionVersionParams() beforehand.
func (o *UpdateActionVersionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body dto.UpdateComponentVersionRequest
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
	rName, rhkName, _ := route.Params.GetOK("name")
	if err := o.bindName(rName, rhkName, route.Formats); err != nil {
		res = append(res, err)
	}

	rOrgName, rhkOrgName, _ := route.Params.GetOK("orgName")
	if err := o.bindOrgName(rOrgName, rhkOrgName, route.Formats); err != nil {
		res = append(res, err)
	}

	rVersion, rhkVersion, _ := route.Params.GetOK("version")
	if err := o.bindVersion(rVersion, rhkVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindName binds and validates parameter Name from path.
func (o *UpdateActionVersionParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Name = raw

	return nil
}

// bindOrgName binds and validates parameter OrgName from path.
func (o *UpdateActionVersionParams) bindOrgName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.OrgName = raw

	return nil
}

// bindVersion binds and validates parameter Version from path.
func (o *UpdateActionVersionParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Version = raw

	return nil
}
