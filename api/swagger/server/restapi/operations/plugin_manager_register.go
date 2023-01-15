// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PluginManagerRegisterHandlerFunc turns a function with the right signature into a plugin manager register handler
type PluginManagerRegisterHandlerFunc func(PluginManagerRegisterParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PluginManagerRegisterHandlerFunc) Handle(params PluginManagerRegisterParams) middleware.Responder {
	return fn(params)
}

// PluginManagerRegisterHandler interface for that can handle valid plugin manager register params
type PluginManagerRegisterHandler interface {
	Handle(PluginManagerRegisterParams) middleware.Responder
}

// NewPluginManagerRegister creates a new http.Handler for the plugin manager register operation
func NewPluginManagerRegister(ctx *middleware.Context, handler PluginManagerRegisterHandler) *PluginManagerRegister {
	return &PluginManagerRegister{Context: ctx, Handler: handler}
}

/*
	PluginManagerRegister swagger:route POST /plugin-manager/register pluginManagerRegister

PluginManagerRegister plugin manager register API
*/
type PluginManagerRegister struct {
	Context *middleware.Context
	Handler PluginManagerRegisterHandler
}

func (o *PluginManagerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPluginManagerRegisterParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PluginManagerRegisterBody plugin manager register body
//
// swagger:model PluginManagerRegisterBody
type PluginManagerRegisterBody struct {

	// Plugin API specification
	APISpec string `json:"api_spec,omitempty"`

	// URL authority to use to connect to the plugin
	// Required: true
	Authority string `json:"authority"`

	// Plugin description.
	// Required: true
	Description string `json:"description"`

	// Plugin identifier.
	// Required: true
	ID int64 `json:"id"`

	// Plugin logo.
	// Required: true
	// Format: byte
	Logo strfmt.Base64 `json:"logo"`

	// Plugin name.
	// Required: true
	Name string `json:"name"`
}

// Validate validates this plugin manager register body
func (o *PluginManagerRegisterBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAuthority(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateLogo(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PluginManagerRegisterBody) validateAuthority(formats strfmt.Registry) error {

	if err := validate.RequiredString("body"+"."+"authority", "body", o.Authority); err != nil {
		return err
	}

	return nil
}

func (o *PluginManagerRegisterBody) validateDescription(formats strfmt.Registry) error {

	if err := validate.RequiredString("body"+"."+"description", "body", o.Description); err != nil {
		return err
	}

	return nil
}

func (o *PluginManagerRegisterBody) validateID(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"id", "body", int64(o.ID)); err != nil {
		return err
	}

	return nil
}

func (o *PluginManagerRegisterBody) validateLogo(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"logo", "body", strfmt.Base64(o.Logo)); err != nil {
		return err
	}

	return nil
}

func (o *PluginManagerRegisterBody) validateName(formats strfmt.Registry) error {

	if err := validate.RequiredString("body"+"."+"name", "body", o.Name); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this plugin manager register body based on context it is used
func (o *PluginManagerRegisterBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PluginManagerRegisterBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PluginManagerRegisterBody) UnmarshalBinary(b []byte) error {
	var res PluginManagerRegisterBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
