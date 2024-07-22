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

// NewGetPublicActionParams creates a new GetPublicActionParams object
// with the default values initialized.
func NewGetPublicActionParams() *GetPublicActionParams {
	var ()
	return &GetPublicActionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPublicActionParamsWithTimeout creates a new GetPublicActionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPublicActionParamsWithTimeout(timeout time.Duration) *GetPublicActionParams {
	var ()
	return &GetPublicActionParams{

		timeout: timeout,
	}
}

// NewGetPublicActionParamsWithContext creates a new GetPublicActionParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPublicActionParamsWithContext(ctx context.Context) *GetPublicActionParams {
	var ()
	return &GetPublicActionParams{

		Context: ctx,
	}
}

// NewGetPublicActionParamsWithHTTPClient creates a new GetPublicActionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPublicActionParamsWithHTTPClient(client *http.Client) *GetPublicActionParams {
	var ()
	return &GetPublicActionParams{
		HTTPClient: client,
	}
}

/*GetPublicActionParams contains all the parameters to send to the API endpoint
for the get public action operation typically these are written to a http.Request
*/
type GetPublicActionParams struct {

	/*Name*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get public action params
func (o *GetPublicActionParams) WithTimeout(timeout time.Duration) *GetPublicActionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get public action params
func (o *GetPublicActionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get public action params
func (o *GetPublicActionParams) WithContext(ctx context.Context) *GetPublicActionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get public action params
func (o *GetPublicActionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get public action params
func (o *GetPublicActionParams) WithHTTPClient(client *http.Client) *GetPublicActionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get public action params
func (o *GetPublicActionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get public action params
func (o *GetPublicActionParams) WithName(name string) *GetPublicActionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get public action params
func (o *GetPublicActionParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *GetPublicActionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
