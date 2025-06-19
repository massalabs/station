package cmd

import (
	"encoding/base64"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
)

func NewExecuteSCHandler(config *config.NetworkInfos) operations.CmdExecuteSCHandler {
	return &executeSC{networkInfos: config}
}

type executeSC struct {
	networkInfos *config.NetworkInfos
}

//nolint:funlen
func (d *executeSC) Handle(params operations.CmdExecuteSCParams) middleware.Responder {
	if params.Body.Bytecode == "" {
		return operations.NewCmdExecuteSCUnprocessableEntity().
			WithPayload(
				&models.Error{
					Message: "Smart contract bytecode is required",
				})
	}

	smartContractByteCode, err := base64.StdEncoding.DecodeString(params.Body.Bytecode)
	if err != nil {
		return operations.NewCmdExecuteSCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	maxCoins, err := strconv.ParseUint(string(params.Body.MaxCoins), 10, 64)
	if err != nil {
		return operations.NewCmdExecuteSCBadRequest().
			WithPayload(
				&models.Error{
					Code:    errorInvalidMaxCoins,
					Message: err.Error(),
				})
	}

	fee, err := strconv.ParseUint(string(params.Body.Fee), 10, 64)
	if err != nil {
		return operations.NewCmdExecuteSCBadRequest().
			WithPayload(
				&models.Error{
					Code:    errorInvalidFee,
					Message: err.Error(),
				})
	}

	maxGas := uint64(sendoperation.MaxGasAllowedExecuteSC)

	if string(params.Body.MaxGas) != "" {
		parsedMaxGas, err := strconv.ParseUint(string(params.Body.MaxGas), 10, 64)
		if err != nil {
			return operations.NewCmdExecuteSCBadRequest().WithPayload(
				&models.Error{
					Code:    errorInvalidMaxGas,
					Message: "Error during max gas conversion: " + err.Error(),
				})
		}

		maxGas = parsedMaxGas
	}

	var datastore []onchain.DatastoreEntry = []onchain.DatastoreEntry{}

	if len(params.Body.Datastore) > 0 {
		datastore = make([]onchain.DatastoreEntry, len(params.Body.Datastore))
		for i, entry := range params.Body.Datastore {
			key, err := base64.StdEncoding.DecodeString(string(entry[0]))
			if err != nil {
				return operations.NewCmdExecuteSCBadRequest().
					WithPayload(
						&models.Error{
							Code:    errorInvalidDatastore,
							Message: err.Error(),
						})
			}
			value, err := base64.StdEncoding.DecodeString(string(entry[1]))
			if err != nil {
				return operations.NewCmdExecuteSCBadRequest().
					WithPayload(
						&models.Error{
							Code:    errorInvalidDatastore,
							Message: err.Error(),
						})
			}
			datastore[i] = onchain.DatastoreEntry{
				Key:   key,
				Value: value,
			}
		}
	}

	headers := signer.CustomHeader{
		Origin: params.HTTPRequest.Header.Get("Origin"),
		Referer: params.HTTPRequest.Header.Get("Referer"),
	}

	operationResponse, err := onchain.ExecuteSC(
		d.networkInfos,
		params.Body.Nickname,
		maxGas,
		maxCoins,
		fee,
		sendoperation.DefaultExpiryInSlot,
		smartContractByteCode,
		datastore,
		signer.NewWalletPlugin(headers),
		"Executing contract bytecode: "+params.Body.Description,
	)
	if err != nil {
		return operations.NewCmdExecuteSCInternalServerError().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	return operations.NewCmdExecuteSCOK().
		WithPayload(&operations.CmdExecuteSCOKBody{
			OperationID: operationResponse.OperationID,
		})
}
