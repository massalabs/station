package gui

import (
	"github.com/pwiecz/go-fltk"
)

func Password(nickname string) string {

	password := ""
	win := fltk.NewWindow(400, 160)
	win.SetLabel("Thyra - " + nickname)
	valuePassword := fltk.NewInput(120, 50, 200, 30, "Password")
	btn := fltk.NewButton(160, 110, 80, 30, "Click")
	btn.SetCallback(func() {
		password = valuePassword.Value()
		win.Destroy()
	})

	win.End()
	win.Show()
	fltk.Run()
	return password
}
