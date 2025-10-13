package main

import (
	"fmt"
	"strings"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
	"github.com/massalabs/station/pkg/store"
)

const (
	walletPluginName = "Massa Wallet"
	dewebPluginName  = "Local DeWeb Provider"
	pluginAuthor     = "Massa Labs"
)

// stationFirstRunSetup process some setup task the first time station is started
func stationFirstRunSetup(configManager *config.MSConfigManager, pluginManager *plugin.Manager) error {
	prefixLog := "Station First Run Setup:"

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("%s could not load config, got error: %s", prefixLog, err)
	}

	// if the setup has already been done, exit the function
	if cfg.StationFirstRunSetupDone {
		logger.Debugf("%s setup already done", prefixLog)
		return nil
	}

	// check if the wallet and deweb plugins are installed
	pluginIds := pluginManager.ID()
	walletInstalled := false
	dewebInstalled := false
	for _, pluginId := range pluginIds {
		plugin, err := pluginManager.Plugin(pluginId)
		if err != nil {
			return fmt.Errorf("could not get plugin %s, got error: %s", pluginId, err)
		}
		if plugin != nil && strings.Contains(plugin.Path, "wallet") {
			walletInstalled = true
		}
		if plugin != nil && strings.Contains(plugin.Path, "deweb") {
			dewebInstalled = true
		}
	}

	/* set the setup done to true to avoid running the setup again
	Following step of the function are not blocking and there is no problem if they fail
	*/
	cfg.StationFirstRunSetupDone = true
	if err := configManager.SaveConfig(cfg); err != nil {
		return fmt.Errorf("%s could not save config 'StationFirstRunSetupDone = true', got error: %s", prefixLog, err)
	}

	if !walletInstalled {
		logger.Infof("%s '%s' plugin not installed, installing...", prefixLog, walletPluginName)
		// install the wallet plugin
		if err := installPlugin(pluginManager, walletPluginName); err != nil {
			logger.Warnf("%s %s", prefixLog, err)
		} else {
			logger.Infof("%s '%s' plugin installed", prefixLog, walletPluginName)
		}

	}

	if !dewebInstalled {
		logger.Infof("%s '%s' plugin not installed, installing...", prefixLog, dewebPluginName)
		// install the deweb plugin
		if err := installPlugin(pluginManager, dewebPluginName); err != nil {
			logger.Warnf("%s %s", prefixLog, err)
		} else {
			logger.Infof("%s '%s' plugin installed", prefixLog, dewebPluginName)
		}
	}

	return nil
}

func installPlugin(pluginManager *plugin.Manager, pluginName string) error {
	pluginStore := store.StoreInstance.FindPluginByName(pluginName)
	if pluginStore == nil {
		return fmt.Errorf("%s plugin not found in store", pluginName)
	}
	pluginUrl, _, _, err := pluginStore.GetDLChecksumAndOs()
	if err != nil {
		return fmt.Errorf("could not retrieve %s plugin URL, got error: %s", pluginName, err)
	}
	if err := pluginManager.Install(pluginUrl); err != nil {
		return fmt.Errorf("could not install %s plugin, got error: %s", pluginName, err)
	}
	return nil
}
