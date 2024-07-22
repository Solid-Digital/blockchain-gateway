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

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// NewCreateBaseVersionParams creates a new CreateBaseVersionParams object
// with the default values initialized.
func NewCreateBaseVersionParams() *CreateBaseVersionParams {
	var ()
	return &CreateBaseVersionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateBaseVersionParamsWithTimeout creates a new CreateBaseVersionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateBaseVersionParamsWithTimeout(timeout time.Duration) *CreateBaseVersionParams {
	var ()
	return &CreateBaseVersionParams{

		timeout: timeout,
	}
}

// NewCreateBaseVersionParamsWithContext creates a new CreateBaseVersionParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateBaseVersionParamsWithContext(ctx context.Context) *CreateBaseVersionParams {
	var ()
	return &CreateBaseVersionParams{

		Context: ctx,
	}
}

// NewCreateBaseVersionParamsWithHTTPClient creates a new CreateBaseVersionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateBaseVersionParamsWithHTTPClient(client *http.Client) *CreateBaseVersionParams {
	var ()
	return &CreateBaseVersionParams{
		HTTPClient: client,
	}
}

/*CreateBaseVersionParams contains all the parameters to send to the API endpoint
for the create base version operation typically these are written to a http.Request
*/
type CreateBaseVersionParams struct {

	/*Body*/
	Body *dto.CreateBaseVersionRequest
	/*Name*/
	Name string
	/*OrgName*/
	OrgName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create base version params
func (o *CreateBaseVersionParams) WithTimeout(timeout time.Duration) *CreateBaseVersionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create base version params
func (o *CreateBaseVersionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create base version params
func (o *CreateBaseVersionParams) WithContext(ctx context.Context) *CreateBaseVersionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create base version params
func (o *CreateBaseVersionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create base version params
func (o *CreateBaseVersionParams) WithHTTPClient(client *http.Client) *CreateBaseVersionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create base version params
func (o *CreateBaseVersionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create base version params
func (o *CreateBaseVersionParams) WithBody(body *dto.CreateBaseVersionRequest) *CreateBaseVersionParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create base version params
func (o *CreateBaseVersionParams) SetBody(body *dto.CreateBaseVersionRequest) {
	o.Body = body
}

// WithName adds the name to the create base version params
func (o *CreateBaseVersionParams) WithName(name string) *CreateBaseVersionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the create base version params
func (o *CreateBaseVersionParams) SetName(name string) {
	o.Name = name
}

// WithOrgName adds the orgName to the create base version params
func (o *CreateBaseVersionParams) WithOrgName(orgName string) *CreateBaseVersionParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the create base version params
func (o *CreateBaseVersionParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WriteToRequest writes these params to a swagger request
func (o *CreateBaseVersionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	// path param orgName
	if err := r.SetPathParam("orgName", o.OrgName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
