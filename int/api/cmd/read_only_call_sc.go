package cmd

import (
	"encoding/base64"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/wallet"
)

func NewReadOnlyCallSCHandler(configManager *config.MSConfigManager) operations.CmdReadOnlyCallSCHandler {
	return &ReadOnlyCallSC{configManager: configManager}
}

type ReadOnlyCallSC struct {
	configManager *config.MSConfigManager
}

func (e *ReadOnlyCallSC) Handle(params operations.CmdReadOnlyCallSCParams) middleware.Responder {
	args, err := base64.StdEncoding.DecodeString(params.Body.Args)
	if err != nil {
		return operations.NewCmdReadOnlyCallSCUnprocessableEntity().
			WithPayload(
				&models.Error{
					Code:    errorInvalidArgs,
					Message: err.Error(),
				})
	}

	acc, err := wallet.Fetch(params.Body.Nickname)
	if err != nil {
		return operations.NewCmdReadOnlyCallSCBadRequest().WithPayload(
			&models.Error{
				Code:    errorInvalidNickname,
				Message: "Error during wallet fetch: " + err.Error(),
			})
	}

	fee := *params.Body.Fee

	if fee == "" {
		fee = "0"
	}

	coins, errResponse := amountToString(params.Body.Coins, uint64(0))
	if errResponse != nil {
		return errResponse
	}

	result, err := sendOperation.ReadOnlyCallSC(
		params.Body.At,
		params.Body.Name,
		args,
		coins,
		fee,
		acc.Address,
		node.NewClient(e.configManager.CurrentNetwork().NodeURL),
	)
	if err != nil {
		return operations.NewCmdReadOnlyCallSCInternalServerError().WithPayload(
			&models.Error{Code: errorSendOperation, Message: "Error: read only callSC failed: " + err.Error()})
	}

	modelResult := CreateReadOnlyResult(*result)

	return operations.NewCmdReadOnlyCallSCOK().WithPayload(&modelResult)
}
