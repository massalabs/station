package config

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/massalabs/station/pkg/logger"
	"gopkg.in/yaml.v2"
)

const (
	MassaStationURL   = "station.massa"
	configFile        = "config.yaml"
	mainnetRPC        = "https://mainnet.massa.net/api/v2"
	mainnetChainID    = 77658377
	buildnetRPC       = "https://buildnet.massa.net/api/v2"
	buildnetChainID   = 77658366
	permissionUrwGrOr = 0o644
	configDirName     = "massa-station"
)

// MSConfigManager represents a manager for network configurations.
type MSConfigManager struct {
	Network     NetworkConfig
	configMutex sync.RWMutex
}

// NetworkManager is an alias for backward compatibility
type NetworkManager = MSConfigManager

// NewNetworkManager creates a new NetworkManager (for backward compatibility)
func NewNetworkManager() (*NetworkManager, error) {
	return GetConfigManager()
}

// NewMSConfigManager creates a new MSConfigManager (for backward compatibility)
func NewMSConfigManager() (*MSConfigManager, error) {
	return GetConfigManager()
}

var (
	instance      *MSConfigManager
	instanceOnce  sync.Once
	instanceMutex sync.RWMutex
)

// GetConfigManager returns the singleton instance of MSConfigManager
func GetConfigManager() (*MSConfigManager, error) {
	var err error
	instanceOnce.Do(func() {
		instance, err = newMSConfigManager()
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// ResetConfigManager resets the singleton instance (useful for testing)
func ResetConfigManager() {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	instance = nil
	instanceOnce = sync.Once{}
}

// newMSConfigManager creates a new instance (internal function)
func newMSConfigManager() (*MSConfigManager, error) {
	logger.Info("Loading Config...")

	// Load network configuration
	userConfig, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	var defaultNetwork RPCInfos
	networks := make([]RPCInfos, 0, len(userConfig.Networks))

	for networkName, networkConfig := range userConfig.Networks {

		networkStatus := NetworkStatusUp
		version, chainID, err := fetchRPCInfos(networkConfig.URL)
		if err != nil {
			logger.Warnf("failed to fetch RPC status for network %s: %v", networkName, err)
			networkStatus = NetworkStatusDown
		}
		network := RPCInfos{
			Name:    networkName,
			NodeURL: networkConfig.URL,
			Version: version,
			ChainID: chainID,
			status:  networkStatus,
		}
		networks = append(networks, network)

		if networkConfig.Default != nil && *networkConfig.Default {
			defaultNetwork = network
		}
	}

	//nolint: exhaustruct
	manager := &MSConfigManager{
		Network: NetworkConfig{
			currentNetwork: &defaultNetwork,
			Networks:       networks,
		},
	}

	// Start the background refresh routine
	StartNetworkRefresh(manager)

	err = manager.SwitchNetwork(defaultNetwork.Name)
	if err != nil {
		return nil, fmt.Errorf("set default network %s: %w", defaultNetwork.Name, err)
	}

	return manager, nil
}

func (n *MSConfigManager) CurrentNetwork() *RPCInfos {
	return n.Network.currentNetwork
}

// Networks returns a slice of available network names (for backward compatibility)
func (n *MSConfigManager) Networks() *[]string {
	n.configMutex.RLock()
	defer n.configMutex.RUnlock()

	options := make([]string, 0, len(n.Network.Networks))
	for _, network := range n.Network.Networks {
		options = append(options, network.Name)
	}
	return &options
}

// SwitchNetwork switches the current network configuration to the specified network.
// rpcName: The name of the network configuration to switch to.
// Returns any error encountered during the switch operation.
func (n *MSConfigManager) SwitchNetwork(rpcName string) error {
	n.configMutex.Lock()
	defer n.configMutex.Unlock()

	// Find the network with the specified name
	var targetNetwork *RPCInfos
	for i := range n.Network.Networks {
		if n.Network.Networks[i].Name == rpcName {
			targetNetwork = &n.Network.Networks[i]
			break
		}
	}

	if targetNetwork == nil {
		return fmt.Errorf("unknown network: %s", rpcName)
	}

	// Set as current network (background routine will keep it updated)
	n.Network.currentNetwork = targetNetwork

	logger.Debugf("Switched to network: %s", rpcName)

	return nil
}

func configDirPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config directory: %w", err)
	}

	path := path.Join(configDir, configDirName)

	// create the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("creating account directory '%s': %w", path, err)
		}
	}

	return path, nil
}

