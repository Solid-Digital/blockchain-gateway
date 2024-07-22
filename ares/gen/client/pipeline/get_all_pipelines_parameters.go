// Code generated by go-swagger; DO NOT EDIT.

package pipeline

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

// NewGetAllPipelinesParams creates a new GetAllPipelinesParams object
// with the default values initialized.
func NewGetAllPipelinesParams() *GetAllPipelinesParams {
	var ()
	return &GetAllPipelinesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAllPipelinesParamsWithTimeout creates a new GetAllPipelinesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAllPipelinesParamsWithTimeout(timeout time.Duration) *GetAllPipelinesParams {
	var ()
	return &GetAllPipelinesParams{

		timeout: timeout,
	}
}

// NewGetAllPipelinesParamsWithContext creates a new GetAllPipelinesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAllPipelinesParamsWithContext(ctx context.Context) *GetAllPipelinesParams {
	var ()
	return &GetAllPipelinesParams{

		Context: ctx,
	}
}

// NewGetAllPipelinesParamsWithHTTPClient creates a new GetAllPipelinesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAllPipelinesParamsWithHTTPClient(client *http.Client) *GetAllPipelinesParams {
	var ()
	return &GetAllPipelinesParams{
		HTTPClient: client,
	}
}

/*GetAllPipelinesParams contains all the parameters to send to the API endpoint
for the get all pipelines operation typically these are written to a http.Request
*/
type GetAllPipelinesParams struct {

	/*OrgName*/
	OrgName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get all pipelines params
func (o *GetAllPipelinesParams) WithTimeout(timeout time.Duration) *GetAllPipelinesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get all pipelines params
func (o *GetAllPipelinesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get all pipelines params
func (o *GetAllPipelinesParams) WithContext(ctx context.Context) *GetAllPipelinesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get all pipelines params
func (o *GetAllPipelinesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get all pipelines params
func (o *GetAllPipelinesParams) WithHTTPClient(client *http.Client) *GetAllPipelinesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get all pipelines params
func (o *GetAllPipelinesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOrgName adds the orgName to the get all pipelines params
func (o *GetAllPipelinesParams) WithOrgName(orgName string) *GetAllPipelinesParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the get all pipelines params
func (o *GetAllPipelinesParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WriteToRequest writes these params to a swagger request
func (o *GetAllPipelinesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param orgName
	if err := r.SetPathParam("orgName", o.OrgName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
