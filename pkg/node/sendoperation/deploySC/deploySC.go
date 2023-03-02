package deploysc

import (
	"io"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain"
)

func CreateCmdDeploySCHandler() func(params operations.CmdDeploySCParams) middleware.Responder {
	//nolint:gocritic
	return func(params operations.CmdDeploySCParams) middleware.Responder {
		return cmdDeploySCHandler(params)
	}
}

func cmdDeploySCHandler(params operations.CmdDeploySCParams) middleware.Responder {
	client := node.NewDefaultClient()

	file, err := io.ReadAll(params.Wasmfile)
	if err != nil {
		return operations.NewCmdDeploySCInternalServerError().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	address, err := onchain.DeploySCV2(client, params.Nickname, file)
	if err != nil {
		return operations.NewCmdDeploySCInternalServerError().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	return operations.NewCmdDeploySCOK().
		WithPayload(address)
}
