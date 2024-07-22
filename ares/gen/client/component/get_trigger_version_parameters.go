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

// NewGetTriggerVersionParams creates a new GetTriggerVersionParams object
// with the default values initialized.
func NewGetTriggerVersionParams() *GetTriggerVersionParams {
	var ()
	return &GetTriggerVersionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetTriggerVersionParamsWithTimeout creates a new GetTriggerVersionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetTriggerVersionParamsWithTimeout(timeout time.Duration) *GetTriggerVersionParams {
	var ()
	return &GetTriggerVersionParams{

		timeout: timeout,
	}
}

// NewGetTriggerVersionParamsWithContext creates a new GetTriggerVersionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetTriggerVersionParamsWithContext(ctx context.Context) *GetTriggerVersionParams {
	var ()
	return &GetTriggerVersionParams{

		Context: ctx,
	}
}

// NewGetTriggerVersionParamsWithHTTPClient creates a new GetTriggerVersionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetTriggerVersionParamsWithHTTPClient(client *http.Client) *GetTriggerVersionParams {
	var ()
	return &GetTriggerVersionParams{
		HTTPClient: client,
	}
}

/*GetTriggerVersionParams contains all the parameters to send to the API endpoint
for the get trigger version operation typically these are written to a http.Request
*/
type GetTriggerVersionParams struct {

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

// WithTimeout adds the timeout to the get trigger version params
func (o *GetTriggerVersionParams) WithTimeout(timeout time.Duration) *GetTriggerVersionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get trigger version params
func (o *GetTriggerVersionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get trigger version params
func (o *GetTriggerVersionParams) WithContext(ctx context.Context) *GetTriggerVersionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get trigger version params
func (o *GetTriggerVersionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get trigger version params
func (o *GetTriggerVersionParams) WithHTTPClient(client *http.Client) *GetTriggerVersionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get trigger version params
func (o *GetTriggerVersionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get trigger version params
func (o *GetTriggerVersionParams) WithName(name string) *GetTriggerVersionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get trigger version params
func (o *GetTriggerVersionParams) SetName(name string) {
	o.Name = name
}

// WithOrgName adds the orgName to the get trigger version params
func (o *GetTriggerVersionParams) WithOrgName(orgName string) *GetTriggerVersionParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the get trigger version params
func (o *GetTriggerVersionParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WithVersion adds the version to the get trigger version params
func (o *GetTriggerVersionParams) WithVersion(version string) *GetTriggerVersionParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the get trigger version params
func (o *GetTriggerVersionParams) SetVersion(version string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *GetTriggerVersionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
