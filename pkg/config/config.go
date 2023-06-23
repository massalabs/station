package config

import (
	"fmt"
	"log"
	"os"
)

var (
	appConfig     AppConfig
	knownNetworks map[string]NetworkConfig
	Version       = "dev"
)

// init initializes the known networks and the app configuration.
func init() {
	// Load network configuration from "config_network.yaml"
	initNetworksData, err := LoadConfig("config_network.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	SetNetworks(initNetworksData)

	// Get AppConfig for the selected network (BuildNet)
	initConfig, err := GetAppConfig(BuildNet)
	if err != nil {
		log.Fatal("Failed to get app configuration:", err)
	}

	SetAppConfig(initConfig)
}

// Config returns the current app configuration.
func Config() AppConfig {
	return appConfig
}

// SetAppConfig sets the app configuration.
func SetAppConfig(config AppConfig) {
	appConfig = config
}

// SetNetworks sets the known networks configuration.
func SetNetworks(networks map[string]NetworkConfig) {
	knownNetworks = networks
}

// Networks returns the map of network configurations.
func Networks() map[string]NetworkConfig {
	return knownNetworks
}

// SwitchNetworkTO updates the app configuration based on the selected network.
func SwitchNetworkTO(selectedNetwork NetworkOption) {
	// Get AppConfig for the selected network
	newConfig, err := GetAppConfig(selectedNetwork)
	if err != nil {
		log.Fatal("Failed to get app configuration:", err)
	}

	// Set the new configuration
	SetAppConfig(newConfig)
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
