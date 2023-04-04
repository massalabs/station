package cmd

import (
	"encoding/base64"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
)

func CreateExecuteFunctionHandler() func(params operations.CmdExecuteFunctionParams) middleware.Responder {
	return func(params operations.CmdExecuteFunctionParams) middleware.Responder {
		return ExecuteFunctionHandler(params)
	}
}

func ExecuteFunctionHandler(params operations.CmdExecuteFunctionParams) middleware.Responder {
	addr, err := base58.CheckDecode(params.Body.At[2:])
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().WithPayload(
			&models.Error{Code: errorCodeUnknownKeyID, Message: "Error : cannot decode Smart contract address : " + err.Error()})
	}

	addr = addr[1:]

	args, err := base64.StdEncoding.DecodeString(params.Body.Args)
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	c := node.NewDefaultClient()

	event, err := onchain.CallFunction(c, params.Body.Nickname, addr, params.Body.Name, args, uint64(params.Body.Coins))
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeSendOperation, Message: "Error : callSC failed " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(event)
}
