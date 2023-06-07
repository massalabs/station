package systray

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra/pkg/config"

	// nolint:typecheck
	"fyne.io/fyne/v2/driver/desktop"
)

func MakeGUI() (fyne.App, *fyne.Menu) {
	stationGUI := app.New()
	menu := fyne.NewMenu("MassaStation")
	_ = stationGUI.NewWindow("MassaStation")

	if desk, ok := stationGUI.(desktop.App); ok {
		icon := fyne.NewStaticResource("logo", embeded.Logo)
		titleMenu := fyne.NewMenuItem("MassaStation", nil)
		homeShortCutMenu := fyne.NewMenuItem("MassaStation home", nil)
		testMenu := fyne.NewMenuItem("Test", nil)

		titleMenu.Disabled = true
		titleMenu.Icon = icon

		testMenu.Action = func() {
			notification := fyne.NewNotification("Test notification", "This is a test notification from MassaStation")
			stationGUI.SendNotification(notification)
		}

		homeShortCutMenu.Action = func() {
			openURL(&stationGUI, fmt.Sprintf("http://%s", config.MassaStationURL))
		}

		menu.Items = append(menu.Items,
			titleMenu,
			fyne.NewMenuItemSeparator(),
			homeShortCutMenu,
			// testMenu,
			fyne.NewMenuItemSeparator(),
		)

		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)
	}

	return stationGUI, menu
}
