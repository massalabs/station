package wallet

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewDelete(walletStorage *sync.Map) operations.MgmtWalletDeleteHandler {
	return &walletDelete{walletStorage: walletStorage}
}

type walletDelete struct {
	walletStorage *sync.Map
}

func (c *walletDelete) Handle(params operations.MgmtWalletDeleteParams) middleware.Responder {
	if len(params.Nickname) == 0 {
		return operations.NewMgmtWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodeWalletDeleteNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	c.walletStorage.Delete(params.Nickname)

	err := wallet.Delete(params.Nickname)
	if err != nil {
		panic(err)
	}

	return operations.NewMgmtWalletDeleteNoContent()
}
