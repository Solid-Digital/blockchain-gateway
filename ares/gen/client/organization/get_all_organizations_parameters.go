// Code generated by go-swagger; DO NOT EDIT.

package organization

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

// NewGetAllOrganizationsParams creates a new GetAllOrganizationsParams object
// with the default values initialized.
func NewGetAllOrganizationsParams() *GetAllOrganizationsParams {

	return &GetAllOrganizationsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAllOrganizationsParamsWithTimeout creates a new GetAllOrganizationsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAllOrganizationsParamsWithTimeout(timeout time.Duration) *GetAllOrganizationsParams {

	return &GetAllOrganizationsParams{

		timeout: timeout,
	}
}

// NewGetAllOrganizationsParamsWithContext creates a new GetAllOrganizationsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAllOrganizationsParamsWithContext(ctx context.Context) *GetAllOrganizationsParams {

	return &GetAllOrganizationsParams{

		Context: ctx,
	}
}

// NewGetAllOrganizationsParamsWithHTTPClient creates a new GetAllOrganizationsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAllOrganizationsParamsWithHTTPClient(client *http.Client) *GetAllOrganizationsParams {

	return &GetAllOrganizationsParams{
		HTTPClient: client,
	}
}

/*GetAllOrganizationsParams contains all the parameters to send to the API endpoint
for the get all organizations operation typically these are written to a http.Request
*/
type GetAllOrganizationsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get all organizations params
func (o *GetAllOrganizationsParams) WithTimeout(timeout time.Duration) *GetAllOrganizationsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get all organizations params
func (o *GetAllOrganizationsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get all organizations params
func (o *GetAllOrganizationsParams) WithContext(ctx context.Context) *GetAllOrganizationsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get all organizations params
func (o *GetAllOrganizationsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get all organizations params
func (o *GetAllOrganizationsParams) WithHTTPClient(client *http.Client) *GetAllOrganizationsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get all organizations params
func (o *GetAllOrganizationsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetAllOrganizationsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
