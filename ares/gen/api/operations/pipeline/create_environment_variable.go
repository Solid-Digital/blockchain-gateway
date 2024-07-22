// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// CreateEnvironmentVariableHandlerFunc turns a function with the right signature into a create environment variable handler
type CreateEnvironmentVariableHandlerFunc func(CreateEnvironmentVariableParams, *dto.User) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateEnvironmentVariableHandlerFunc) Handle(params CreateEnvironmentVariableParams, principal *dto.User) middleware.Responder {
	return fn(params, principal)
}

// CreateEnvironmentVariableHandler interface for that can handle valid create environment variable params
type CreateEnvironmentVariableHandler interface {
	Handle(CreateEnvironmentVariableParams, *dto.User) middleware.Responder
}

// NewCreateEnvironmentVariable creates a new http.Handler for the create environment variable operation
func NewCreateEnvironmentVariable(ctx *middleware.Context, handler CreateEnvironmentVariableHandler) *CreateEnvironmentVariable {
	return &CreateEnvironmentVariable{Context: ctx, Handler: handler}
}

/*CreateEnvironmentVariable swagger:route POST /orgs/{orgName}/pipelines/{pipelineName}/environments/{envName}/vars pipeline createEnvironmentVariable

Create a new variable for a pipeline environment

*/
type CreateEnvironmentVariable struct {
	Context *middleware.Context
	Handler CreateEnvironmentVariableHandler
}

func (o *CreateEnvironmentVariable) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateEnvironmentVariableParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *dto.User
	if uprinc != nil {
		principal = uprinc.(*dto.User) // this is really a dto.User, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
