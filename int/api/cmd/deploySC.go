package cmd

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"io"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
)

//go:embed sc/deployer.wasm
var deployerSCByteCode []byte

func NewDeploySCHandler(config *config.NetworkInfos) operations.CmdDeploySCHandler {
	return &deploySC{networkInfos: config}
}

type deploySC struct {
	networkInfos *config.NetworkInfos
}

func (d *deploySC) Handle(params operations.CmdDeploySCParams) middleware.Responder {

	_smartContractBytes, err := base64.StdEncoding.DecodeString(params.Body.SmartContract)
	smartContractReader := bytes.NewReader(_smartContractBytes)
	smartContractByteCode, err := io.ReadAll(smartContractReader)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	_parameters, err := base64.StdEncoding.DecodeString(params.Body.Parameters)
	parameterReader := bytes.NewReader(_parameters)
	parameters, err := io.ReadAll(parameterReader)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	operationResponse, events, err := onchain.DeploySC(
		d.networkInfos,
		params.Body.Nickname,
		sendoperation.MaxGasAllowedExecuteSC, // default
		*params.Body.MaxCoins, // maxCoins 
		*params.Body.Coins,    // smart contract deployment "fee"
		sendoperation.DefaultExpiryInSlot,
		parameters, 
		smartContractByteCode,
		deployerSCByteCode,
		&signer.WalletPlugin{},
		"Deploying website",
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
