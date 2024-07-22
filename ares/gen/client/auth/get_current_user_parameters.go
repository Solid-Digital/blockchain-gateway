// Code generated by go-swagger; DO NOT EDIT.

package auth

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

// NewGetCurrentUserParams creates a new GetCurrentUserParams object
// with the default values initialized.
func NewGetCurrentUserParams() *GetCurrentUserParams {

	return &GetCurrentUserParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetCurrentUserParamsWithTimeout creates a new GetCurrentUserParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetCurrentUserParamsWithTimeout(timeout time.Duration) *GetCurrentUserParams {

	return &GetCurrentUserParams{

		timeout: timeout,
	}
}

// NewGetCurrentUserParamsWithContext creates a new GetCurrentUserParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetCurrentUserParamsWithContext(ctx context.Context) *GetCurrentUserParams {

	return &GetCurrentUserParams{

		Context: ctx,
	}
}

// NewGetCurrentUserParamsWithHTTPClient creates a new GetCurrentUserParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetCurrentUserParamsWithHTTPClient(client *http.Client) *GetCurrentUserParams {

	return &GetCurrentUserParams{
		HTTPClient: client,
	}
}

/*GetCurrentUserParams contains all the parameters to send to the API endpoint
for the get current user operation typically these are written to a http.Request
*/
type GetCurrentUserParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get current user params
func (o *GetCurrentUserParams) WithTimeout(timeout time.Duration) *GetCurrentUserParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get current user params
func (o *GetCurrentUserParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get current user params
func (o *GetCurrentUserParams) WithContext(ctx context.Context) *GetCurrentUserParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get current user params
func (o *GetCurrentUserParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get current user params
func (o *GetCurrentUserParams) WithHTTPClient(client *http.Client) *GetCurrentUserParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get current user params
func (o *GetCurrentUserParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetCurrentUserParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
