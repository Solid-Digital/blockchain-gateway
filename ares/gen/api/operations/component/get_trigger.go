// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// GetTriggerHandlerFunc turns a function with the right signature into a get trigger handler
type GetTriggerHandlerFunc func(GetTriggerParams, *dto.User) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTriggerHandlerFunc) Handle(params GetTriggerParams, principal *dto.User) middleware.Responder {
	return fn(params, principal)
}

// GetTriggerHandler interface for that can handle valid get trigger params
type GetTriggerHandler interface {
	Handle(GetTriggerParams, *dto.User) middleware.Responder
}

// NewGetTrigger creates a new http.Handler for the get trigger operation
func NewGetTrigger(ctx *middleware.Context, handler GetTriggerHandler) *GetTrigger {
	return &GetTrigger{Context: ctx, Handler: handler}
}

/*GetTrigger swagger:route GET /orgs/{orgName}/triggers/{name} component getTrigger

Get Trigger Detail

Get a trigger component with all its versions

*/
type GetTrigger struct {
	Context *middleware.Context
	Handler GetTriggerHandler
}

func (o *GetTrigger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTriggerParams()

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
