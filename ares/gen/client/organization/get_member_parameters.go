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

// NewGetMemberParams creates a new GetMemberParams object
// with the default values initialized.
func NewGetMemberParams() *GetMemberParams {
	var ()
	return &GetMemberParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetMemberParamsWithTimeout creates a new GetMemberParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetMemberParamsWithTimeout(timeout time.Duration) *GetMemberParams {
	var ()
	return &GetMemberParams{

		timeout: timeout,
	}
}

// NewGetMemberParamsWithContext creates a new GetMemberParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetMemberParamsWithContext(ctx context.Context) *GetMemberParams {
	var ()
	return &GetMemberParams{

		Context: ctx,
	}
}

// NewGetMemberParamsWithHTTPClient creates a new GetMemberParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetMemberParamsWithHTTPClient(client *http.Client) *GetMemberParams {
	var ()
	return &GetMemberParams{
		HTTPClient: client,
	}
}

/*GetMemberParams contains all the parameters to send to the API endpoint
for the get member operation typically these are written to a http.Request
*/
type GetMemberParams struct {

	/*Email*/
	Email strfmt.Email
	/*OrgName*/
	OrgName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get member params
func (o *GetMemberParams) WithTimeout(timeout time.Duration) *GetMemberParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get member params
func (o *GetMemberParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get member params
func (o *GetMemberParams) WithContext(ctx context.Context) *GetMemberParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get member params
func (o *GetMemberParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get member params
func (o *GetMemberParams) WithHTTPClient(client *http.Client) *GetMemberParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get member params
func (o *GetMemberParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEmail adds the email to the get member params
func (o *GetMemberParams) WithEmail(email strfmt.Email) *GetMemberParams {
	o.SetEmail(email)
	return o
}

// SetEmail adds the email to the get member params
func (o *GetMemberParams) SetEmail(email strfmt.Email) {
	o.Email = email
}

// WithOrgName adds the orgName to the get member params
func (o *GetMemberParams) WithOrgName(orgName string) *GetMemberParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the get member params
func (o *GetMemberParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WriteToRequest writes these params to a swagger request
func (o *GetMemberParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param email
	if err := r.SetPathParam("email", o.Email.String()); err != nil {
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
