package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Struct to return the password and an error to know if the user submitted or canceled the form.
type PasswordEntry struct {
	ClearPassword string
	Err           error
}

func AskPassword(nickname string, app *fyne.App) (string, error) {
	return Password(nickname, app)
}

func AskPasswordDeleteWallet(nickname string, app *fyne.App) string {
	return PasswordDeleteWallet(nickname, app)
}

// inspired by https://hackernoon.com/asyncawait-in-golang-an-introductory-guide-ol1e34sg

// Thyra password input dialog.
func PasswordDialog(nickname string, app *fyne.App) chan PasswordEntry {
	passwordEntry := make(chan PasswordEntry)
	window := (*app).NewWindow("Massa - Thyra")
	width := 250.0
	height := 90.0

	window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})

	password := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Password", password),
	}

	//nolint:exhaustruct
	form := &widget.Form{
		Items: items,
		OnSubmit: func() {
			window.Hide()
			passwordEntry <- PasswordEntry{password.Text, nil}
		},
		OnCancel: func() {
			passwordEntry <- PasswordEntry{ClearPassword: "", Err: errors.New("password entry cancelled by the user")}
			window.Hide()
		},
		SubmitText: "Submit",
		CancelText: "Cancel transaction",
	}

	window.SetContent(form)
	window.CenterOnScreen()
	window.Canvas().Focus(password)
	window.Show()

	return passwordEntry
}

func PasswordDeleteDialog(nickname string, app *fyne.App) chan string {
	passwordText := make(chan string)

	window := (*app).NewWindow("Massa - Thyra")

	width := 250.0
	height := 80.0

	window.Resize(fyne.Size{Width: float32(width), Height: float32(height)})

	password := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Password", password),
	}

	//nolint:exhaustruct
	form := &widget.Form{
		Items: items,
		OnSubmit: func() {
			window.Hide()
			passwordText <- password.Text
		},
		OnCancel: func() {
			passwordText <- ""
			window.Hide()
		},
		SubmitText: "Delete",
		CancelText: "Cancel",
	}
	spacer := layout.NewSpacer()
	text1 := widget.NewLabel(`Delete "` + nickname + `" Wallet ?`)
	title := container.New(layout.NewHBoxLayout(), spacer, text1, spacer)
	text2 := widget.NewLabel("If you delete a wallet, you will lose your MAS associated to it and ")
	text3 := widget.NewLabel("won't be able to edit websites linked to this wallet anymore ")
	content := container.New(layout.NewVBoxLayout(), text2, text3, spacer)
	centeredForm := container.New(layout.NewVBoxLayout(), spacer, form, spacer)
	window.SetContent(container.New(layout.NewVBoxLayout(), title, spacer, content, spacer, centeredForm, spacer))
	window.CenterOnScreen()
	window.Canvas().Focus(password)
	window.Show()

	return passwordText
}

// This function is blocking, it returns when the user submit or cancel the form.
func Password(nickname string, app *fyne.App) (string, error) {
	PasswordEntry := <-PasswordDialog(nickname, app)

	return PasswordEntry.ClearPassword, PasswordEntry.Err
}

func PasswordDeleteWallet(nickname string, app *fyne.App) string {
	return <-PasswordDeleteDialog(nickname, app)
}
