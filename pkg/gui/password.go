package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func AskPassword(nickname string, app *fyne.App) string {
	return Password(nickname, app)
}

func AskPasswordDeleteWallet(nickname string, app *fyne.App) string {
	return PasswordDeleteWallet(nickname, app)

}

// inspired by https://hackernoon.com/asyncawait-in-golang-an-introductory-guide-ol1e34sg

// Thyra password input dialog.
func PasswordDialog(nickname string, app *fyne.App) chan string {
	passwordText := make(chan string)

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
			passwordText <- password.Text
		},
		OnCancel: func() {
			passwordText <- ""
			window.Hide()
		},
		SubmitText: "Submit",
		CancelText: "Cancel transaction",
	}

	window.SetContent(form)
	window.CenterOnScreen()
	window.Canvas().Focus(password)
	window.Show()

	return passwordText
}

func PasswordDeleteDialog(nickname string, app *fyne.App) chan string {
	passwordText := make(chan string)

	window := (*app).NewWindow("Massa - Thyra")

	width := 250.0
	height := 250.0

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
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	text1 := canvas.NewText("Delete Wallet ?", black)
	text1.TextSize = 25
	title := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), text1, layout.NewSpacer())
	text2 := canvas.NewText("If you delete a wallet, you will lose your MAS associated to it and ", black)
	text3 := canvas.NewText("won't be able to edit websites linked to this wallet anymore", black)

	content := container.New(layout.NewVBoxLayout(), text2, layout.NewSpacer(), text3)

	centeredForm := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), form, layout.NewSpacer())
	window.SetContent(container.New(layout.NewVBoxLayout(), title, content, centeredForm))

	//window.SetContent(content)
	window.CenterOnScreen()
	window.Canvas().Focus(password)
	window.Show()

	return passwordText
}

// This function is blocking, it returns when the user submit or cancel the form.
func Password(nickname string, app *fyne.App) string {
	return <-PasswordDialog(nickname, app)
}

func PasswordDeleteWallet(nickname string, app *fyne.App) string {
	return <-PasswordDeleteDialog(nickname, app)
}
