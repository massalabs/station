package wallet

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewGet(walletStorage *sync.Map) operations.MgmtWalletGetHandler {
	return &walletGet{walletStorage: walletStorage}
}

type walletGet struct {
	walletStorage *sync.Map
}

// TODO Clean the struct mapping here + correct KeyPairs not returned
func (c *walletGet) Handle(params operations.MgmtWalletGetParams) middleware.Responder {

	wallets, err := wallet.ReadWallets()
	var walll []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		walletss := &models.Wallet{
			Nickname: &wallets[i].Nickname,
			Address:  &wallets[i].Address,
			KeyPairs: []*models.WalletKeyPairsItems0{}}
		walll = append(walll, walletss)
	}

	if err != nil {
		panic(err)
	}

	return operations.NewMgmtWalletGetOK().WithPayload(walll)
}
