// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CmdExecuteFunctionHandlerFunc turns a function with the right signature into a cmd execute function handler
type CmdExecuteFunctionHandlerFunc func(CmdExecuteFunctionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CmdExecuteFunctionHandlerFunc) Handle(params CmdExecuteFunctionParams) middleware.Responder {
	return fn(params)
}

// CmdExecuteFunctionHandler interface for that can handle valid cmd execute function params
type CmdExecuteFunctionHandler interface {
	Handle(CmdExecuteFunctionParams) middleware.Responder
}

// NewCmdExecuteFunction creates a new http.Handler for the cmd execute function operation
func NewCmdExecuteFunction(ctx *middleware.Context, handler CmdExecuteFunctionHandler) *CmdExecuteFunction {
	return &CmdExecuteFunction{Context: ctx, Handler: handler}
}

/*
	CmdExecuteFunction swagger:route POST /cmd/executeFunction cmdExecuteFunction

CmdExecuteFunction cmd execute function API
*/
type CmdExecuteFunction struct {
	Context *middleware.Context
	Handler CmdExecuteFunctionHandler
}

func (o *CmdExecuteFunction) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCmdExecuteFunctionParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// CmdExecuteFunctionBody cmd execute function body
//
// swagger:model CmdExecuteFunctionBody
type CmdExecuteFunctionBody struct {

	// Arguments to pass to the function.
	Args string `json:"args,omitempty"`

	// Smart contract address exporting the function to call.
	// Required: true
	At *string `json:"at"`

	// coins
	Coins *CmdExecuteFunctionParamsBodyCoins `json:"coins,omitempty"`

	// Set the expiry duration (in number of slots) of the transaction.
	Expiry *int64 `json:"expiry,omitempty"`

	// Set the fee amount (in massa) that will be given to the block creator.
	Fee float64 `json:"fee,omitempty"`

	// gaz
	Gaz *CmdExecuteFunctionParamsBodyGaz `json:"gaz,omitempty"`

	// Defines the key to used to sign the transaction.
	KeyID *string `json:"keyId,omitempty"`

	// Function name to call.
	// Required: true
	Name *string `json:"name"`
}

func (o *CmdExecuteFunctionBody) UnmarshalJSON(b []byte) error {
	type CmdExecuteFunctionBodyAlias CmdExecuteFunctionBody
	var t CmdExecuteFunctionBodyAlias
	if err := json.Unmarshal([]byte("{\"args\":\"\",\"at\":\"A1MrqLgWq5XXDpTBH6fzXHUg7E8M5U2fYDAF3E1xnUSzyZuKpMh\",\"coins\":{\"parallel\":0,\"sequential\":0},\"expiry\":3,\"fee\":0,\"gaz\":{\"limit\":700000000,\"price\":0},\"keyId\":\"default\",\"name\":\"test\"}"), &t); err != nil {
		return err
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	*o = CmdExecuteFunctionBody(t)
	return nil
}

// Validate validates this cmd execute function body
func (o *CmdExecuteFunctionBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAt(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateCoins(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateGaz(formats); err != nil {
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

func (o *CmdExecuteFunctionBody) validateAt(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"at", "body", o.At); err != nil {
		return err
	}

	return nil
}

func (o *CmdExecuteFunctionBody) validateCoins(formats strfmt.Registry) error {
	if swag.IsZero(o.Coins) { // not required
		return nil
	}

	if o.Coins != nil {
		if err := o.Coins.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "coins")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("body" + "." + "coins")
			}
			return err
		}
	}

	return nil
}

func (o *CmdExecuteFunctionBody) validateGaz(formats strfmt.Registry) error {
	if swag.IsZero(o.Gaz) { // not required
		return nil
	}

	if o.Gaz != nil {
		if err := o.Gaz.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "gaz")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("body" + "." + "gaz")
			}
			return err
		}
	}

	return nil
}

func (o *CmdExecuteFunctionBody) validateName(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"name", "body", o.Name); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this cmd execute function body based on the context it is used
func (o *CmdExecuteFunctionBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateCoins(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateGaz(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CmdExecuteFunctionBody) contextValidateCoins(ctx context.Context, formats strfmt.Registry) error {

	if o.Coins != nil {
		if err := o.Coins.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "coins")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("body" + "." + "coins")
			}
			return err
		}
	}

	return nil
}

func (o *CmdExecuteFunctionBody) contextValidateGaz(ctx context.Context, formats strfmt.Registry) error {

	if o.Gaz != nil {
		if err := o.Gaz.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "gaz")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("body" + "." + "gaz")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *CmdExecuteFunctionBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CmdExecuteFunctionBody) UnmarshalBinary(b []byte) error {
	var res CmdExecuteFunctionBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CmdExecuteFunctionParamsBodyCoins Coins to be send from caller to smart contract address.
//
// swagger:model CmdExecuteFunctionParamsBodyCoins
type CmdExecuteFunctionParamsBodyCoins struct {

	// Number of parallel coins to transfer from the caller to the smart contract address.
	Parallel float64 `json:"parallel,omitempty"`

	// Number of sequential coins to transfer from the caller to the smart contract address.
	Sequential float64 `json:"sequential,omitempty"`
}

// Validate validates this cmd execute function params body coins
func (o *CmdExecuteFunctionParamsBodyCoins) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this cmd execute function params body coins based on context it is used
func (o *CmdExecuteFunctionParamsBodyCoins) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CmdExecuteFunctionParamsBodyCoins) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CmdExecuteFunctionParamsBodyCoins) UnmarshalBinary(b []byte) error {
	var res CmdExecuteFunctionParamsBodyCoins
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// CmdExecuteFunctionParamsBodyGaz Gaz attibutes. Gaz is a virtual resource consumed by node while running smart contract.
//
// swagger:model CmdExecuteFunctionParamsBodyGaz
type CmdExecuteFunctionParamsBodyGaz struct {

	// Maximum number of gaz unit that a node will be able consume.
	Limit *int64 `json:"limit,omitempty"`

	// Price of a gaz unit.
	Price float64 `json:"price,omitempty"`
}

// Validate validates this cmd execute function params body gaz
func (o *CmdExecuteFunctionParamsBodyGaz) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this cmd execute function params body gaz based on context it is used
func (o *CmdExecuteFunctionParamsBodyGaz) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *CmdExecuteFunctionParamsBodyGaz) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CmdExecuteFunctionParamsBodyGaz) UnmarshalBinary(b []byte) error {
	var res CmdExecuteFunctionParamsBodyGaz
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
