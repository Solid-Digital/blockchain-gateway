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

// NewGetAllEnvironmentVariablesParams creates a new GetAllEnvironmentVariablesParams object
// with the default values initialized.
func NewGetAllEnvironmentVariablesParams() *GetAllEnvironmentVariablesParams {
	var ()
	return &GetAllEnvironmentVariablesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAllEnvironmentVariablesParamsWithTimeout creates a new GetAllEnvironmentVariablesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAllEnvironmentVariablesParamsWithTimeout(timeout time.Duration) *GetAllEnvironmentVariablesParams {
	var ()
	return &GetAllEnvironmentVariablesParams{

		timeout: timeout,
	}
}

// NewGetAllEnvironmentVariablesParamsWithContext creates a new GetAllEnvironmentVariablesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAllEnvironmentVariablesParamsWithContext(ctx context.Context) *GetAllEnvironmentVariablesParams {
	var ()
	return &GetAllEnvironmentVariablesParams{

		Context: ctx,
	}
}

// NewGetAllEnvironmentVariablesParamsWithHTTPClient creates a new GetAllEnvironmentVariablesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAllEnvironmentVariablesParamsWithHTTPClient(client *http.Client) *GetAllEnvironmentVariablesParams {
	var ()
	return &GetAllEnvironmentVariablesParams{
		HTTPClient: client,
	}
}

/*GetAllEnvironmentVariablesParams contains all the parameters to send to the API endpoint
for the get all environment variables operation typically these are written to a http.Request
*/
type GetAllEnvironmentVariablesParams struct {

	/*EnvName*/
	EnvName string
	/*OrgName*/
	OrgName string
	/*PipelineName*/
	PipelineName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) WithTimeout(timeout time.Duration) *GetAllEnvironmentVariablesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) WithContext(ctx context.Context) *GetAllEnvironmentVariablesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) WithHTTPClient(client *http.Client) *GetAllEnvironmentVariablesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEnvName adds the envName to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) WithEnvName(envName string) *GetAllEnvironmentVariablesParams {
	o.SetEnvName(envName)
	return o
}

// SetEnvName adds the envName to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) SetEnvName(envName string) {
	o.EnvName = envName
}

// WithOrgName adds the orgName to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) WithOrgName(orgName string) *GetAllEnvironmentVariablesParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WithPipelineName adds the pipelineName to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) WithPipelineName(pipelineName string) *GetAllEnvironmentVariablesParams {
	o.SetPipelineName(pipelineName)
	return o
}

// SetPipelineName adds the pipelineName to the get all environment variables params
func (o *GetAllEnvironmentVariablesParams) SetPipelineName(pipelineName string) {
	o.PipelineName = pipelineName
}

// WriteToRequest writes these params to a swagger request
func (o *GetAllEnvironmentVariablesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param envName
	if err := r.SetPathParam("envName", o.EnvName); err != nil {
		return err
	}

	// path param orgName
	if err := r.SetPathParam("orgName", o.OrgName); err != nil {
		return err
	}

	// path param pipelineName
	if err := r.SetPathParam("pipelineName", o.PipelineName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
