package gui

import (
	"fmt"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/gui/embedded"
	"github.com/massalabs/station/int/gui/update"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	menuLabelUpdate       = "Update Massa Station"
	menuLabelQuit         = "Quit"
	menuLabelCheckUpdates = "Check for updates"
)

func makeSystray(
	app *application.App,
	server *api.Server,
	networkManager *config.NetworkManager,
	pluginManager *plugin.Manager,
) (*application.SystemTray, *application.Menu, *application.MenuItem) {
	systray := app.NewSystemTray()
	systray.SetIcon(embedded.Logo)

	menu, updateButton := makeMenu(app, systray, server, networkManager, pluginManager)

	systray.SetMenu(menu)

	return systray, menu, updateButton
}

// makeMenu creates the system tray menu
// NOTE: The menu items are created in the order they are added to the menu.
func makeMenu(
	app *application.App,
	systray *application.SystemTray,
	server *api.Server,
	networkManager *config.NetworkManager,
	pluginManager *plugin.Manager,
) (*application.Menu, *application.MenuItem) {
	trayMenu := app.NewMenu()

	addMenuItem(trayMenu, "Open Massa Station", true, func(ctx *application.Context) {
		app.GetWindowByName("main").Focus()
	})

	trayMenu.AddSeparator()

	makeNetworkSubMenu(trayMenu, networkManager)

	trayMenu.AddSeparator()
	updateButton := addMenuItem(trayMenu, menuLabelUpdate, true, func(ctx *application.Context) {
		update.ShowUpdateDialog(func() {
			quitApp(app, server, pluginManager)
		})
	}).SetHidden(true).SetEnabled(false)

	trayMenu.AddSeparator()
	addMenuItem(trayMenu, fmt.Sprintf("Version: %s", config.Version), false, nil)
	addMenuItem(trayMenu, menuLabelCheckUpdates, config.Version != "dev", func(ctx *application.Context) {
		update.Check(systray, updateButton)
	})

	trayMenu.AddSeparator()
	addMenuItem(trayMenu, menuLabelQuit, true, func(ctx *application.Context) {
		quitApp(app, server, pluginManager)
	})

	return trayMenu, updateButton
}

func addMenuItem(
	menu *application.Menu,
	label string,
	enabled bool,
	onClick func(ctx *application.Context),
) *application.MenuItem {
	item := menu.Add(label)

	item.SetEnabled(enabled)

	if onClick != nil {
		item.OnClick(onClick)
	}

	return item
}

// makeNetworkSubMenu creates the network submenu
// NOTE: This is a mockup, the network submenu isn't ready yet.
func makeNetworkSubMenu(trayMenu *application.Menu, _ *config.NetworkManager) {
	subMenu := trayMenu.AddSubmenu("Network")
	addMenuItem(subMenu, "Available nodes", false, nil)
	subMenu.AddSeparator()

	radioCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		menuItem.SetChecked(true)
		subMenu.Update()
		logger.Debug("MOCKUP: Network changed to ", menuItem.Label())
	}

	subMenu.AddRadio("My node", true).OnClick(radioCallback)
	subMenu.AddRadio("Mainnet Node", false).OnClick(radioCallback)
	subMenu.AddRadio("Buildnet Node", false).OnClick(radioCallback)
}

// quitApp quits the app gracefully, after stopping the API server and all plugins.
func quitApp(app *application.App, server *api.Server, pluginManager *plugin.Manager) {
	if pluginManager != nil {
		logger.Debug("SysTray: Quitting... Stopping all plugins")
		pluginManager.StopAll()
	}

	if server != nil {
		logger.Debug("SysTray: Quitting... Stopping server")
		server.Stop()
	}

	if app != nil {
		logger.Debug("SysTray: Quitting... Quitting Wails app")
		app.Quit()
	}
}
