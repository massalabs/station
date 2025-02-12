package cmd

import (
	_ "embed"
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

//go:embed sc/deployer.wasm
var deployerSCByteCode []byte

func NewDeploySCHandler(config *config.NetworkInfos) operations.CmdDeploySCHandler {
	return &deploySC{networkInfos: config}
}

type deploySC struct {
	networkInfos *config.NetworkInfos
}

//nolint:funlen
func (d *deploySC) Handle(params operations.CmdDeploySCParams) middleware.Responder {
	if params.Body.SmartContract == "" {
		return operations.NewCmdDeploySCUnprocessableEntity().
			WithPayload(
				&models.Error{
					Message: "Smart contract bytecode is required",
				})
	}

	smartContractByteCode, err := base64.StdEncoding.DecodeString(params.Body.SmartContract)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	parameters, err := base64.StdEncoding.DecodeString(params.Body.Parameters)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    errorInvalidArgs,
					Message: err.Error(),
				})
	}

	maxCoins, err := strconv.ParseUint(*params.Body.MaxCoins, 10, 64)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    errorInvalidMaxCoins,
					Message: err.Error(),
				})
	}

	coins := uint64(0)
	if params.Body.Coins != nil {
		coins, err = strconv.ParseUint(*params.Body.Coins, 10, 64)
		if err != nil {
			return operations.NewCmdDeploySCBadRequest().
				WithPayload(
					&models.Error{
						Code:    errorInvalidCoin,
						Message: err.Error(),
					})
		}
	}

	fee, err := strconv.ParseUint(*params.Body.Fee, 10, 64)
	if err != nil {
		return operations.NewCmdDeploySCBadRequest().
			WithPayload(
				&models.Error{
					Code:    errorInvalidFee,
					Message: err.Error(),
				})
	}

	operationResponse, events, err := onchain.DeploySC(
		d.networkInfos,
		params.Body.Nickname,
		sendoperation.MaxGasAllowedExecuteSC, // default
		maxCoins,                             // maxCoins
		coins,                                // Coins to send for storage
		fee,                                  // operation fee
		sendoperation.DefaultExpiryInSlot,
		parameters,
		smartContractByteCode,
		deployerSCByteCode,
		&signer.WalletPlugin{},
		"Deploying contract: "+params.Body.Description,
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
