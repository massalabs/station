// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ThyraRegistryHandlerFunc turns a function with the right signature into a thyra registry handler
type ThyraRegistryHandlerFunc func(ThyraRegistryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ThyraRegistryHandlerFunc) Handle(params ThyraRegistryParams) middleware.Responder {
	return fn(params)
}

// ThyraRegistryHandler interface for that can handle valid thyra registry params
type ThyraRegistryHandler interface {
	Handle(ThyraRegistryParams) middleware.Responder
}

// NewThyraRegistry creates a new http.Handler for the thyra registry operation
func NewThyraRegistry(ctx *middleware.Context, handler ThyraRegistryHandler) *ThyraRegistry {
	return &ThyraRegistry{Context: ctx, Handler: handler}
}

/* ThyraRegistry swagger:route GET /thyra/registry/{resource} thyraRegistry

ThyraRegistry thyra registry API

*/
type ThyraRegistry struct {
	Context *middleware.Context
	Handler ThyraRegistryHandler
}

func (o *ThyraRegistry) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewThyraRegistryParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
