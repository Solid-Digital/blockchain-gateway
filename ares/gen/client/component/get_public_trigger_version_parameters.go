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

// NewGetPublicTriggerVersionParams creates a new GetPublicTriggerVersionParams object
// with the default values initialized.
func NewGetPublicTriggerVersionParams() *GetPublicTriggerVersionParams {
	var ()
	return &GetPublicTriggerVersionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPublicTriggerVersionParamsWithTimeout creates a new GetPublicTriggerVersionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPublicTriggerVersionParamsWithTimeout(timeout time.Duration) *GetPublicTriggerVersionParams {
	var ()
	return &GetPublicTriggerVersionParams{

		timeout: timeout,
	}
}

// NewGetPublicTriggerVersionParamsWithContext creates a new GetPublicTriggerVersionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPublicTriggerVersionParamsWithContext(ctx context.Context) *GetPublicTriggerVersionParams {
	var ()
	return &GetPublicTriggerVersionParams{

		Context: ctx,
	}
}

// NewGetPublicTriggerVersionParamsWithHTTPClient creates a new GetPublicTriggerVersionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPublicTriggerVersionParamsWithHTTPClient(client *http.Client) *GetPublicTriggerVersionParams {
	var ()
	return &GetPublicTriggerVersionParams{
		HTTPClient: client,
	}
}

/*GetPublicTriggerVersionParams contains all the parameters to send to the API endpoint
for the get public trigger version operation typically these are written to a http.Request
*/
type GetPublicTriggerVersionParams struct {

	/*Name*/
	Name string
	/*Version*/
	Version string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get public trigger version params
func (o *GetPublicTriggerVersionParams) WithTimeout(timeout time.Duration) *GetPublicTriggerVersionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get public trigger version params
func (o *GetPublicTriggerVersionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get public trigger version params
func (o *GetPublicTriggerVersionParams) WithContext(ctx context.Context) *GetPublicTriggerVersionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get public trigger version params
func (o *GetPublicTriggerVersionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get public trigger version params
func (o *GetPublicTriggerVersionParams) WithHTTPClient(client *http.Client) *GetPublicTriggerVersionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get public trigger version params
func (o *GetPublicTriggerVersionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get public trigger version params
func (o *GetPublicTriggerVersionParams) WithName(name string) *GetPublicTriggerVersionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get public trigger version params
func (o *GetPublicTriggerVersionParams) SetName(name string) {
	o.Name = name
}

// WithVersion adds the version to the get public trigger version params
func (o *GetPublicTriggerVersionParams) WithVersion(version string) *GetPublicTriggerVersionParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the get public trigger version params
func (o *GetPublicTriggerVersionParams) SetVersion(version string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *GetPublicTriggerVersionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
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
