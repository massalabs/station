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

func NewDelete(walletStorage *sync.Map, app *fyne.App) operations.MgmtWalletDeleteHandler {
	return &walletDelete{walletStorage: walletStorage, app: app}
}

type walletDelete struct {
	walletStorage *sync.Map
	app           *fyne.App
}

//nolint:nolintlint,ireturn
func (c *walletDelete) Handle(params operations.MgmtWalletDeleteParams) middleware.Responder {

	walletLoaded, err := wallet.Load(params.Nickname)
	if err != nil {
		return createInternalServerError(errorCodeGetWallet, err.Error())
	}

	if len(params.Nickname) == 0 {
		return operations.NewMgmtWalletDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodeWalletDeleteNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	password := gui.AskPasswordDeleteWallet(params.Nickname, c.app)

	err = walletLoaded.Unprotect(password, 0)
	if err != nil {
		return createInternalServerError(errorCodeWalletWrongPassword, err.Error())
	}

	c.walletStorage.Delete(params.Nickname)

	err = wallet.Delete(params.Nickname)
	if err != nil {
		return createInternalServerError(errorCodeWalletDeleteFile, err.Error())

	}

	return operations.NewMgmtWalletDeleteNoContent()
}

//nolint:nolintlint,ireturn
func createInternalServerError(errorCode string, errorMessage string) middleware.Responder {
	return operations.NewWebsiteCreatorPrepareInternalServerError().
		WithPayload(
			&models.Error{
				Code:    errorCode,
				Message: errorMessage,
			})
}
