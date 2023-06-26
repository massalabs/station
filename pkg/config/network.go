package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

const (
	MassaStationURL = "station.massa"
)

// NetworkManager represents a manager for network configurations.
type NetworkManager struct {
	appConfig     AppConfig                // Current network configuration
	knownNetworks map[string]NetworkConfig // Known network configurations
	networkMutex  sync.RWMutex             // Mutex for concurrent access to network configuration
}

// Verify at compilation time that NetworkManager implements NetworkManagerInterface.
//
//nolint:exhaustruct
var _ NetworkManagerInterface = &NetworkManager{}

// NewNetworkManager creates a new instance of NetworkManager.
// It loads the initial network configurations from the specified file and sets the default network configuration.
// configFile: The path to the YAML configuration file.
// Returns the initialized NetworkManager and any error encountered during initialization.
func NewNetworkManager(configFile string) (*NetworkManager, error) {
	//nolint: exhaustruct
	networkManager := &NetworkManager{
		appConfig:     AppConfig{},
		knownNetworks: make(map[string]NetworkConfig),
	}

	// Load network configuration from file
	initNetworksData, err := LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	networkManager.SetNetworks(initNetworksData)

	var defaultNetwork string

	for networkName, networkConfig := range initNetworksData {
		if networkConfig.Default {
			defaultNetwork = networkName

			break
		}
	}

	if defaultNetwork == "" {
		return nil, fmt.Errorf("default network not found")
	}

	// Get AppConfig for the selected network
	initConfig, ok := networkManager.knownNetworks[defaultNetwork]
	if !ok {
		return nil, fmt.Errorf("selected network '%s' not found", defaultNetwork)
	}

	appConfig := AppConfig{
		NodeURL:    initConfig.URLs[0],
		DNSAddress: initConfig.DNS,
		Network:    defaultNetwork,
	}

	networkManager.SetNetwork(appConfig)

	return networkManager, nil
}

func (n *NetworkManager) SetNetworks(networks map[string]NetworkConfig) {
	n.knownNetworks = networks
}

func (n *NetworkManager) Networks() *[]string {
	options := make([]string, 0, len(n.knownNetworks))

	for network := range n.knownNetworks {
		options = append(options, network)
	}

	return &options
}

func (n *NetworkManager) NetworkFromString(networkName string) (*AppConfig, error) {
	knownNetworks := *n.Networks()

	for _, name := range knownNetworks {
		if name == networkName {
			config, ok := n.knownNetworks[networkName]
			if !ok {
				return nil, fmt.Errorf("failed to find configuration for network '%s'", networkName)
			}

			appConfig := &AppConfig{
				NodeURL:    config.URLs[0],
				DNSAddress: config.DNS,
				Network:    networkName,
			}

			return appConfig, nil
		}
	}

	return nil, fmt.Errorf("invalid network option: '%s'", networkName)
}

func (n *NetworkManager) SetNetwork(config AppConfig) {
	n.networkMutex.Lock()
	defer n.networkMutex.Unlock()

	n.appConfig = config
}

func (n *NetworkManager) Network() *AppConfig {
	n.networkMutex.RLock()
	defer n.networkMutex.RUnlock()

	return &n.appConfig
}

func (n *NetworkManager) SwitchNetwork(selectedNetworkStr string) error {
	config, ok := n.knownNetworks[selectedNetworkStr]
	if !ok {
		return fmt.Errorf("unknown network option: %s", selectedNetworkStr)
	}

	n.SetNetwork(AppConfig{
		NodeURL:    config.URLs[0],
		DNSAddress: config.DNS,
		Network:    selectedNetworkStr,
	})

	log.Printf("Switched to network: %s\n", selectedNetworkStr)
	log.Printf("Current config: %+v\n", n.Network())

	return nil
}

// LoadConfig loads network configurations from a YAML file.
// filename: The path to the YAML configuration file.
// Returns the loaded network configurations and any error encountered during loading.
func LoadConfig(filename string) (map[string]NetworkConfig, error) {
	var networksData map[string]NetworkConfig

	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal the YAML data into the networksData variable
	err = yaml.Unmarshal(yamlFile, &networksData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	return networksData, nil
}
