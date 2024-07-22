// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// GetCurrentUserHandlerFunc turns a function with the right signature into a get current user handler
type GetCurrentUserHandlerFunc func(GetCurrentUserParams, *dto.User) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCurrentUserHandlerFunc) Handle(params GetCurrentUserParams, principal *dto.User) middleware.Responder {
	return fn(params, principal)
}

// GetCurrentUserHandler interface for that can handle valid get current user params
type GetCurrentUserHandler interface {
	Handle(GetCurrentUserParams, *dto.User) middleware.Responder
}

// NewGetCurrentUser creates a new http.Handler for the get current user operation
func NewGetCurrentUser(ctx *middleware.Context, handler GetCurrentUserHandler) *GetCurrentUser {
	return &GetCurrentUser{Context: ctx, Handler: handler}
}

/*GetCurrentUser swagger:route GET /auth/user auth getCurrentUser

Get Current User Profile

Get the current user's profile.

*/
type GetCurrentUser struct {
	Context *middleware.Context
	Handler GetCurrentUserHandler
}

func (o *GetCurrentUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetCurrentUserParams()

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
