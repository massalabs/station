package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func AskPassword(nickname string, app *fyne.App) (string, bool) {
	return Password(nickname, app)
}

// inspired by https://hackernoon.com/asyncawait-in-golang-an-introductory-guide-ol1e34sg

// Thyra password input dialog.
func PasswordDialog(nickname string, app *fyne.App) (chan string, chan bool) {
	passwordText := make(chan string)
	isSubmitted := make(chan bool)
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
			isSubmitted <- true
		},
		OnCancel: func() {
			passwordText <- ""
			window.Hide()
			isSubmitted <- false
		},
		SubmitText: "Submit",
		CancelText: "Cancel transaction",
	}

	window.SetContent(form)
	window.CenterOnScreen()
	window.Canvas().Focus(password)
	window.Show()

	return passwordText, isSubmitted
}

// This function is blocking, it returns when the user submit or cancel the form.
func Password(nickname string, app *fyne.App) (string, bool) {
	passwordText, isSubmitted := PasswordDialog(nickname, app)
	return <-passwordText, <-isSubmitted
}
