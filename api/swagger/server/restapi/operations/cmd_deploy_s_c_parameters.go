// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CmdDeploySCMaxParseMemory sets the maximum size in bytes for
// the multipart form parser for this operation.
//
// The default value is 32 MB.
// The multipart parser stores up to this + 10MB.
var CmdDeploySCMaxParseMemory int64 = 32 << 20

// NewCmdDeploySCParams creates a new CmdDeploySCParams object
// with the default values initialized.
func NewCmdDeploySCParams() CmdDeploySCParams {

	var (
		// initialize parameters with default values

		coinsDefault     = uint64(0)
		datastoreDefault = string("")
		expiryDefault    = uint64(2)
		feeDefault       = uint64(0)
		gazLimitDefault  = uint64(7e+08)
		gazPriceDefault  = uint64(0)
	)

	return CmdDeploySCParams{
		Coins: &coinsDefault,

		Datastore: &datastoreDefault,

		Expiry: &expiryDefault,

		Fee: &feeDefault,

		GazLimit: &gazLimitDefault,

		GazPrice: &gazPriceDefault,
	}
}

// CmdDeploySCParams contains all the bound params for the cmd deploy s c operation
// typically these are obtained from a http.Request
//
// swagger:parameters cmdDeploySC
type CmdDeploySCParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Set the number of coins that will be sent along the deployment call.
	  Minimum: 0
	  In: formData
	  Default: 0
	*/
	Coins *uint64
	/*Datastore that will be sent along the smart contract.
	  In: formData
	  Default: ""
	*/
	Datastore *string
	/*Set the expiry duration (in number of slots) of the transaction.
	  Minimum: 0
	  In: formData
	  Default: 2
	*/
	Expiry *uint64
	/*Set the fee amount (in massa) that will be given to the block creator.
	  Minimum: 0
	  In: formData
	  Default: 0
	*/
	Fee *uint64
	/*Maximum number of gaz unit that a node will be able to consume.
	  Minimum: 0
	  In: formData
	  Default: 7e+08
	*/
	GazLimit *uint64
	/*Price of a gaz unit.
	  Minimum: 0
	  In: formData
	  Default: 0
	*/
	GazPrice *uint64
	/*Smart contract file in a Wasm format.
	  Required: true
	  In: formData
	*/
	SmartContract io.ReadCloser
	/*Name of the wallet used to deploy the smart contract.
	  Required: true
	  In: formData
	*/
	WalletNickname string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCmdDeploySCParams() beforehand.
func (o *CmdDeploySCParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(CmdDeploySCMaxParseMemory); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}
	fds := runtime.Values(r.Form)

	fdCoins, fdhkCoins, _ := fds.GetOK("coins")
	if err := o.bindCoins(fdCoins, fdhkCoins, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDatastore, fdhkDatastore, _ := fds.GetOK("datastore")
	if err := o.bindDatastore(fdDatastore, fdhkDatastore, route.Formats); err != nil {
		res = append(res, err)
	}

	fdExpiry, fdhkExpiry, _ := fds.GetOK("expiry")
	if err := o.bindExpiry(fdExpiry, fdhkExpiry, route.Formats); err != nil {
		res = append(res, err)
	}

	fdFee, fdhkFee, _ := fds.GetOK("fee")
	if err := o.bindFee(fdFee, fdhkFee, route.Formats); err != nil {
		res = append(res, err)
	}

	fdGazLimit, fdhkGazLimit, _ := fds.GetOK("gazLimit")
	if err := o.bindGazLimit(fdGazLimit, fdhkGazLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	fdGazPrice, fdhkGazPrice, _ := fds.GetOK("gazPrice")
	if err := o.bindGazPrice(fdGazPrice, fdhkGazPrice, route.Formats); err != nil {
		res = append(res, err)
	}

	smartContract, smartContractHeader, err := r.FormFile("smartContract")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "smartContract", err))
	} else if err := o.bindSmartContract(smartContract, smartContractHeader); err != nil {
		// Required: true
		res = append(res, err)
	} else {
		o.SmartContract = &runtime.File{Data: smartContract, Header: smartContractHeader}
	}

	fdWalletNickname, fdhkWalletNickname, _ := fds.GetOK("walletNickname")
	if err := o.bindWalletNickname(fdWalletNickname, fdhkWalletNickname, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCoins binds and validates parameter Coins from formData.
func (o *CmdDeploySCParams) bindCoins(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCmdDeploySCParams()
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("coins", "formData", "uint64", raw)
	}
	o.Coins = &value

	if err := o.validateCoins(formats); err != nil {
		return err
	}

	return nil
}

