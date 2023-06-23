package config

import (
	"fmt"
	"io/ioutil"

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


// LoadConfig reads the YAML configuration file and returns a map of network configurations.
// The keys of the map represent the network names, and the values contain the corresponding network configuration.
func LoadConfig(filename string) (map[string]NetworkConfig, error) {
	var networksData map[string]NetworkConfig

	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal the YAML data into the configData variable
	err = yaml.Unmarshal(yamlFile, &networksData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	return networksData, nil
}
