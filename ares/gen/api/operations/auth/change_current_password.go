// Code generated by go-swagger; DO NOT EDIT.

package auth

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

// ChangeCurrentPasswordHandlerFunc turns a function with the right signature into a change current password handler
type ChangeCurrentPasswordHandlerFunc func(ChangeCurrentPasswordParams, *dto.User) middleware.Responder

// Handle executing the request and returning a response
func (fn ChangeCurrentPasswordHandlerFunc) Handle(params ChangeCurrentPasswordParams, principal *dto.User) middleware.Responder {
	return fn(params, principal)
}

// ChangeCurrentPasswordHandler interface for that can handle valid change current password params
type ChangeCurrentPasswordHandler interface {
	Handle(ChangeCurrentPasswordParams, *dto.User) middleware.Responder
}

// NewChangeCurrentPassword creates a new http.Handler for the change current password operation
func NewChangeCurrentPassword(ctx *middleware.Context, handler ChangeCurrentPasswordHandler) *ChangeCurrentPassword {
	return &ChangeCurrentPassword{Context: ctx, Handler: handler}
}

/*ChangeCurrentPassword swagger:route POST /auth/password/change auth changeCurrentPassword

Change the current user's password

*/
type ChangeCurrentPassword struct {
	Context *middleware.Context
	Handler ChangeCurrentPasswordHandler
}

func (o *ChangeCurrentPassword) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewChangeCurrentPasswordParams()

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
