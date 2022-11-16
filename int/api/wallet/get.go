package wallet

import (
	"strconv"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/ledger"
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
	client := node.NewDefaultClient()

	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewMgmtWalletGetInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletGetWallets,
				Message: err.Error(),
			})
	}

	var addresses []string = make([]string, len(wallets))
	for i, wallet := range wallets {
		addresses[i] = wallet.Address
	}

	var wal []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		var address, err = ledger.Addresses(client, []string{wallets[i].Address})
		if err != nil {
			return operations.NewMgmtWalletGetInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeWalletGetBalance,
					Message: err.Error(),
				})
		}

		var balance, error = strconv.ParseFloat(address[0].CandidateBalanceInfo, 64)
		if error != nil {
			return operations.NewMgmtWalletGetInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeWalletGetBalance,
					Message: error.Error(),
				})
		}

		walletss := &models.Wallet{
			Nickname: &wallets[i].Nickname,
			Address:  &wallets[i].Address,
			KeyPairs: []*models.WalletKeyPairsItems0{},
			Balance:  balance,
		}

		wal = append(wal, walletss)
	}

	return operations.NewMgmtWalletGetOK().WithPayload(wal)
}
