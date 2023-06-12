package systray

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/massalabs/thyra/pkg/config"

	// nolint:typecheck
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/massalabs/thyra/int/systray/embedded"
	"github.com/massalabs/thyra/int/systray/utils"
)

func MakeGUI() (fyne.App, *fyne.Menu) {
	stationGUI := app.New()
	menu := fyne.NewMenu("MassaStation")

	if desk, ok := stationGUI.(desktop.App); ok {
		icon := fyne.NewStaticResource("logo", embedded.Logo)
		titleMenu := fyne.NewMenuItem("MassaStation", nil)
		homeShortCutMenu := fyne.NewMenuItem("Open MassaStation", nil)
		testMenu := fyne.NewMenuItem("Test", nil)

		titleMenu.Disabled = true
		titleMenu.Icon = icon

		testMenu.Action = func() {
			notification := fyne.NewNotification("Test notification", "This is a test notification from MassaStation")
			stationGUI.SendNotification(notification)
		}

		homeShortCutMenu.Action = func() {
			utils.OpenURL(&stationGUI, fmt.Sprintf("httpS://%s", config.MassaStationURL))
		}

		menu.Items = append(menu.Items,
			titleMenu,
			fyne.NewMenuItemSeparator(),
			homeShortCutMenu,
			// testMenu,
		)

		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)
	}

	return stationGUI, menu
}
