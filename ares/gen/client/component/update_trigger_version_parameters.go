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

// NewUpdateTriggerVersionParams creates a new UpdateTriggerVersionParams object
// with the default values initialized.
func NewUpdateTriggerVersionParams() *UpdateTriggerVersionParams {
	var ()
	return &UpdateTriggerVersionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateTriggerVersionParamsWithTimeout creates a new UpdateTriggerVersionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateTriggerVersionParamsWithTimeout(timeout time.Duration) *UpdateTriggerVersionParams {
	var ()
	return &UpdateTriggerVersionParams{

		timeout: timeout,
	}
}

// NewUpdateTriggerVersionParamsWithContext creates a new UpdateTriggerVersionParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateTriggerVersionParamsWithContext(ctx context.Context) *UpdateTriggerVersionParams {
	var ()
	return &UpdateTriggerVersionParams{

		Context: ctx,
	}
}

// NewUpdateTriggerVersionParamsWithHTTPClient creates a new UpdateTriggerVersionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateTriggerVersionParamsWithHTTPClient(client *http.Client) *UpdateTriggerVersionParams {
	var ()
	return &UpdateTriggerVersionParams{
		HTTPClient: client,
	}
}

/*UpdateTriggerVersionParams contains all the parameters to send to the API endpoint
for the update trigger version operation typically these are written to a http.Request
*/
type UpdateTriggerVersionParams struct {

	/*Body*/
	Body *dto.UpdateComponentVersionRequest
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

// WithTimeout adds the timeout to the update trigger version params
func (o *UpdateTriggerVersionParams) WithTimeout(timeout time.Duration) *UpdateTriggerVersionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update trigger version params
func (o *UpdateTriggerVersionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update trigger version params
func (o *UpdateTriggerVersionParams) WithContext(ctx context.Context) *UpdateTriggerVersionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update trigger version params
func (o *UpdateTriggerVersionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update trigger version params
func (o *UpdateTriggerVersionParams) WithHTTPClient(client *http.Client) *UpdateTriggerVersionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update trigger version params
func (o *UpdateTriggerVersionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update trigger version params
func (o *UpdateTriggerVersionParams) WithBody(body *dto.UpdateComponentVersionRequest) *UpdateTriggerVersionParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update trigger version params
func (o *UpdateTriggerVersionParams) SetBody(body *dto.UpdateComponentVersionRequest) {
	o.Body = body
}

// WithName adds the name to the update trigger version params
func (o *UpdateTriggerVersionParams) WithName(name string) *UpdateTriggerVersionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the update trigger version params
func (o *UpdateTriggerVersionParams) SetName(name string) {
	o.Name = name
}

// WithOrgName adds the orgName to the update trigger version params
func (o *UpdateTriggerVersionParams) WithOrgName(orgName string) *UpdateTriggerVersionParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the update trigger version params
func (o *UpdateTriggerVersionParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WithVersion adds the version to the update trigger version params
func (o *UpdateTriggerVersionParams) WithVersion(version string) *UpdateTriggerVersionParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the update trigger version params
func (o *UpdateTriggerVersionParams) SetVersion(version string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateTriggerVersionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param version
	if err := r.SetPathParam("version", o.Version); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
