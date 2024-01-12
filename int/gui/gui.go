package gui

import (
	"embed"

	"github.com/massalabs/station/int/api"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/gui/embedded"
	"github.com/massalabs/station/int/gui/update"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// App is the GUI application.
type App struct {
	// The application
	app *application.App
	// The main window of the app
	window *application.WebviewWindow
	// The app system tray
	systray *application.SystemTray
	// The menu attached to the system tray
	menu *application.Menu
	// The update button, used to update Massa Station from the system tray
	updateButton *application.MenuItem
}

// NewApp creates a new App.
func NewApp(
	server *api.Server,
	networkManager *config.NetworkManager,
	pluginManager *plugin.Manager,
	assets embed.FS,
) *App {
	app := application.New(application.Options{
		Name:        "Massa Station",
		Description: "Your gateway to the decentralized web",
		Assets: application.AssetOptions{
			FS: assets,
		},
		Icon: embedded.Logo,
		PanicHandler: func(any) {
			logger.Error("Wails Panicked - Please check the logs for more information")
			quitApp(nil, server, pluginManager)
		},
	})
	if app == nil {
		logger.Fatal("Unable to create application")

		return nil
	}

	app.On(events.Common.ApplicationStarted, func(event *application.Event) {
		server.Start(networkManager, pluginManager)

		err := pluginManager.RunAll()
		if err != nil {
			logger.Fatalf("while running all plugins: %w", err)
		}
	})

	window := makeWindow(app)
	if window == nil {
		logger.Fatal("Unable to create window")

		return nil
	}

	systray, menu, updateButton := makeSystray(app, server, networkManager, pluginManager)

	if systray == nil {
		logger.Fatal("Unable to create system tray")

		return nil
	} else if menu == nil {
		logger.Fatal("Unable to create menu")

		return nil
	}

	return &App{
		app:          app,
		window:       window,
		menu:         menu,
		systray:      systray,
		updateButton: updateButton,
	}
}

// Starts the app.
func (a *App) Run() error {
	update.StartUpdateCheck(a.systray, a.updateButton)

	err := a.app.Run()
	if err != nil {
		logger.Fatalf("Unable to run app: %s", err.Error())
	}

	return nil
}
