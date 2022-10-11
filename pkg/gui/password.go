package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

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

// This function is blocking, it returns when the user submit or cancel the form.
func Password(nickname string, app *fyne.App) string {
	return <-PasswordDialog(nickname, app)
}
