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

const fileModeUserRW = 0o600

//nolint:nolintlint,ireturn
func NewImport(walletStorage *sync.Map) operations.MgmtWalletImportHandler {
	return &wImport{walletStorage: walletStorage}
}

type wImport struct {
	walletStorage *sync.Map
}

//nolint:nolintlint,ireturn,funlen
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
	for index := 0; index < len(params.Body.KeyPairs); index++ {
		keyPairs[index].PrivateKey, err = base58.CheckDecode(*params.Body.KeyPairs[index].PrivateKey)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		keyPairs[index].PublicKey, err = base58.CheckDecode(*params.Body.KeyPairs[index].PublicKey)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		salt, err := base58.CheckDecode(*params.Body.KeyPairs[index].Salt)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		copy(keyPairs[index].Salt[:], salt)

		nonce, err := base58.CheckDecode(*params.Body.KeyPairs[index].Nonce)
		if err != nil {
			return operations.NewMgmtWalletCreateUnprocessableEntity()
		}

		copy(keyPairs[index].Nonce[:], nonce)

		keyPairs[index].Protected = true
	}

	newWallet := wallet.Wallet{
		Version:  0,
		Nickname: *params.Body.Nickname,
		Address:  *params.Body.Address,
		KeyPairs: keyPairs,
	}

	c.walletStorage.Store(newWallet.Nickname, newWallet)

	bytesOutput, err := json.Marshal(newWallet)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletImportNew,
				Message: err.Error(),
			})
	}

	err = os.WriteFile("wallet_"+*params.Body.Nickname+".json", bytesOutput, fileModeUserRW)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletImportNew,
				Message: err.Error(),
			})
	}

	return operations.NewMgmtWalletImportNoContent()
}
