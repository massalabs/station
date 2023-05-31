package systray

import (
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	//nolint:typecheck
	"fyne.io/fyne/v2/driver/desktop"
)

func openURL(app *fyne.App, urlToOpen string) {
	u, err := url.Parse(urlToOpen)
	if err != nil {
		log.Fatal(err)
	}

	err = (*app).OpenURL(u)
	if err != nil {
		log.Fatal(err)
	}
}

func MakeGUI(logo []byte) (fyne.App, *fyne.Menu) {
	stationGUI := app.New()
	menu := fyne.NewMenu("MassaStation")

	if desk, ok := stationGUI.(desktop.App); ok {
		icon := fyne.NewStaticResource("logo", logo)
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
			openURL(&stationGUI, "http://station.massa/thyra/home")
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
