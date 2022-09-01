// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// MgmtWalletGetterHandlerFunc turns a function with the right signature into a mgmt wallet getter handler
type MgmtWalletGetterHandlerFunc func(MgmtWalletGetterParams) middleware.Responder

// Handle executing the request and returning a response
func (fn MgmtWalletGetterHandlerFunc) Handle(params MgmtWalletGetterParams) middleware.Responder {
	return fn(params)
}

// MgmtWalletGetterHandler interface for that can handle valid mgmt wallet getter params
type MgmtWalletGetterHandler interface {
	Handle(MgmtWalletGetterParams) middleware.Responder
}

// NewMgmtWalletGetter creates a new http.Handler for the mgmt wallet getter operation
func NewMgmtWalletGetter(ctx *middleware.Context, handler MgmtWalletGetterHandler) *MgmtWalletGetter {
	return &MgmtWalletGetter{Context: ctx, Handler: handler}
}

/* MgmtWalletGetter swagger:route GET /mgmt/wallet/{nickname} mgmtWalletGetter

MgmtWalletGetter mgmt wallet getter API

*/
type MgmtWalletGetter struct {
	Context *middleware.Context
	Handler MgmtWalletGetterHandler
}

func (o *MgmtWalletGetter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewMgmtWalletGetterParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
