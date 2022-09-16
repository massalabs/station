package wallet

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/wallet"
)

//nolint:nolintlint,ireturn
func NewGet(walletStorage *sync.Map) operations.MgmtWalletGetHandler {
	return &walletGet{walletStorage: walletStorage}
}

type walletGet struct {
	walletStorage *sync.Map
}

//nolint:nolintlint,ireturn
func (c *walletGet) Handle(params operations.MgmtWalletGetParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewMgmtWalletGetInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletGetWallets,
				Message: err.Error(),
			})
	}

	var wal []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		walletss := &models.Wallet{
			Nickname: &wallets[i].Nickname,
			Address:  &wallets[i].Address,
			KeyPairs: []*models.WalletKeyPairsItems0{},
		}

		wal = append(wal, walletss)
	}

	return operations.NewMgmtWalletGetOK().WithPayload(wal)
}
