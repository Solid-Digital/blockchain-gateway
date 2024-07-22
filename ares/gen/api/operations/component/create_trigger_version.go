// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// CreateTriggerVersionHandlerFunc turns a function with the right signature into a create trigger version handler
type CreateTriggerVersionHandlerFunc func(CreateTriggerVersionParams, *dto.User) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateTriggerVersionHandlerFunc) Handle(params CreateTriggerVersionParams, principal *dto.User) middleware.Responder {
	return fn(params, principal)
}

// CreateTriggerVersionHandler interface for that can handle valid create trigger version params
type CreateTriggerVersionHandler interface {
	Handle(CreateTriggerVersionParams, *dto.User) middleware.Responder
}

// NewCreateTriggerVersion creates a new http.Handler for the create trigger version operation
func NewCreateTriggerVersion(ctx *middleware.Context, handler CreateTriggerVersionHandler) *CreateTriggerVersion {
	return &CreateTriggerVersion{Context: ctx, Handler: handler}
}

/*CreateTriggerVersion swagger:route POST /orgs/{orgName}/triggers/{name}/versions component createTriggerVersion

Create trigger version

If the trigger that this version belongs to doesn't already exist, it will be created

*/
type CreateTriggerVersion struct {
	Context *middleware.Context
	Handler CreateTriggerVersionHandler
}

func (o *CreateTriggerVersion) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateTriggerVersionParams()

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
