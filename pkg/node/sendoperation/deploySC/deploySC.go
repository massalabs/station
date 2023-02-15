package deploysc

import (
	"io"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/int/api/websites"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain"
)

func CreateCmdDeploySCHandler(
	app *fyne.App,
) func(params operations.CmdDeploySCParams) middleware.Responder {
	return func(params operations.CmdDeploySCParams) middleware.Responder {
		return cmdDeploySCHandler(params, app)
	}
}

func cmdDeploySCHandler(params operations.CmdDeploySCParams, app *fyne.App) middleware.Responder {
	client := node.NewDefaultClient()

	wallet, _ := websites.LoadAndUnprotectWallet(params.Nickname, app)

	file, err := io.ReadAll(params.Wasmfile)
	if err != nil {
		return operations.NewCmdDeploySCInternalServerError().
			WithPayload(
				&models.Error{
					Code:    err.Error(),
					Message: err.Error(),
				})
	}

	address, err := onchain.DeploySC(client, *wallet, file)
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
