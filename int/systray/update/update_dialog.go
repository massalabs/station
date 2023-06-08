package update

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/massalabs/thyra/int/systray/utils"
)

const (
	// updateWindowTitle is the title of the update dialog. It is used to find the dialog window or to create it.
	updateWindowTitle = "MassaStation update available"

	// downloadURL is the URL to the latest release on GitHub.
	downloadURL = "https://github.com/massalabs/thyra/releases/latest"
)

// Finds the update dialog window by its title.
func findUpdateWindow(app *fyne.App) fyne.Window {
	var window fyne.Window

	for _, w := range (*app).Driver().AllWindows() {
		if w.Title() == updateWindowTitle {
			window = w

			break
		}
	}

	return window
}

// Creates the update dialog.
func createUpdateDialog(window fyne.Window, app *fyne.App) dialog.Dialog {
	dialog := dialog.NewCustomConfirm(
		"An update is available for MassaStation. Do you want to update now?",
		"Update",
		"Cancel",
		widget.NewLabel("If you choose to update, MassaStation will be closed and every running task will be stopped."),
		func(b bool) {
			window.Close()
			if b {
				utils.OpenURL(app, downloadURL)
				(*app).Quit()
			}
		},
		window,
	)

	return dialog
}

// Opens the update dialog.
//
// If the dialog is already opened, it will be raised.
func OpenUpdateDialog(app *fyne.App) {
	window := findUpdateWindow(app)

	// If the dialog is already opened, we just want to raise it.
	if window != nil {
		window.Show()
		window.RequestFocus()

		return
	}

	window = (*app).NewWindow(updateWindowTitle)
	dialog := createUpdateDialog(window, app)

	// Resize the window to fit the dialog.
	window.Resize(dialog.MinSize())
	window.SetFixedSize(true)

	window.Show()
	dialog.Show()
}
