package systray

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	//nolint:goimports,nolintlint
	"fyne.io/fyne/v2/app"
	//nolint:typecheck,nolintlint
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/systray/embedded"
	"github.com/massalabs/station/int/systray/utils"
)

func MakeGUI() (fyne.App, *fyne.Menu) {
	stationGUI := app.New()
	menu := fyne.NewMenu("MassaStation")

	msURL := fmt.Sprintf("https://%s", config.MassaStationURL)

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
			utils.OpenURL(&stationGUI, msURL)
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

	// Set up macOS dock click handler
	// if runtime.GOOS == "darwin" {

	// 	fmt.Println(">>>>>>>>>>>>>MassaStation is running on macOS")
	// 	stationGUI.Lifecycle().SetOnEnteredForeground(func() {
	// 		fmt.Println(">>>>>>>>>>>>>MassaStation is in the foreground")
	// 		logger.Debugf(">>>>>>>>>>>>>Dock icon clicked211111112")

	// 		openURL(msURL)
	// 		utils.OpenURL(&stationGUI, msURL)
	// 	})
	// 	go handleMacDockClick()
	// }

	stationGUI.Run()
	return stationGUI, menu
}

func openURL(url string) {
	cmd := exec.Command("open", url)

	if err := cmd.Start(); err != nil {
		println("Error opening URL:", err.Error())
	}
}

// 🖥 macOS: Detect Dock icon click
// func handleMacDockClick() {

// 	macos.RunApp(func(app appkit.Application, delegate *appkit.ApplicationDelegate) {
// 		// app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
// 		// app.ActivateIgnoringOtherApps(true)

// 		// Set the delegate to handle the application lifecycle
// 		delegate.SetApplicationShouldHandleReopenHasVisibleWindows(func(app appkit.Application, _ bool) bool {
// 			// Handle the Dock icon click
// 			// This is where you can open the URL or perform any action
// 			// openURL("https://example.com")
// 			fmt.Println(">>>>SetApplicationShouldHandleReopenHasVisibleWindows><YOOO Dock icon clicked")
// 			logger.Debugf(">>>>>>>>>>>>>Dock icon clicked")
// 			openURL("https://www.google.com")
// 			return true
// 		})

// 		delegate.SetApplicationDidBecomeActive(func(notification foundation.Notification) {
// 			logger.Debugf("Dock icon clicked222")

// 			fmt.Println(">>>>>SetApplicationDidBecomeActive<YOOO Dock icon clicked")
// 			openURL("https://www.google.com")
// 		})

// 		// Set the delegate to handle the application lifecycle
// 		// app.SetDelegate(delegate)

// 	})

// 	//ApplicationOpenURLs/

// 	// cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
// 	// 	cocoa.NSApp_SetDelegate(cocoa.AppDelegate{
// 	// 		ApplicationShouldHandleReopen: func(_ objc.Object, _ bool) bool {
// 	// 			openURL("https://example.com")
// 	// 			return true
// 	// 		},
// 	// 	})
// 	// }).Run()
// }
