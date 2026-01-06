package systray

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fyneDesktop "fyne.io/fyne/v2/driver/desktop"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/int/systray/embedded"
	"github.com/massalabs/station/int/systray/utils"
)

// networkMenuRefresher holds references needed to refresh the network menu.
type networkMenuRefresher struct {
	networkMenuItems []*fyne.MenuItem
	networkSubMenu   *fyne.Menu
	topLevelMenu     *fyne.Menu
	desk             fyneDesktop.App
	cfgMgr           *config.MSConfigManager
}

// buildNetworkMenuItems creates menu items for each network and wires the click handler.
// The onClick factory receives the network name and returns the action closure.
func buildNetworkMenuItems(
	networks []config.RPCInfos,
	current *config.RPCInfos,
	onClick func(name string) func(),
) []*fyne.MenuItem {
	items := make([]*fyne.MenuItem, 0, len(networks))

	for _, netInfo := range networks {
		networkName := netInfo.Name

		item := fyne.NewMenuItem(networkName, nil)
		if current != nil && current.Name == networkName {
			item.Checked = true
		}

		if onClick != nil {
			item.Action = onClick(networkName)
		}

		items = append(items, item)
	}

	return items
}

// refreshNetworkMenu rebuilds the network menu to reflect current networks and checkmarks.
func (n *networkMenuRefresher) refreshNetworkMenu() {
	if n == nil || n.cfgMgr == nil || n.networkSubMenu == nil {
		return
	}

	current := n.cfgMgr.CurrentNetwork()
	networks := n.cfgMgr.GetNetworkInfos()

	// Sort networks alphabetically by name
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	})

	// Rebuild menu items from current network list
	n.networkMenuItems = buildNetworkMenuItems(
		networks,
		current,
		func(name string) func() {
			return func() {
				// Switch the backend network. The registered callback on the
				// config manager will trigger a menu refresh.
				if err := n.cfgMgr.SwitchNetwork(name); err != nil {
					return
				}
			}
		},
	)

	// Update the submenu with new items
	n.networkSubMenu.Items = n.networkMenuItems
	n.networkSubMenu.Refresh()

	// Refresh top-level menu and re-set system tray menu
	if n.topLevelMenu != nil {
		n.topLevelMenu.Refresh()
	}
	if n.desk != nil {
		n.desk.SetSystemTrayMenu(n.topLevelMenu)
	}
}

// New creates the GUI and returns the systray menu.
func New() (fyne.App, *fyne.Menu) {
	stationGUI := app.New()

	if desk, ok := stationGUI.(fyneDesktop.App); ok {
		icon := fyne.NewStaticResource("logo", embedded.Logo)

		// Try to load network configuration for the "Network" submenu.
		var networkSubMenu *fyne.Menu
		var topLevelMenu *fyne.Menu
		var refresher *networkMenuRefresher
		var cfgMgr *config.MSConfigManager
		if mgr, err := config.GetConfigManager(); err == nil && mgr != nil {
			cfgMgr = mgr
			current := mgr.CurrentNetwork()
			var networkMenuItems []*fyne.MenuItem

			// Retrieve a lock-safe copy of the networks.
			networks := mgr.GetNetworkInfos()

			// Sort networks alphabetically by name.
			sort.Slice(networks, func(i, j int) bool {
				return networks[i].Name < networks[j].Name
			})

			// Build the initial set of menu items.
			networkMenuItems = buildNetworkMenuItems(
				networks,
				current,
				func(name string) func() {
					return func() {
						// Switch the backend network.
						if err := mgr.SwitchNetwork(name); err != nil {
							// If switching fails, leave the UI unchanged.
							return
						}
					}
				},
			)

			// Create the network submenu before wiring the refresher.
			if len(networkMenuItems) > 0 {
				networkSubMenu = fyne.NewMenu("Network", networkMenuItems...)

				// Create refresher for callback-based updates.
				refresher = &networkMenuRefresher{
					networkMenuItems: networkMenuItems,
					networkSubMenu:   networkSubMenu,
					cfgMgr:           mgr,
				}

				// Now set up the action closures with access to the submenu and desktop app.
				for i, netInfo := range networks {
					networkName := netInfo.Name
					item := networkMenuItems[i]

					// Capture networkName by value in the closure.
					item.Action = func(name string) func() {
						return func() {
							// Switch the backend network.
							if err := mgr.SwitchNetwork(name); err != nil {
								// If switching fails, leave the UI unchanged.
								return
							}
						}
					}(networkName)
				}
			}
		}

		homeShortCutMenu := fyne.NewMenuItem("Open MassaStation", nil)
		homeShortCutMenu.Action = func() {
			utils.OpenURL(&stationGUI, "https://"+config.MassaStationURL)
		}

		// Build the top-level systray menu in an idiomatic Fyne way.
		menuItems := []*fyne.MenuItem{
			fyne.NewMenuItemSeparator(),
			homeShortCutMenu,
		}

		// If we have networks, add a "Network" submenu.
		if networkSubMenu != nil {
			networkItem := fyne.NewMenuItem("Network", nil)
			networkItem.ChildMenu = networkSubMenu

			menuItems = append(
				menuItems,
				fyne.NewMenuItemSeparator(),
				networkItem,
			)
		}

		menu := fyne.NewMenu("MassaStation", menuItems...)
		topLevelMenu = menu

		// Update refresher with top-level menu and desktop app references
		if refresher != nil && cfgMgr != nil {
			refresher.topLevelMenu = topLevelMenu
			refresher.desk = desk

			// Register the refresher callback with the config manager
			cfgMgr.SetNetworkChangeCallback(func() {
				refresher.refreshNetworkMenu()
			})
		}

		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)

		return stationGUI, menu
	}

	return stationGUI, nil
}
