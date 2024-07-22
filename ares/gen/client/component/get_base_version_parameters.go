// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetBaseVersionParams creates a new GetBaseVersionParams object
// with the default values initialized.
func NewGetBaseVersionParams() *GetBaseVersionParams {
	var ()
	return &GetBaseVersionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetBaseVersionParamsWithTimeout creates a new GetBaseVersionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetBaseVersionParamsWithTimeout(timeout time.Duration) *GetBaseVersionParams {
	var ()
	return &GetBaseVersionParams{

		timeout: timeout,
	}
}

// NewGetBaseVersionParamsWithContext creates a new GetBaseVersionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetBaseVersionParamsWithContext(ctx context.Context) *GetBaseVersionParams {
	var ()
	return &GetBaseVersionParams{

		Context: ctx,
	}
}

// NewGetBaseVersionParamsWithHTTPClient creates a new GetBaseVersionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetBaseVersionParamsWithHTTPClient(client *http.Client) *GetBaseVersionParams {
	var ()
	return &GetBaseVersionParams{
		HTTPClient: client,
	}
}

/*GetBaseVersionParams contains all the parameters to send to the API endpoint
for the get base version operation typically these are written to a http.Request
*/
type GetBaseVersionParams struct {

	/*Name*/
	Name string
	/*OrgName*/
	OrgName string
	/*Version*/
	Version string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get base version params
func (o *GetBaseVersionParams) WithTimeout(timeout time.Duration) *GetBaseVersionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get base version params
func (o *GetBaseVersionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get base version params
func (o *GetBaseVersionParams) WithContext(ctx context.Context) *GetBaseVersionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get base version params
func (o *GetBaseVersionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get base version params
func (o *GetBaseVersionParams) WithHTTPClient(client *http.Client) *GetBaseVersionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get base version params
func (o *GetBaseVersionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get base version params
func (o *GetBaseVersionParams) WithName(name string) *GetBaseVersionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get base version params
func (o *GetBaseVersionParams) SetName(name string) {
	o.Name = name
}

// WithOrgName adds the orgName to the get base version params
func (o *GetBaseVersionParams) WithOrgName(orgName string) *GetBaseVersionParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the get base version params
func (o *GetBaseVersionParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WithVersion adds the version to the get base version params
func (o *GetBaseVersionParams) WithVersion(version string) *GetBaseVersionParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the get base version params
func (o *GetBaseVersionParams) SetVersion(version string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *GetBaseVersionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	// path param orgName
	if err := r.SetPathParam("orgName", o.OrgName); err != nil {
		return err
	}

	// path param version
	if err := r.SetPathParam("version", o.Version); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
