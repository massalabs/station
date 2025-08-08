package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
)

type NetworkStatus string

const (
	networkRefreshInterval                 = 5 * time.Minute
	NetworkStatusUp          NetworkStatus = "up"
	NetworkStatusDown        NetworkStatus = "down"
	legacy_networkConfigFile               = "config_network.yaml"
)

type RPCInfos struct {
	Name    string
	NodeURL string
	Version string
	ChainID uint64
	status  NetworkStatus
}

// NetworkInfos is an alias for backward compatibility
type NetworkInfos = RPCInfos

type NetworkConfig struct {
	currentNetwork *RPCInfos
	Networks       []RPCInfos
}

func fetchRPCInfos(url string) (string, uint64, error) {
	client := node.NewClient(url)

	status, err := node.Status(client)
	if err != nil {
		return "", 0, fmt.Errorf("failed to retrieve node status: %w", err)
	}

	version, err := node.GetVersionDigits(status)
	if err != nil {
		return "", 0, fmt.Errorf("getting network version: %w", err)
	}

	return version, uint64(*status.ChainID), nil
}

// StartNetworkRefresh starts a background routine that refreshes network information every 10 minutes
func StartNetworkRefresh(configManager *MSConfigManager) {
	go func() {
		ticker := time.NewTicker(networkRefreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			refreshNetworks(configManager)
		}
	}()
}

// refreshNetworks updates information for all networks in the config manager
func refreshNetworks(configManager *MSConfigManager) {
	configManager.configMutex.Lock()
	defer configManager.configMutex.Unlock()

	for i := range configManager.Network.Networks {
		refreshNetworkInfo(&configManager.Network.Networks[i])
	}
}

// refreshNetworkInfo updates information for a single network
func refreshNetworkInfo(network *RPCInfos) {
	version, chainID, err := fetchRPCInfos(network.NodeURL)
	if err != nil {
		logger.Warnf("Failed to refresh network %s: %v", network.Name, err)
		network.status = NetworkStatusDown
		return
	}

	network.Version = version
	network.ChainID = chainID
	network.status = NetworkStatusUp

	logger.Debugf("Refreshed network %s: version=%s, chainID=%d, status=%s",
		network.Name, version, chainID, network.status)
}

// checkLegacy checks if the legacy network config file exists and warns if present
func checkLegacyNetworkConfigFile() {
	configDir, err := configDirPath()
	if err != nil {
		logger.Debugf("Failed to get config directory path while checking for legacy file: %v", err)
		return
	}

	legacyFilePath := path.Join(configDir, legacy_networkConfigFile)

	if _, err := os.Stat(legacyFilePath); err == nil {
		logger.Warnf("Legacy network config file '%s' found at %s", legacy_networkConfigFile, legacyFilePath)
		logger.Warn("This file is not used anymore and can be safely removed")
	}
}

// checkDuplicateNames checks that there are no networks using the same name
func checkDuplicateNames(config *ConfigFile) error {
	if config == nil || config.Networks == nil {
		return nil
	}

	// Since Networks is a map[string]RPCConfigItem, duplicate keys would be overwritten
	// during YAML unmarshaling. However, we can check for case-insensitive duplicates
	// or provide warnings about potential naming conflicts.

	seenNames := make(map[string]string) // lowercase -> original case

	for networkName := range config.Networks {
		lowerName := strings.ToLower(networkName)
		if existingName, exists := seenNames[lowerName]; exists {
			return fmt.Errorf("duplicate network names detected: '%s' and '%s' (case-insensitive conflict)", existingName, networkName)
		}
		seenNames[lowerName] = networkName
	}

	return nil
}
