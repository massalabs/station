package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type WalletInfoEntry struct {
	ClearPassword string
	WalletName    string
	PrivateKey    string
	Err           error
}

func AskWalletInfo(app *fyne.App) (string, string, string, error) {
	return WalletInfo(app)
}

func WalletInfo(app *fyne.App) (string, string, string, error) {
	WalletInfoEntry := <-LoadWalletDialog(app)

	return WalletInfoEntry.ClearPassword, WalletInfoEntry.WalletName, WalletInfoEntry.PrivateKey, WalletInfoEntry.Err
}

func LoadWalletDialog(app *fyne.App) chan WalletInfoEntry {
	walletInfoEntry := make(chan WalletInfoEntry)

	window := (*app).NewWindow("Massa - Thyra")

	width := 700.0
	height := 100.0

	window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})

	walletName := widget.NewEntry()
	password := widget.NewPasswordEntry()
	privateKey := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Wallet name", walletName),
		widget.NewFormItem("Password", password),
		widget.NewFormItem("Private key", privateKey),
	}

	//nolint:exhaustruct
	form := &widget.Form{
		Items: items,
		OnSubmit: func() {
			window.Hide()
			//nolint:gofumpt
			walletInfoEntry <- WalletInfoEntry{ClearPassword: password.Text,
				WalletName: walletName.Text,
				PrivateKey: privateKey.Text,
				Err:        nil}
		},
		OnCancel: func() {
			//nolint:gofumpt
			walletInfoEntry <- WalletInfoEntry{ClearPassword: "",
				WalletName: "", PrivateKey: "",
				Err: errors.New("wallet loading cancelled by the user")}
			window.Hide()
		},
		SubmitText: "Load",
		CancelText: "Cancel",
	}
	spacer := layout.NewSpacer()
	text := widget.NewLabel(`Load a Wallet`)
	title := container.New(layout.NewHBoxLayout(), spacer, text, spacer)
	centeredForm := container.New(layout.NewVBoxLayout(), spacer, form, spacer)
	window.SetContent(container.New(layout.NewVBoxLayout(), title, spacer, centeredForm, spacer))
	window.CenterOnScreen()
	window.Canvas().Focus(password)
	window.Show()

	return walletInfoEntry
}
