package config

import (
	"fmt"
	"os"
)

//nolint:gochecknoglobals
var Version = "dev"

type NetworkManager struct {
	appConfig     AppConfig
	knownNetworks map[string]NetworkConfig
}

func NewNetworkManager(configFile string) (*NetworkManager, error) {
	networkManager := &NetworkManager{
		appConfig: AppConfig{
			Network:    "",
			NodeURL:    "",
			DNSAddress: "",
		},
		knownNetworks: make(map[string]NetworkConfig),
	}

	// Load network configuration from file
	initNetworksData, err := LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	networkManager.SetNetworks(initNetworksData)

	// Get AppConfig for the selected network (BuildNet)
	initConfig, err := networkManager.GetAppConfig(BuildNet)
	if err != nil {
		return nil, fmt.Errorf("failed to get app configuration: %w", err)
	}

	networkManager.SetAppConfig(initConfig)

	return networkManager, nil
}

func (n *NetworkManager) Network() AppConfig {
	return n.appConfig
}

func (n *NetworkManager) SetAppConfig(config AppConfig) {
	n.appConfig = config
}

func (n *NetworkManager) SetNetworks(networks map[string]NetworkConfig) {
	n.knownNetworks = networks
}

func (n *NetworkManager) Networks() map[string]NetworkConfig {
	return n.knownNetworks
}

func (n *NetworkManager) GetNetworkOptions() []NetworkOption {
	options := make([]NetworkOption, 0, len(n.knownNetworks))

	for network := range n.knownNetworks {
		switch network {
		case "testnet":
			options = append(options, TestNet)
		case "buildnet":
			options = append(options, BuildNet)
		case "labnet":
			options = append(options, LabNet)
		}
	}

	return options
}

func (n *NetworkManager) GetAppConfig(selectedNetwork NetworkOption) (AppConfig, error) {
	// Convert the NetworkOption to string for lookup
	selectedNetworkStr := selectedNetwork.String()

	config, ok := n.knownNetworks[selectedNetworkStr]
	if !ok {
		return AppConfig{}, fmt.Errorf("selected network '%s' not found", selectedNetworkStr)
	}

	appConfig := AppConfig{
		// we will Implement later the logic to choose the appropriate URL based on the selected network configuration.
		NodeURL:    config.URLs[0],
		DNSAddress: config.DNS,
		Network:    selectedNetworkStr,
	}

	return appConfig, nil
}

func (n *NetworkManager) SwitchNetwork(selectedNetwork NetworkOption) error {
	// Get AppConfig for the selected network
	newConfig, err := n.GetAppConfig(selectedNetwork)
	if err != nil {
		return fmt.Errorf("failed to get app configuration: %w", err)
	}

	// Set the new configuration
	n.SetAppConfig(newConfig)

	return nil
}

// GetConfigDir returns the config directory for the current OS.
func GetConfigDir() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(confDir)
	if err != nil {
		return "", fmt.Errorf("unable to read config directory: %s: %w", confDir, err)
	}

	return confDir, nil
}

// GetCertDir returns the cert directory for the current OS.
func GetCertDir() (string, error) {
	certDir, err := getCertDir()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(certDir)
	if err != nil {
		return "", fmt.Errorf("unable to read cert directory: %s: %w", certDir, err)
	}

	return certDir, nil
}
