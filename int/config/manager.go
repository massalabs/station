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
