package wallet

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/wallet"
)

func NewImport(walletStorage *sync.Map) operations.MgmtWalletImportHandler {
	return &wImport{walletStorage: walletStorage}
}

type wImport struct {
	walletStorage *sync.Map
}

func (c *wImport) Handle(params operations.MgmtWalletImportParams) middleware.Responder {
	var err error

	_, ok := c.walletStorage.Load(*params.Body.Nickname)
	if ok {
		return operations.NewMgmtWalletImportInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletAlreadyExists,
				Message: "Error: a wallet with the same nickname already exists.",
			})
	}

	keyPairs := make([]wallet.KeyPair, len(params.Body.KeyPairs))
	for i := 0; i < len(params.Body.KeyPairs); i++ {
		keyPairs[i].PrivateKey, err = base58.CheckDecode(*params.Body.KeyPairs[i].PrivateKey)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		keyPairs[i].PublicKey, err = base58.CheckDecode(*params.Body.KeyPairs[i].PublicKey)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		salt, err := base58.CheckDecode(*params.Body.KeyPairs[i].Salt)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		copy(keyPairs[i].Salt[:], salt)

		nonce, err := base58.CheckDecode(*params.Body.KeyPairs[i].Nonce)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		copy(keyPairs[i].Nonce[:], nonce)

		keyPairs[i].Protected = true
	}

	newWallet := wallet.Wallet{Version: 0, Nickname: *params.Body.Nickname, Address: *params.Body.Address, KeyPairs: keyPairs}

	c.walletStorage.Store(newWallet.Nickname, newWallet)

	bytesOutput, err := json.Marshal(newWallet)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletImportNew,
				Message: err.Error(),
			})
	}

	err = os.WriteFile("wallet_"+*params.Body.Nickname+".json", bytesOutput, 0o644)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletImportNew,
				Message: err.Error(),
			})
	}

	return operations.NewMgmtWalletImportNoContent()
}
