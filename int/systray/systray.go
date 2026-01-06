package systray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fyneDesktop "fyne.io/fyne/v2/driver/desktop"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/systray/embedded"
	"github.com/massalabs/station/int/systray/utils"
)

func MakeGUI() (fyne.App, *fyne.Menu) {
	stationGUI := app.New()

	if desk, ok := stationGUI.(fyneDesktop.App); ok {
		icon := fyne.NewStaticResource("logo", embedded.Logo)

		homeShortCutMenu := fyne.NewMenuItem("Open MassaStation", nil)
		homeShortCutMenu.Action = func() {
			utils.OpenURL(&stationGUI, "https://"+config.MassaStationURL)
		}

		menu := fyne.NewMenu(
			"MassaStation",
			fyne.NewMenuItemSeparator(),
			homeShortCutMenu,
		)

		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)

		return stationGUI, menu
	}

	return stationGUI, nil
}
