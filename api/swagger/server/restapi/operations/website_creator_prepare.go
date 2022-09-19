// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// WebsiteCreatorPrepareHandlerFunc turns a function with the right signature into a website creator prepare handler
type WebsiteCreatorPrepareHandlerFunc func(WebsiteCreatorPrepareParams) middleware.Responder

// Handle executing the request and returning a response
func (fn WebsiteCreatorPrepareHandlerFunc) Handle(params WebsiteCreatorPrepareParams) middleware.Responder {
	return fn(params)
}

// WebsiteCreatorPrepareHandler interface for that can handle valid website creator prepare params
type WebsiteCreatorPrepareHandler interface {
	Handle(WebsiteCreatorPrepareParams) middleware.Responder
}

// NewWebsiteCreatorPrepare creates a new http.Handler for the website creator prepare operation
func NewWebsiteCreatorPrepare(ctx *middleware.Context, handler WebsiteCreatorPrepareHandler) *WebsiteCreatorPrepare {
	return &WebsiteCreatorPrepare{Context: ctx, Handler: handler}
}

/*
	WebsiteCreatorPrepare swagger:route PUT /websiteCreator/prepare websiteCreatorPrepare

WebsiteCreatorPrepare website creator prepare API
*/
type WebsiteCreatorPrepare struct {
	Context *middleware.Context
	Handler WebsiteCreatorPrepareHandler
}

func (o *WebsiteCreatorPrepare) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewWebsiteCreatorPrepareParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
