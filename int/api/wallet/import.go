package wallet

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/gui"
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
	var err error
	password, walletName, pk, err := gui.AskWalletInfo(c.app)
	if err != nil {
		panic(err)
	}
	fmt.Println(password, walletName, pk)

	return operations.NewMgmtWalletImportNoContent()
}
