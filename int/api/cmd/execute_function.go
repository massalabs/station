package cmd

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	callSC "github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/wallet"
)

//nolint:nolintlint,ireturn
func ExecuteFunctionHandler(params operations.CmdExecuteFunctionParams) middleware.Responder {
	addr, err := base58.CheckDecode(params.Body.At[1:])
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().WithPayload(
			&models.Error{Code: errorCodeUnknownKeyID, Message: "Error : cannot decode Smart contract address : " + err.Error()})
	}

	addr = addr[1:]

	wallet, err := wallet.Load(params.Body.Nickname)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeGetWallet,
				Message: "Error cannot load wallet : " + err.Error(),
			})
	}

	err = wallet.Unprotect(params.HTTPRequest.Header.Get("Authorization"), 0)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeWalletWrongPassword, Message: "Error : cannot uncipher the wallet : " + err.Error()})
	}

	callSC := callSC.New(
		addr,
		params.Body.Name,
		[]byte(params.Body.Args),
		uint64(params.Body.Gaz.Price),
		uint64(*params.Body.Gaz.Limit),
		uint64(params.Body.Coins.Sequential),
		uint64(params.Body.Coins.Parallel))

	c := node.NewDefaultClient()

	operationID, err := sendOperation.Call(
		c,
		sendOperation.DefaultSlotsDuration,
		uint64(params.Body.Fee),
		callSC,
		wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeSendOperation, Message: "Error : callSC operations not sent " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(operationID)
}
