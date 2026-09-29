package systray

import (
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/systray/utils"
	"github.com/massalabs/station/pkg/logger"
)

const (
	enableHTTPSButtonLabel = "Enable HTTPS..."

	// httpsWindowTitle is the title of the HTTPS dialog. It is used to find the dialog window or to create it.
	httpsWindowTitle = "MassaStation is not secured"

	httpsNotTrustedMessage  = "The MassaStation certificate is not trusted by your system: MassaStation is served over http."
	httpsRetryFailedMessage = "The MassaStation certificate could not be trusted: MassaStation is still served over http."
)

// httpsEnabler retries to trust the MassaStation certificate and updates the systray accordingly.
type httpsEnabler struct {
	app         *fyne.App
	systrayMenu *fyne.Menu
	menuItem    *fyne.MenuItem
	retrying    atomic.Bool
}

// StartHTTPSCheck warns the user when the MassaStation certificates are not trusted by the operating system
// (e.g. the authorization prompt was dismissed) and adds a systray menu item to retry.
//
// A dialog is used rather than a notification: on macOS, notifications may be silently dropped
// (not permitted, unbundled binary...).
func StartHTTPSCheck(app *fyne.App, systrayMenu *fyne.Menu) {
	if config.IsHTTPSTrusted() {
		return
	}

	enabler := &httpsEnabler{app: app, systrayMenu: systrayMenu}

	// If systrayMenu is nil (e.g., desktop interface assertion failed), the dialog is the only way to retry.
	if systrayMenu != nil {
		enabler.menuItem = fyne.NewMenuItem(enableHTTPSButtonLabel, enabler.retry)
		systrayMenu.Items = append(systrayMenu.Items, enabler.menuItem)
		systrayMenu.Refresh()
	}

	enabler.openDialog(httpsNotTrustedMessage)
}

// retry retries to trust the certificate. On success, the systray menu item is removed and MassaStation is opened
// over https. On failure, the dialog is shown again.
func (e *httpsEnabler) retry() {
	// The authorization prompt would be displayed several times if clicked again while it is shown.
	if !e.retrying.CompareAndSwap(false, true) {
		return
	}

	// Retrying runs system commands waiting for the user to authorize: don't block the UI.
	go func() {
		defer e.retrying.Store(false)

		err := config.RetryCertificateTrust()
		if err != nil {
			logger.Warnf("failed to trust the MassaStation certificate: %s", err)

			fyne.Do(func() {
				e.openDialog(httpsRetryFailedMessage)
			})

			return
		}

		fyne.Do(func() {
			if e.menuItem != nil {
				removeMenuItem(e.systrayMenu, e.menuItem)
			}

			utils.OpenURL(e.app, config.StationURL())
		})
	}()
}

// openDialog opens the HTTPS dialog with the given message. If the dialog is already opened, it is raised.
func (e *httpsEnabler) openDialog(message string) {
	for _, w := range (*e.app).Driver().AllWindows() {
		if w.Title() == httpsWindowTitle {
			w.Show()
			w.RequestFocus()

			return
		}
	}

	window := (*e.app).NewWindow(httpsWindowTitle)

	confirm := dialog.NewCustomConfirm(
		message,
		"Enable HTTPS",
		"Continue with HTTP",
		widget.NewLabel("Enabling HTTPS requires your password to trust the MassaStation certificate."),
		func(enable bool) {
			window.Close()

			if enable {
				e.retry()
			}
		},
		window,
	)

	// Resize the window to fit the dialog.
	window.Resize(confirm.MinSize())
	window.SetFixedSize(true)

	window.Show()
	confirm.Show()
}

// removeMenuItem removes the given item from the menu.
func removeMenuItem(menu *fyne.Menu, item *fyne.MenuItem) {
	for i, menuItem := range menu.Items {
		if menuItem == item {
			menu.Items = append(menu.Items[:i], menu.Items[i+1:]...)
			menu.Refresh()

			return
		}
	}
}
