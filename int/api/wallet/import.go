package wallet

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/gui"
	"github.com/massalabs/thyra/pkg/wallet"
)

const fileModeUserRW = 0o600

//nolint:nolintlint,ireturn
func NewImport(walletStorage *sync.Map, app *fyne.App) operations.MgmtWalletImportHandler {
	return &wImport{walletStorage: walletStorage, app: app}
}

type wImport struct {
	walletStorage *sync.Map
	app           *fyne.App
}

//nolint:nolintlint,ireturn,funlen
func (c *wImport) Handle(params operations.MgmtWalletImportParams) middleware.Responder {
	password, walletName, privateKey, err := gui.AskWalletInfo(c.app)
	if err != nil {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNew,
				Message: err.Error(),
			})
	}

	if len(walletName) == 0 {
		return operations.NewMgmtWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	_, inStore := c.walletStorage.Load(walletName)
	if inStore {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletAlreadyExists,
				Message: "Error: a wallet with the same nickname already exists.",
			})
	}

	if len(password) == 0 {
		return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletCreateNoPassword,
				Message: "Error: password field is mandatory.",
			})
	}

	newWallet, RequestError := wallet.Imported(walletName, privateKey)
	if RequestError.Err != nil {
		switch {
		case RequestError.StatusCode == wallet.ImportedAlreadyImported:
			return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeWalletAlreadyImported,
					Message: RequestError.Err.Error(),
				})
		case RequestError.StatusCode == wallet.ImportedEncodingB58Error:
			return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeWalletEncodingB58E,
					Message: RequestError.Err.Error(),
				})
		case RequestError.StatusCode == wallet.ImportedLoadingWalletsError:
			return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeWalletLoadingWallets,
					Message: RequestError.Err.Error(),
				})
		default:
			return operations.NewMgmtWalletCreateInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCodeWalletCreateNew,
					Message: RequestError.Err.Error(),
				})
		}
	}

	return CreateNewWallet(&walletName, &password, c.walletStorage, newWallet)
}
