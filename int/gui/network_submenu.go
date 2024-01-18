package gui

import (
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type NetworkSubMenu struct {
	// The network submenu
	submenu *application.Menu

	// The network manager
	networkManager *config.NetworkManager

	// The network submenu items
	networkItems []*application.MenuItem
}

// makeNetworkSubMenu creates a network submenu.
func makeNetworkSubMenu(
	menu *application.Menu,
	networkManager *config.NetworkManager,
) *NetworkSubMenu {
	subMenu := menu.AddSubmenu("Network")
	addMenuItem(subMenu, "Available nodes", false, nil)
	subMenu.AddSeparator()

	radioCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		if err := networkManager.SwitchNetwork(menuItem.Label()); err != nil {
			logger.Error("Systray: Failed to switch network: ", err)
		}
	}

	networkItems := make([]*application.MenuItem, 0)

	for _, nodeName := range *networkManager.Networks() {
		networkItem := subMenu.AddRadio(nodeName, networkManager.Network().Network == nodeName)
		networkItem.OnClick(radioCallback)

		networkItems = append(networkItems, networkItem)
	}

	subMenu.AddSeparator()
	addMenuItem(subMenu, "Add network", false, nil)

	return &NetworkSubMenu{
		submenu:        subMenu,
		networkManager: networkManager,
		networkItems:   networkItems,
	}
}

// Update updates the network submenu with the current network.
func (n *NetworkSubMenu) Update() {
	network := n.networkManager.Network()

	for _, menuItem := range n.networkItems {
		if menuItem.Label() == network.Network {
			menuItem.SetChecked(true)
		} else {
			menuItem.SetChecked(false)
		}
	}
}
