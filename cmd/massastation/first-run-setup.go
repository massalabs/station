package main

import (
	"fmt"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin"
	"github.com/massalabs/station/pkg/store"
)

const (
	walletPluginName = "Massa Wallet"
	dewebPluginName  = "Local DeWeb Provider"
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

	/* set the setup done to true to avoid running the setup again
	Following step of the function are not blocking and there is no problem if they fail
	*/
	cfg.StationFirstRunSetupDone = true
	if err := configManager.SaveConfig(cfg); err != nil {
		return fmt.Errorf("%s could not save config 'StationFirstRunSetupDone = true', got error: %s", prefixLog, err)
	}

	if err := installPlugin(pluginManager, walletPluginName); err != nil {
		logger.Warnf("%s %s", prefixLog, err)
	} else {
		logger.Infof("%s '%s' plugin installed", prefixLog, walletPluginName)
	}

	if err := installPlugin(pluginManager, dewebPluginName); err != nil {
		logger.Warnf("%s %s", prefixLog, err)
	} else {
		logger.Infof("%s '%s' plugin installed", prefixLog, dewebPluginName)
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
