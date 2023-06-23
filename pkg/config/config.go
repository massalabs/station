package config

import (
	"fmt"
	"os"
)
var Version = "dev"


type NetworkManager struct {
	appConfig     AppConfig
	knownNetworks map[string]NetworkConfig
}

func NewNetworkManager(configFile string) (*NetworkManager, error) {
	nm := &NetworkManager{}

	// Load network configuration from file
	initNetworksData, err := LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	nm.SetNetworks(initNetworksData)

	// Get AppConfig for the selected network (BuildNet)
	initConfig, err := nm.GetAppConfig(BuildNet)
	if err != nil {
		return nil, fmt.Errorf("failed to get app configuration: %w", err)
	}

	nm.SetAppConfig(initConfig)

	return nm, nil
}

func (nm *NetworkManager) Network() AppConfig {
	return nm.appConfig
}

func (nm *NetworkManager) SetAppConfig(config AppConfig) {
	nm.appConfig = config
}

func (nm *NetworkManager) SetNetworks(networks map[string]NetworkConfig) {
	nm.knownNetworks = networks
}

func (nm *NetworkManager) Networks() map[string]NetworkConfig {
	return nm.knownNetworks
}

func (nm *NetworkManager) GetNetworkOptions() []NetworkOption {
	options := make([]NetworkOption, 0, len(nm.knownNetworks))

	for network := range nm.knownNetworks {
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

func (nm *NetworkManager) GetAppConfig(selectedNetwork NetworkOption) (AppConfig, error) {
	// Convert the NetworkOption to string for lookup
	selectedNetworkStr := selectedNetwork.String()

	config, ok := nm.knownNetworks[selectedNetworkStr]
	if !ok {
		return AppConfig{}, fmt.Errorf("selected network '%s' not found", selectedNetworkStr)
	}

	appConfig := AppConfig{
		// TODO: Add logic to choose which URL to use
		NodeURL:    config.URLs[0],
		DNSAddress: config.DNS,
		Network:    selectedNetworkStr,
	}

	return appConfig, nil
}

func (nm *NetworkManager) SwitchNetwork(selectedNetwork NetworkOption) error {
	// Get AppConfig for the selected network
	newConfig, err := nm.GetAppConfig(selectedNetwork)
	if err != nil {
		return fmt.Errorf("failed to get app configuration: %w", err)
	}

	// Set the new configuration
	nm.SetAppConfig(newConfig)

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
