package gui

import "github.com/wailsapp/wails/v3/pkg/application"

const (
	WindowName       = "main"
	windowTitle      = "Massa Station"
	windowWidth      = 1280
	windowHeight     = 720
	windowDefaultURL = "/"
)

// makeWindow creates the main window of the app.
func makeWindow(app *application.App) *application.WebviewWindow {
	return app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:   WindowName,
		Title:  windowTitle,
		Width:  windowWidth,
		Height: windowHeight,
		URL:    windowDefaultURL,
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()

			return false
		},
	})
}