// validateCoins carries on validations for parameter Coins
func (o *CmdDeploySCParams) validateCoins(formats strfmt.Registry) error {

	if err := validate.MinimumUint("coins", "formData", *o.Coins, 0, false); err != nil {
		return err
	}

	return nil
}

// bindDatastore binds and validates parameter Datastore from formData.
func (o *CmdDeploySCParams) bindDatastore(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCmdDeploySCParams()
		return nil
	}
	o.Datastore = &raw

	return nil
}

// bindExpiry binds and validates parameter Expiry from formData.
func (o *CmdDeploySCParams) bindExpiry(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCmdDeploySCParams()
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("expiry", "formData", "uint64", raw)
	}
	o.Expiry = &value

	if err := o.validateExpiry(formats); err != nil {
		return err
	}

	return nil
}

// validateExpiry carries on validations for parameter Expiry
func (o *CmdDeploySCParams) validateExpiry(formats strfmt.Registry) error {

	if err := validate.MinimumUint("expiry", "formData", *o.Expiry, 0, false); err != nil {
		return err
	}

	return nil
}

// bindFee binds and validates parameter Fee from formData.
func (o *CmdDeploySCParams) bindFee(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCmdDeploySCParams()
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("fee", "formData", "uint64", raw)
	}
	o.Fee = &value

	if err := o.validateFee(formats); err != nil {
		return err
	}

	return nil
}

// validateFee carries on validations for parameter Fee
func (o *CmdDeploySCParams) validateFee(formats strfmt.Registry) error {

	if err := validate.MinimumUint("fee", "formData", *o.Fee, 0, false); err != nil {
		return err
	}

	return nil
}

// bindGazLimit binds and validates parameter GazLimit from formData.
func (o *CmdDeploySCParams) bindGazLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCmdDeploySCParams()
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("gazLimit", "formData", "uint64", raw)
	}
	o.GazLimit = &value

	if err := o.validateGazLimit(formats); err != nil {
		return err
	}

	return nil
}

// validateGazLimit carries on validations for parameter GazLimit
func (o *CmdDeploySCParams) validateGazLimit(formats strfmt.Registry) error {

	if err := validate.MinimumUint("gazLimit", "formData", *o.GazLimit, 0, false); err != nil {
		return err
	}

	return nil
}

// bindGazPrice binds and validates parameter GazPrice from formData.
func (o *CmdDeploySCParams) bindGazPrice(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCmdDeploySCParams()
		return nil
	}

	value, err := swag.ConvertUint64(raw)
	if err != nil {
		return errors.InvalidType("gazPrice", "formData", "uint64", raw)
	}
	o.GazPrice = &value

	if err := o.validateGazPrice(formats); err != nil {
		return err
	}

	return nil
}

// validateGazPrice carries on validations for parameter GazPrice
func (o *CmdDeploySCParams) validateGazPrice(formats strfmt.Registry) error {

	if err := validate.MinimumUint("gazPrice", "formData", *o.GazPrice, 0, false); err != nil {
		return err
	}

	return nil
}

// bindSmartContract binds file parameter SmartContract.
//
// The only supported validations on files are MinLength and MaxLength
func (o *CmdDeploySCParams) bindSmartContract(file multipart.File, header *multipart.FileHeader) error {
	return nil
}

// bindWalletNickname binds and validates parameter WalletNickname from formData.
func (o *CmdDeploySCParams) bindWalletNickname(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("walletNickname", "formData", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("walletNickname", "formData", raw); err != nil {
		return err
	}
	o.WalletNickname = raw

	return nil
}
