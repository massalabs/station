package genericInfo

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
)

func CreateMgmtWalletBalanceHandler() func(params operations.MgmtWalletBalanceParams) middleware.Responder {
	return func(params operations.MgmtWalletBalanceParams) middleware.Responder {
		return MgmtWalletBalanceHandler(params)
	}
}

func MgmtWalletBalanceHandler(params operations.MgmtWalletBalanceParams) middleware.Responder {
	client := node.NewDefaultClient()

	balance, err := node.FetchBalance(client, params.Address)
	if err != nil {
		return operations.NewMgmtPluginsListInternalServerError().WithPayload(
			&models.Error{Code: "FetchingBalance", Message: "Error : Cannot get Balance from Address: " + err.Error()})
	}

	return operations.NewMgmtWalletBalanceOK().WithPayload(&models.Data{Data: balance.Candidate.String()})
}
