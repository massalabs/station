package cmd

import (
	"encoding/base64"
	"io"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/wallet"
)

func NewReadOnlyExecuteSCHandler(configManager *config.MSConfigManager) operations.CmdReadOnlyExecuteSCHandler {
	return &ReadOnlyExecuteSC{configManager: configManager}
}

type ReadOnlyExecuteSC struct {
	configManager *config.MSConfigManager
}

func (e *ReadOnlyExecuteSC) Handle(params operations.CmdReadOnlyExecuteSCParams) middleware.Responder {
	file, err := io.ReadAll(params.Bytecode)
	if err != nil {
		return operations.NewCmdReadOnlyExecuteSCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	datastore, err := base64.StdEncoding.DecodeString(*params.Datastore)
	if err != nil {
		return operations.NewCmdReadOnlyExecuteSCUnprocessableEntity().
			WithPayload(
				&models.Error{
					Code:    errorInvalidArgs,
					Message: err.Error(),
				})
	}

	if len(datastore) == 0 {
		datastore = nil
	}

	acc, err := wallet.Fetch(params.Nickname)
	if err != nil {
		return operations.NewCmdReadOnlyExecuteSCBadRequest().WithPayload(
			&models.Error{
				Code:    errorInvalidNickname,
				Message: "Error during wallet fetch: " + err.Error(),
			})
	}

	coins, errResponse := amountToString(models.Amount(*params.Coins), uint64(0))
	if errResponse != nil {
		return errResponse
	}

	result, err := sendOperation.ReadOnlyExecuteSC(
		file,
		datastore,
		coins,
		*params.Fee,
		acc.Address,
		node.NewClient(e.configManager.CurrentNetwork().NodeURL),
	)
	if err != nil {
		return operations.NewCmdReadOnlyExecuteSCInternalServerError().WithPayload(
			&models.Error{Code: errorSendOperation, Message: "Error: read only callSC failed: " + err.Error()})
	}

	modelResult := CreateReadOnlyResult(*result)

	return operations.NewCmdReadOnlyExecuteSCOK().WithPayload(&modelResult)
}