func LoadConfig() (*ConfigFile, error) {
	configDir, err := configDirPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config dir: %w", err)
	}

	// Check for legacy config file and warn if present
	checkLegacyNetworkConfigFile()

	filePath := path.Join(configDir, configFile)

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		writeDefaultConfig(filePath)
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var configData ConfigFile

	// Unmarshal the YAML data into the configData variable
	err = yaml.Unmarshal(yamlFile, &configData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	// Check for duplicate network names
	err = checkDuplicateNames(&configData)
	if err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &configData, nil
}

// getConfigFilePath returns the full path to the user's config file
func getConfigFilePath() (string, error) {
	configDir, err := configDirPath()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	filePath := path.Join(configDir, configFile)
	return filePath, nil
}

// saveConfig writes the provided configuration to disk
func saveConfig(cfg *ConfigFile) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	ymlFile, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, ymlFile, permissionUrwGrOr); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddNetwork adds a new network to both memory and persistent configuration
func (n *MSConfigManager) AddNetwork(name, url string, makeDefault bool) error {
	n.configMutex.Lock()
	defer n.configMutex.Unlock()

	if name == "" || url == "" {
		return fmt.Errorf("name and url are required")
	}

	// Load current persisted configuration
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if cfg.Networks == nil {
		cfg.Networks = map[string]RPCConfigItem{}
	}

	if _, exists := cfg.Networks[name]; exists {
		return fmt.Errorf("network already exists: %s", name)
	}

	// If default requested, clear defaults on others
	if makeDefault {
		for k, v := range cfg.Networks {
			v.Default = nil
			cfg.Networks[k] = v
		}
	}

	cfg.Networks[name] = RPCConfigItem{URL: url}
	if makeDefault {
		cfg.Networks[name] = RPCConfigItem{URL: url, Default: boolPtr(true)}
	}

	if err := saveConfig(cfg); err != nil {
		return err
	}

	// Update in-memory state
	version, chainID, fetchErr := fetchRPCInfos(url)
	status := NetworkStatusUp
	if fetchErr != nil {
		logger.Warnf("failed to fetch RPC status for network %s: %v", name, fetchErr)
		status = NetworkStatusDown
	}
	newNet := RPCInfos{
		Name:    name,
		NodeURL: url,
		Version: version,
		ChainID: chainID,
		status:  status,
	}
	n.Network.Networks = append(n.Network.Networks, newNet)

	if makeDefault {
		n.Network.currentNetwork = &n.Network.Networks[len(n.Network.Networks)-1]
	}

	return nil
}

