// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetPublicActionVersionParams creates a new GetPublicActionVersionParams object
// no default values defined in spec.
func NewGetPublicActionVersionParams() GetPublicActionVersionParams {

	return GetPublicActionVersionParams{}
}

// GetPublicActionVersionParams contains all the bound params for the get public action version operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetPublicActionVersion
type GetPublicActionVersionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	Name string
	/*
	  Required: true
	  In: path
	*/
	Version string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetPublicActionVersionParams() beforehand.
func (o *GetPublicActionVersionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rName, rhkName, _ := route.Params.GetOK("name")
	if err := o.bindName(rName, rhkName, route.Formats); err != nil {
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
func (o *GetPublicActionVersionParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Name = raw

	return nil
}

// bindVersion binds and validates parameter Version from path.
func (o *GetPublicActionVersionParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.Version = raw

	return nil
}
