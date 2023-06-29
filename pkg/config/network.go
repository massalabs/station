package config

import (
	"embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v2"
)

const (
	MassaStationURL = "station.massa"
)

type AppConfig struct {
	Network    string
	NodeURL    string
	DNSAddress string
}

// NetworkConfig represents the configuration of a network.
//
//nolint:tagliatelle
type NetworkConfig struct {
	DNS     string   `yaml:"DNS"`
	URLs    []string `yaml:"URLs"`
	Default bool     `yaml:"Default"`
}

// NetworkManager represents a manager for network configurations.
type NetworkManager struct {
	appConfig     AppConfig                // Current network configuration
	knownNetworks map[string]NetworkConfig // Known network configurations
	networkMutex  sync.RWMutex             // Mutex for concurrent access to network configuration
}

// NewNetworkManager creates a new instance of NetworkManager.
// It loads the initial network configurations from the specified file and sets the default network configuration.
// Returns the initialized NetworkManager and any error encountered during initialization.
func NewNetworkManager() (*NetworkManager, error) {
	//nolint: exhaustruct
	networkManager := &NetworkManager{
		appConfig:     AppConfig{},
		knownNetworks: make(map[string]NetworkConfig),
	}

	// Load network configuration
	initNetworksData, err := LoadConfig()
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

// SetNetworks sets the known networks for the NetworkManager.
func (n *NetworkManager) SetNetworks(networks map[string]NetworkConfig) {
	n.knownNetworks = networks
}

// Networks retrieves a pointer to a slice of known networks from the NetworkManager.
// It returns a pointer to a slice containing the names of the known networks.
// The slice will be updated if the known networks are modified.
func (n *NetworkManager) Networks() *[]string {
	options := make([]string, 0, len(n.knownNetworks))

	for network := range n.knownNetworks {
		options = append(options, network)
	}

	return &options
}

// NetworkFromString retrieves the network configuration corresponding to the given network name.
// It returns the network configuration represented by an AppConfig struct.
// An error is returned if the network configuration is not found or if the provided network name is invalid.
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

// SetNetwork sets the current network configuration for the NetworkManager.
func (n *NetworkManager) SetNetwork(config AppConfig) {
	n.networkMutex.Lock()
	defer n.networkMutex.Unlock()

	n.appConfig = config
}

// Network returns the current network configuration.
// It returns the network configuration represented by an AppConfig struct.
func (n *NetworkManager) Network() *AppConfig {
	return &n.appConfig
}

// SwitchNetwork switches the current network configuration to the specified network.
// selectedNetworkStr: The name of the network configuration to switch to.
// Returns any error encountered during the switch operation.
func (n *NetworkManager) SwitchNetwork(selectedNetworkStr string) error {
	cfg, ok := n.knownNetworks[selectedNetworkStr]
	if !ok {
		return fmt.Errorf("unknown network option: %s", selectedNetworkStr)
	}

	n.SetNetwork(AppConfig{
		NodeURL:    cfg.URLs[0],
		DNSAddress: cfg.DNS,
		Network:    selectedNetworkStr,
	})

	Logger.Debugf("Switched to network: %s\n", selectedNetworkStr)
	Logger.Debugf("Current config: %+v\n", n.Network())

	return nil
}

//nolint:typecheck,nolintlint
//go:embed config_network.yaml
var configData embed.FS

// LoadConfig loads network configurations from an embedded YAML file.
// Returns the loaded network configurations and any error encountered during loading.
func LoadConfig() (map[string]NetworkConfig, error) {
	var networksData map[string]NetworkConfig

	// Read the embedded YAML file
	yamlFile, err := configData.ReadFile("config_network.yaml")
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
