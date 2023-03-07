package cmd

import (
	"encoding/base64"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
)

func CreateExecuteFunctionHandler(app *fyne.App) func(params operations.CmdExecuteFunctionParams) middleware.Responder {
	return func(params operations.CmdExecuteFunctionParams) middleware.Responder {
		return ExecuteFunctionHandler(params, app)
	}
}

func ExecuteFunctionHandler(params operations.CmdExecuteFunctionParams, app *fyne.App) middleware.Responder {
	addr, err := base58.CheckDecode(params.Body.At[2:])
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().WithPayload(
			&models.Error{Code: errorCodeUnknownKeyID, Message: "Error : cannot decode Smart contract address : " + err.Error()})
	}

	addr = addr[2:]

	args, err := base64.StdEncoding.DecodeString(params.Body.Args)
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	// callSC := callSC.New(
	// 	addr,
	// 	params.Body.Name,
	// 	args,
	// 	uint64(*params.Body.Gaz.Limit),
	// 	uint64(params.Body.Coins))

	c := node.NewDefaultClient()

	event, err := onchain.CallFunctionV2(c, params.Body.Nickname, addr, params.Body.Name, args, uint64(params.Body.Coins))

	// operationID, err := onchain.CallFunctionV2(
	// 	c,
	// 	sendOperation.DefaultSlotsDuration,
	// 	uint64(params.Body.Fee),
	// 	callSC,
	// 	wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeSendOperation, Message: "Error : callSC failed " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(event)
}

func createInternalServerError(errorCode string, errorMessage string) middleware.Responder {
	return operations.NewCmdExecuteFunctionInternalServerError().
		WithPayload(
			&models.Error{
				Code:    errorCode,
				Message: errorMessage,
			})
}
