package config

import (
	"fmt"
	"io/ioutil"
	"log"

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

var networksData map[string]NetworkConfig
var appConfig AppConfig

func init() {
	// Load network configuration from "config_network.yaml"
	var err error
	networksData, err = LoadConfig("config_network.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Get AppConfig for the selected network (BuildNet)
	appConfig, err = GetAppConfig(BuildNet)
	if err != nil {
		log.Fatal("Failed to get app configuration:", err)
	}
}


func Config() AppConfig {
	return appConfig
}

func Networks() map[string]NetworkConfig {
	return networksData
}

func UpdateConfig(selectedNetwork NetworkOption) {
	// Get AppConfig for the selected network
	newConfig, err := GetAppConfig(selectedNetwork)
	if err != nil {
		log.Fatal("Failed to get app configuration:", err)
	}

	// Update appConfig with the new configuration
	appConfig = newConfig
}


type NetworkConfig struct {
	DNS  string   `yaml:"DNS"`
	URLs []string `yaml:"URLs"`
}

type NetworkOption int

const (
	TestNet NetworkOption = iota
	BuildNet
	LabNet
)

var networkOptionNames = [...]string{
	"testnet",
	"buildnet",
	"labnet",
}

func (option NetworkOption) String() string {
	return networkOptionNames[option]
}

func GetNetworkOptions() []NetworkOption {
	options := make([]NetworkOption, 0, len(Networks()))
	for network := range Networks() {
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

func GetAppConfig(selectedNetwork NetworkOption) (AppConfig, error) {
	// Convert the NetworkOption to string for lookup
	selectedNetworkStr := selectedNetwork.String()

	config, ok := Networks()[selectedNetworkStr]
	if !ok {
		return AppConfig{}, fmt.Errorf("selected network '%s' not found", selectedNetworkStr)
	}

	appConfig := AppConfig{
		NodeURL:    config.URLs[0],
		DNSAddress: config.DNS,
		Network:    selectedNetworkStr,
	}

	return appConfig, nil
}

// LoadConfig reads the YAML configuration file and returns a map of network configurations.
// The keys of the map represent the network names, and the values contain the corresponding network configuration.
// If the configuration file is successfully loaded and parsed, the map of network configurations is returned along with nil error.
// If there is an error reading the file or parsing the YAML data, the function returns nil map and the encountered error.
func LoadConfig(filename string) (map[string]NetworkConfig, error) {
	var networksData map[string]NetworkConfig

	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the configData variable
	err = yaml.Unmarshal(yamlFile, &networksData)
	if err != nil {
		return nil, err
	}
	
	return networksData, nil
}

// NodeURL returns a list of available URLs for a given network.
// If the network is found in the configuration, it returns the corresponding URLs.
// If the network is not found, it returns an empty list.
func NodeURL(network string, networkConfig map[string]NetworkConfig) []string {
	if networkConfig, ok := networkConfig[network]; ok {
		return networkConfig.URLs
	}

	// Network not found, return an empty list
	return []string{}
}

// DNSAddress returns the DNS address for a given network.
// If the network is found in the configuration, it returns the corresponding DNS address.
// If the network is not found, it returns an empty string.
func DNSAddress(network string, networkConfig map[string]NetworkConfig) string {
	if networkConfig, ok := networkConfig[network]; ok {
		return networkConfig.DNS
	}
	// Network not found, return an empty string
	return ""
}