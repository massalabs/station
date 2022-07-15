package wallet

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

func NewDelete(walletStorage *sync.Map) operations.MgmtWalletDeleteHandler {
	return &walletDelete{walletStorage: walletStorage}
}

type walletDelete struct {
	walletStorage *sync.Map
}

func (c *walletDelete) Handle(params operations.MgmtWalletDeleteParams) middleware.Responder {
	if len(params.Nickname) == 0 {
		e := errorCodeWalletDeleteNoNickname
		msg := "Error: nickname field is mandatory."

		return operations.NewMgmtWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	_, ok := c.walletStorage.Load(params.Nickname)
	if !ok {
		e := errorCodeWalletNoWallet
		msg := "Error: a wallet with the same nickname already exists."

		return operations.NewMgmtWalletDeleteNotFound().WithPayload(
			&models.Error{
				Code:    &e,
				Message: &msg,
			})
	}

	c.walletStorage.Delete(params.Nickname)

	return operations.NewMgmtWalletDeleteNoContent()
}
