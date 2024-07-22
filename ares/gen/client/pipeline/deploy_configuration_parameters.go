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

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// NewDeployConfigurationParams creates a new DeployConfigurationParams object
// with the default values initialized.
func NewDeployConfigurationParams() *DeployConfigurationParams {
	var ()
	return &DeployConfigurationParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeployConfigurationParamsWithTimeout creates a new DeployConfigurationParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeployConfigurationParamsWithTimeout(timeout time.Duration) *DeployConfigurationParams {
	var ()
	return &DeployConfigurationParams{

		timeout: timeout,
	}
}

// NewDeployConfigurationParamsWithContext creates a new DeployConfigurationParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeployConfigurationParamsWithContext(ctx context.Context) *DeployConfigurationParams {
	var ()
	return &DeployConfigurationParams{

		Context: ctx,
	}
}

// NewDeployConfigurationParamsWithHTTPClient creates a new DeployConfigurationParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeployConfigurationParamsWithHTTPClient(client *http.Client) *DeployConfigurationParams {
	var ()
	return &DeployConfigurationParams{
		HTTPClient: client,
	}
}

/*DeployConfigurationParams contains all the parameters to send to the API endpoint
for the deploy configuration operation typically these are written to a http.Request
*/
type DeployConfigurationParams struct {

	/*Body*/
	Body *dto.DeployConfigurationRequest
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

// WithTimeout adds the timeout to the deploy configuration params
func (o *DeployConfigurationParams) WithTimeout(timeout time.Duration) *DeployConfigurationParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the deploy configuration params
func (o *DeployConfigurationParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the deploy configuration params
func (o *DeployConfigurationParams) WithContext(ctx context.Context) *DeployConfigurationParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the deploy configuration params
func (o *DeployConfigurationParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the deploy configuration params
func (o *DeployConfigurationParams) WithHTTPClient(client *http.Client) *DeployConfigurationParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the deploy configuration params
func (o *DeployConfigurationParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the deploy configuration params
func (o *DeployConfigurationParams) WithBody(body *dto.DeployConfigurationRequest) *DeployConfigurationParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the deploy configuration params
func (o *DeployConfigurationParams) SetBody(body *dto.DeployConfigurationRequest) {
	o.Body = body
}

// WithEnvName adds the envName to the deploy configuration params
func (o *DeployConfigurationParams) WithEnvName(envName string) *DeployConfigurationParams {
	o.SetEnvName(envName)
	return o
}

// SetEnvName adds the envName to the deploy configuration params
func (o *DeployConfigurationParams) SetEnvName(envName string) {
	o.EnvName = envName
}

// WithOrgName adds the orgName to the deploy configuration params
func (o *DeployConfigurationParams) WithOrgName(orgName string) *DeployConfigurationParams {
	o.SetOrgName(orgName)
	return o
}

// SetOrgName adds the orgName to the deploy configuration params
func (o *DeployConfigurationParams) SetOrgName(orgName string) {
	o.OrgName = orgName
}

// WithPipelineName adds the pipelineName to the deploy configuration params
func (o *DeployConfigurationParams) WithPipelineName(pipelineName string) *DeployConfigurationParams {
	o.SetPipelineName(pipelineName)
	return o
}

// SetPipelineName adds the pipelineName to the deploy configuration params
func (o *DeployConfigurationParams) SetPipelineName(pipelineName string) {
	o.PipelineName = pipelineName
}

// WriteToRequest writes these params to a swagger request
func (o *DeployConfigurationParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

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
