package deploysc

import (
	"io"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain"
)

func Handler(params operations.CmdDeploySCParams) middleware.Responder {
	client := node.NewDefaultClient()

	file, err := io.ReadAll(params.Wasmfile)
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
	address, err := onchain.DeploySCV2(client,
		params.Nickname,
		*params.GazLimit,
		*params.Coins,
		*params.Fee,
		*params.Expiry,
		file)
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
