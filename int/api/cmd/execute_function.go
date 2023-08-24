package cmd

import (
	"encoding/base64"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
)

func NewExecuteFunctionHandler(config *config.AppConfig) operations.CmdExecuteFunctionHandler {
	return &executeFunction{config: config}
}

type executeFunction struct {
	config *config.AppConfig
}

func (e *executeFunction) Handle(params operations.CmdExecuteFunctionParams) middleware.Responder {
	args, err := base64.StdEncoding.DecodeString(params.Body.Args)
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().
			WithPayload(
				&models.Error{
					Code:    errorCodeInvalidArgs,
					Message: err.Error(),
				})
	}

	c := node.NewClient(e.config.NodeURL)

	operationWithEventResponse, err := onchain.CallFunction(
		c,
		params.Body.Nickname,
		params.Body.At,
		params.Body.Name,
		args,
		uint64(params.Body.Coins),
		uint64(*params.Body.Expiry),
		params.Body.Async,
		sendOperation.OperationBatch{NewBatch: false, CorrelationID: ""},
		&signer.WalletPlugin{},
	)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeSendOperation, Message: "Error: callSC failed: " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(
		&operations.CmdExecuteFunctionOKBody{
			FirstEvent:  &models.Events{Data: operationWithEventResponse.Event, Address: params.Body.At},
			OperationID: operationWithEventResponse.OperationResponse.OperationID,
		},
	)
}
