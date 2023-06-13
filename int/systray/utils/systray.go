package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func SetSystemTrayIcon(app *fyne.App, icon *fyne.StaticResource) {
	if desk, ok := (*app).(desktop.App); ok {
		desk.SetSystemTrayIcon(icon)
	}
}
