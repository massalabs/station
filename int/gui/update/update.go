package update

import (
	"strings"
	"time"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/gui/embedded"
	"github.com/massalabs/station/int/update"
	"github.com/massalabs/station/pkg/logger"
	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	updateCheckInterval = 1 * time.Hour
	downloadURL         = "https://station.massa.net"

	updateDialogTitle   = "An update is available for Massa Station. Do you want to update it now ?"
	updateDialogMessage = "If you choose to update, Massa Station will be closed and every running tasks will be stopped."
)

// Checks for updates.
// If an update is available, it changes the icon to notify the user
// and shows the update button in the system tray menu.
func Check(systray *application.SystemTray, updateButton *application.MenuItem) {
	if !update.Check() {
		return
	}

	systray.SetIcon(embedded.NotificationLogo)
	updateButton.SetHidden(false)
	updateButton.SetEnabled(true)
}

// Shows the update dialog.
// If the user chooses to update, it opens the download URL and shuts down the app.
// Otherwise, it does nothing.
func ShowUpdateDialog(appShutDown func()) {
	updateDialog := application.QuestionDialog()
	updateDialog.SetTitle(updateDialogTitle)
	updateDialog.SetMessage(updateDialogMessage)
	updateDialog.AddButton("Yes").OnClick(func() {
		logger.Info("Update Dialog: User chose to update")

		if browser.OpenURL(downloadURL) != nil {
			logger.Error("Update Dialog: An error occurred while opening the download URL")
		}

		appShutDown()
	})
	updateDialog.AddButton("No").SetAsDefault().OnClick(func() {
		logger.Debug("Update Dialog: User chose not to update")
	})

	logger.Debug("Update Dialog: Showing dialog")

	updateDialog.Show()
}

// Starts a ticker that checks for updates every updateCheckInterval.
func StartUpdateCheck(systray *application.SystemTray, updateButton *application.MenuItem) {
	// We don't want to check for updates in dev mode.
	if strings.Contains(config.Version, "dev") {
		return
	}

	// We check for updates on startup.
	Check(systray, updateButton)

	go func() {
		ticker := time.NewTicker(updateCheckInterval)
		for range ticker.C {
			Check(systray, updateButton)
		}
	}()
}
