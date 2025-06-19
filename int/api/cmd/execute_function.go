package cmd

import (
	"encoding/base64"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
)

func NewExecuteFunctionHandler(config *config.NetworkInfos) operations.CmdExecuteFunctionHandler {
	return &executeFunction{networkInfos: config}
}

type executeFunction struct {
	networkInfos *config.NetworkInfos
}

//nolint:funlen
func (e *executeFunction) Handle(params operations.CmdExecuteFunctionParams) middleware.Responder {
	fee, errResponse := amountToUint64(params.Body.Fee, uint64(sendOperation.DefaultFee))
	if errResponse != nil {
		return errResponse
	}

	maxGas := uint64(0)

	if string(params.Body.MaxGas) != "" {
		parsedMaxGas, err := strconv.ParseUint(string(params.Body.MaxGas), 10, 64)
		if err != nil {
			return operations.NewCmdExecuteFunctionBadRequest().WithPayload(
				&models.Error{
					Code:    errorInvalidMaxGas,
					Message: "Error during max gas conversion: " + err.Error(),
				})
		}

		maxGas = parsedMaxGas
	}

	expiry := uint64(sendOperation.DefaultExpiryInSlot)

	if params.Body.Expiry != nil {
		expiry = uint64(*params.Body.Expiry)
	}

	asyncReq := false

	if params.Body.Async != nil {
		asyncReq = *params.Body.Async
	}

	args, err := base64.StdEncoding.DecodeString(params.Body.Args)
	if err != nil {
		return operations.NewCmdExecuteFunctionUnprocessableEntity().
			WithPayload(
				&models.Error{
					Code:    errorInvalidArgs,
					Message: err.Error(),
				})
	}

	headers := signer.CustomHeader{
		Origin:  params.HTTPRequest.Header.Get("Origin"),
		Referer: params.HTTPRequest.Header.Get("Referer"),
	}

	operationWithEventResponse, err := onchain.CallFunction(
		e.networkInfos,
		params.Body.Nickname,
		params.Body.At,
		params.Body.Name,
		args,
		fee,
		maxGas,
		uint64(params.Body.Coins),
		expiry,
		asyncReq,
		signer.NewWalletPlugin(headers),
		params.Body.Description,
	)
	if err != nil {
		return operations.NewCmdExecuteFunctionInternalServerError().WithPayload(
			&models.Error{Code: errorSendOperation, Message: "Error: callSC failed: " + err.Error()})
	}

	return operations.NewCmdExecuteFunctionOK().WithPayload(
		&operations.CmdExecuteFunctionOKBody{
			FirstEvent:  &models.Events{Data: operationWithEventResponse.Event, Address: params.Body.At},
			OperationID: operationWithEventResponse.OperationResponse.OperationID,
		},
	)
}
