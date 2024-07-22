// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// CreateBaseHandlerFunc turns a function with the right signature into a create base handler
type CreateBaseHandlerFunc func(CreateBaseParams, *dto.User) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateBaseHandlerFunc) Handle(params CreateBaseParams, principal *dto.User) middleware.Responder {
	return fn(params, principal)
}

// CreateBaseHandler interface for that can handle valid create base params
type CreateBaseHandler interface {
	Handle(CreateBaseParams, *dto.User) middleware.Responder
}

// NewCreateBase creates a new http.Handler for the create base operation
func NewCreateBase(ctx *middleware.Context, handler CreateBaseHandler) *CreateBase {
	return &CreateBase{Context: ctx, Handler: handler}
}

/*CreateBase swagger:route POST /orgs/{orgName}/bases component createBase

Create new base

*/
type CreateBase struct {
	Context *middleware.Context
	Handler CreateBaseHandler
}

func (o *CreateBase) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateBaseParams()

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
