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
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewCreateActionVersionParams creates a new CreateActionVersionParams object
// with the default values initialized.
func NewCreateActionVersionParams() *CreateActionVersionParams {
	var ()
	return &CreateActionVersionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateActionVersionParamsWithTimeout creates a new CreateActionVersionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateActionVersionParamsWithTimeout(timeout time.Duration) *CreateActionVersionParams {
	var ()
	return &CreateActionVersionParams{

		timeout: timeout,
	}
}

// NewCreateActionVersionParamsWithContext creates a new CreateActionVersionParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateActionVersionParamsWithContext(ctx context.Context) *CreateActionVersionParams {
	var ()
	return &CreateActionVersionParams{

		Context: ctx,
	}
}

// NewCreateActionVersionParamsWithHTTPClient creates a new CreateActionVersionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateActionVersionParamsWithHTTPClient(client *http.Client) *CreateActionVersionParams {
	var ()
	return &CreateActionVersionParams{
		HTTPClient: client,
	}
}

/*CreateActionVersionParams contains all the parameters to send to the API endpoint
for the create action version operation typically these are written to a http.Request
*/
type CreateActionVersionParams struct {

	/*ActionFile
	  the action version file

	*/
	ActionFile runtime.NamedReadCloser
	/*Description
	  short description of this action version

	*/
	Description *string
	/*ExampleConfig
	  default config for this action version

	*/
	ExampleConfig *string
	/*InputSchema
	  json encoded string containing the input specification of the action

	*/
	InputSchema *string
	/*Name*/
	Name string
	/*OrgName*/
	OrgName string
	/*OutputSchema
	  json encoded string containing the output specification of the action

	*/
	OutputSchema *string
	/*Public
	  describes whether or not this action version is public

	*/
	Public *bool
	/*Readme
	  readme for this action version

	*/
	Readme *string
	/*Version
	  version string for this action version

	*/
	Version *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create action version params
func (o *CreateActionVersionParams) WithTimeout(timeout time.Duration) *CreateActionVersionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create action version params
func (o *CreateActionVersionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create action version params
func (o *CreateActionVersionParams) WithContext(ctx context.Context) *CreateActionVersionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create action version params
func (o *CreateActionVersionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create action version params
func (o *CreateActionVersionParams) WithHTTPClient(client *http.Client) *CreateActionVersionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create action version params
func (o *CreateActionVersionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithActionFile adds the actionFile to the create action version params
func (o *CreateActionVersionParams) WithActionFile(actionFile runtime.NamedReadCloser) *CreateActionVersionParams {
	o.SetActionFile(actionFile)
	return o
}

// SetActionFile adds the actionFile to the create action version params
func (o *CreateActionVersionParams) SetActionFile(actionFile runtime.NamedReadCloser) {
	o.ActionFile = actionFile
}

// WithDescription adds the description to the create action version params
func (o *CreateActionVersionParams) WithDescription(description *string) *CreateActionVersionParams {
	o.SetDescription(description)
	return o
}

// SetDescription adds the description to the create action version params
func (o *CreateActionVersionParams) SetDescription(description *string) {
	o.Description = description
}

// WithExampleConfig adds the exampleConfig to the create action version params
func (o *CreateActionVersionParams) WithExampleConfig(exampleConfig *string) *CreateActionVersionParams {
	o.SetExampleConfig(exampleConfig)
	return o
}

// SetExampleConfig adds the exampleConfig to the create action version params
func (o *CreateActionVersionParams) SetExampleConfig(exampleConfig *string) {
	o.ExampleConfig = exampleConfig
}

// WithInputSchema adds the inputSchema to the create action version params
func (o *CreateActionVersionParams) WithInputSchema(inputSchema *string) *CreateActionVersionParams {
	o.SetInputSchema(inputSchema)
	return o
}

// SetInputSchema adds the inputSchema to the create action version params
func (o *CreateActionVersionParams) SetInputSchema(inputSchema *string) {
	o.InputSchema = inputSchema
}

// WithName adds the name to the create action version params
func (o *CreateActionVersionParams) WithName(name string) *CreateActionVersionParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the create action version params
func (o *CreateActionVersionParams) SetName(name string) {
	o.Name = name
}

// WithOrgName adds the orgName to the create action version params
func (o *CreateActionVersionParams) WithOrgName(orgName string) *CreateActionVersionParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the create action version params
func (o *CreateActionVersionParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WithOutputSchema adds the outputSchema to the create action version params
func (o *CreateActionVersionParams) WithOutputSchema(outputSchema *string) *CreateActionVersionParams {
	o.SetOutputSchema(outputSchema)
	return o
}

// SetOutputSchema adds the outputSchema to the create action version params
func (o *CreateActionVersionParams) SetOutputSchema(outputSchema *string) {
	o.OutputSchema = outputSchema
}

// WithPublic adds the public to the create action version params
func (o *CreateActionVersionParams) WithPublic(public *bool) *CreateActionVersionParams {
	o.SetPublic(public)
	return o
}

// SetPublic adds the public to the create action version params
func (o *CreateActionVersionParams) SetPublic(public *bool) {
	o.Public = public
}

// WithReadme adds the readme to the create action version params
func (o *CreateActionVersionParams) WithReadme(readme *string) *CreateActionVersionParams {
	o.SetReadme(readme)
	return o
}

// SetReadme adds the readme to the create action version params
func (o *CreateActionVersionParams) SetReadme(readme *string) {
	o.Readme = readme
}

// WithVersion adds the version to the create action version params
func (o *CreateActionVersionParams) WithVersion(version *string) *CreateActionVersionParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the create action version params
func (o *CreateActionVersionParams) SetVersion(version *string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *CreateActionVersionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ActionFile != nil {

		if o.ActionFile != nil {

			// form file param actionFile
			if err := r.SetFileParam("actionFile", o.ActionFile); err != nil {
				return err
			}

		}

	}

	if o.Description != nil {

		// form param description
		var frDescription string
		if o.Description != nil {
			frDescription = *o.Description
		}
		fDescription := frDescription
		if fDescription != "" {
			if err := r.SetFormParam("description", fDescription); err != nil {
				return err
			}
		}

	}

	if o.ExampleConfig != nil {

		// form param exampleConfig
		var frExampleConfig string
		if o.ExampleConfig != nil {
			frExampleConfig = *o.ExampleConfig
		}
		fExampleConfig := frExampleConfig
		if fExampleConfig != "" {
			if err := r.SetFormParam("exampleConfig", fExampleConfig); err != nil {
				return err
			}
		}

	}

	if o.InputSchema != nil {

		// form param inputSchema
		var frInputSchema string
		if o.InputSchema != nil {
			frInputSchema = *o.InputSchema
		}
		fInputSchema := frInputSchema
		if fInputSchema != "" {
			if err := r.SetFormParam("inputSchema", fInputSchema); err != nil {
				return err
			}
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

	if o.OutputSchema != nil {

		// form param outputSchema
		var frOutputSchema string
		if o.OutputSchema != nil {
			frOutputSchema = *o.OutputSchema
		}
		fOutputSchema := frOutputSchema
		if fOutputSchema != "" {
			if err := r.SetFormParam("outputSchema", fOutputSchema); err != nil {
				return err
			}
		}

	}

	if o.Public != nil {

		// form param public
		var frPublic bool
		if o.Public != nil {
			frPublic = *o.Public
		}
		fPublic := swag.FormatBool(frPublic)
		if fPublic != "" {
			if err := r.SetFormParam("public", fPublic); err != nil {
				return err
			}
		}

	}

	if o.Readme != nil {

		// form param readme
		var frReadme string
		if o.Readme != nil {
			frReadme = *o.Readme
		}
		fReadme := frReadme
		if fReadme != "" {
			if err := r.SetFormParam("readme", fReadme); err != nil {
				return err
			}
		}

	}

	if o.Version != nil {

		// form param version
		var frVersion string
		if o.Version != nil {
			frVersion = *o.Version
		}
		fVersion := frVersion
		if fVersion != "" {
			if err := r.SetFormParam("version", fVersion); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
