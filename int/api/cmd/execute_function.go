package cmd

import (
	"encoding/base64"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/signer"
	"github.com/massalabs/thyra/pkg/onchain"
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
					Code:    err.Error(),
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
		sendOperation.OperationBatch{NewBatch: false, CorrelationID: ""},
		&signer.WalletPlugin{},
	)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorCodeSendOperation, Message: "Error: callSC failed: " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(operationWithEventResponse.Event)
}
