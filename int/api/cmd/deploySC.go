package cmd

import (
	"encoding/base64"
	"io"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
)

func NewDeploySCHandler(config *config.AppConfig) operations.CmdDeploySCHandler {
	return &deploySC{config: config}
}

type deploySC struct {
	config *config.AppConfig
}

func (d *deploySC) Handle(params operations.CmdDeploySCParams) middleware.Responder {
	client := node.NewClient(d.config.NodeURL)

	file, err := io.ReadAll(params.SmartContract)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}
	/* All the pointers below cannot be null as the swagger hydrate
	each one with their default value defined in swagger.yml,
	if no values are provided for these parameters.
	*/
	decodedDatastore, err := base64.StdEncoding.DecodeString(*params.Datastore)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	if len(decodedDatastore) == 0 {
		decodedDatastore = nil
	}

	operationResponse, events, err := onchain.DeploySC(client,
		params.WalletNickname,
		*params.GasLimit,
		*params.Coins,
		*params.Fee,
		*params.Expiry,
		file,
		decodedDatastore,
		sendoperation.OperationBatch{NewBatch: false, CorrelationID: ""},
		&signer.WalletPlugin{},
	)
	if err != nil {
		return operations.NewCmdDeploySCInternalServerError().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	scAddress, _ := onchain.FindDeployedAddress(events)

	return operations.NewCmdDeploySCOK().
		WithPayload(&operations.CmdDeploySCOKBody{
			OperationID: operationResponse.OperationID,
			FirstEvent: &models.Events{
				Data:    events[0].Data,
				Address: scAddress,
			},
		})
}