// EditNetwork edits an existing network. If newName is provided, the network is renamed.
// If makeDefault is provided and true, this network becomes the default in the persisted configuration
// and the current network is switched in memory as well.
func (n *MSConfigManager) EditNetwork(currentName string, newURL *string, makeDefault *bool, newName *string) error {
	n.configMutex.Lock()
	defer n.configMutex.Unlock()

	if currentName == "" {
		return fmt.Errorf("currentName is required")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if cfg.Networks == nil {
		return fmt.Errorf("no networks configured")
	}

	item, exists := cfg.Networks[currentName]
	if !exists {
		return fmt.Errorf("unknown network: %s", currentName)
	}

	targetName := currentName
	if newName != nil && *newName != "" && *newName != currentName {
		if _, nameTaken := cfg.Networks[*newName]; nameTaken {
			return fmt.Errorf("a network with the new name already exists: %s", *newName)
		}
		targetName = *newName
	}

	// Update URL if requested
	if newURL != nil && *newURL != "" {
		item.URL = *newURL
	}

	// Handle default flag
	if makeDefault != nil {
		if *makeDefault {
			// Clear all defaults
			for k, v := range cfg.Networks {
				v.Default = nil
				cfg.Networks[k] = v
			}
			item.Default = boolPtr(true)
		} else {
			item.Default = nil
		}
	}

	// Apply rename if needed
	if targetName != currentName {
		delete(cfg.Networks, currentName)
		cfg.Networks[targetName] = item
	} else {
		cfg.Networks[currentName] = item
	}

	if err := saveConfig(cfg); err != nil {
		return err
	}

	// Update in-memory slice
	targetIdx := -1
	for i := range n.Network.Networks {
		if n.Network.Networks[i].Name == currentName {
			targetIdx = i
			break
		}
	}
	if targetIdx == -1 {
		// If it's not loaded yet (unlikely), append it
		version := n.Network.currentNetwork.Version
		chainID := n.Network.currentNetwork.ChainID
		url := item.URL
		if newURL != nil && *newURL != "" {
			url = *newURL
		}
		n.Network.Networks = append(n.Network.Networks, RPCInfos{
			Name:    targetName,
			NodeURL: url,
			Version: version,
			ChainID: chainID,
			status:  NetworkStatusUp,
		})
		targetIdx = len(n.Network.Networks) - 1
	} else {
		// Update fields
		n.Network.Networks[targetIdx].Name = targetName
		if newURL != nil && *newURL != "" {
			n.Network.Networks[targetIdx].NodeURL = *newURL
			version, chainID, fetchErr := fetchRPCInfos(*newURL)
			if fetchErr != nil {
				logger.Warnf("failed to refresh edited network %s: %v", targetName, fetchErr)
				n.Network.Networks[targetIdx].status = NetworkStatusDown
			} else {
				n.Network.Networks[targetIdx].Version = version
				n.Network.Networks[targetIdx].ChainID = chainID
				n.Network.Networks[targetIdx].status = NetworkStatusUp
			}
		}
	}

	// Switch current network if default requested or if current was renamed
	if makeDefault != nil && *makeDefault {
		n.Network.currentNetwork = &n.Network.Networks[targetIdx]
	} else if n.Network.currentNetwork != nil && n.Network.currentNetwork.Name == currentName && targetName != currentName {
		n.Network.currentNetwork = &n.Network.Networks[targetIdx]
	}

	return nil
}

// DeleteNetwork removes a network from both memory and persistent configuration
func (n *MSConfigManager) DeleteNetwork(name string) error {
	n.configMutex.Lock()
	defer n.configMutex.Unlock()

	if name == "" {
		return fmt.Errorf("name is required")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if cfg.Networks == nil {
		return fmt.Errorf("no networks configured")
	}
	if _, exists := cfg.Networks[name]; !exists {
		return fmt.Errorf("unknown network: %s", name)
	}
	if len(cfg.Networks) == 1 {
		return fmt.Errorf("cannot delete the last remaining network")
	}

	if len(n.Network.Networks) <= 1 {
		return fmt.Errorf("cannot delete the last remaining network")
	}

	deletingCurrent := n.Network.currentNetwork != nil && n.Network.currentNetwork.Name == name

	deletingDefault := false
	if item, ok := cfg.Networks[name]; ok && item.Default != nil && *item.Default {
		deletingDefault = true
	}

	// Delete from persisted config
	delete(cfg.Networks, name)

	// If we deleted the default network, pick an arbitrary remaining item as new default
	if deletingDefault {
		for k, v := range cfg.Networks {
			v.Default = boolPtr(true)
			cfg.Networks[k] = v
			break
		}
	}

	if err := saveConfig(cfg); err != nil {
		return err
	}

	// Remove from in-memory slice
	idx := -1
	for i := range n.Network.Networks {
		if n.Network.Networks[i].Name == name {
			idx = i
			break
		}
	}
	if idx != -1 {
		n.Network.Networks = append(n.Network.Networks[:idx], n.Network.Networks[idx+1:]...)
	}

	// If current was deleted, switch to the first remaining network
	if deletingCurrent {
		if len(n.Network.Networks) == 0 {
			return fmt.Errorf("no remaining networks after deletion")
		}
		n.Network.currentNetwork = &n.Network.Networks[0]
	}

	return nil
}
