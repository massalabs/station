package wallet

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/gui"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/wallet"
)

const fileModeUserRW = 0o600

//nolint:nolintlint,ireturn
func NewImport(app *fyne.App) operations.MgmtWalletImportHandler {
	return &wImport{app: app}
}

type wImport struct {
	walletStorage *sync.Map
	app           *fyne.App
}

//nolint:nolintlint,ireturn,funlen
func (c *wImport) Handle(params operations.MgmtWalletImportParams) middleware.Responder {
	password, walletName, pk, err := gui.AskWalletInfo(c.app)
	if err != nil {
		panic(err)
	}

	newWallet, err := wallet.Imported(walletName, pk)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNew,
				Message: err.Error(),
			})
	}

	err = newWallet.Protect(password, 0)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNew,
				Message: err.Error(),
			})
	}

	bytesOutput, err := json.Marshal(newWallet)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNew,
				Message: err.Error(),
			})
	}

	err = os.WriteFile(wallet.GetWalletDirectory()+"wallet_"+walletName+".json", bytesOutput, fileModeUserRW)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNew,
				Message: err.Error(),
			})
	}
	fmt.Println(password, walletName, pk, newWallet)

	c.walletStorage.Store(newWallet.Nickname, newWallet)

	privK := base58.CheckEncode(newWallet.KeyPairs[0].PrivateKey)
	pubK := base58.CheckEncode(newWallet.KeyPairs[0].PublicKey)
	salt := base58.CheckEncode(newWallet.KeyPairs[0].Salt[:])
	nonce := base58.CheckEncode(newWallet.KeyPairs[0].Nonce[:])

	return operations.NewMgmtWalletCreateOK().WithPayload(
		&models.Wallet{
			Nickname: &newWallet.Nickname,
			Address:  &newWallet.Address,
			KeyPairs: []*models.WalletKeyPairsItems0{{
				PrivateKey: &privK,
				PublicKey:  &pubK,
				Salt:       &salt,
				Nonce:      &nonce,
			}},
			Balance: 0,
		})
}
