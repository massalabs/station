package config

import (
	"embed"
	"fmt"
	"regexp"
	"sync"

	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
	"gopkg.in/yaml.v2"
)

const (
	MassaStationURL = "station.massa"
)

type NetworkInfos struct {
	Network string
	NodeURL string
	Version string
	ChainID uint64
}

// NetworkConfig represents the configuration of a network.
//
//nolint:tagliatelle
type NetworkConfig struct {
	URLs    []string `yaml:"URLs"`
	Default bool     `yaml:"Default"`
	ChainID uint64   `yaml:"ChainID"`
}

// NetworkManager represents a manager for network configurations.
type NetworkManager struct {
	networkInfos  NetworkInfos             // Current network configuration
	knownNetworks map[string]NetworkConfig // Known network configurations
	networkMutex  sync.RWMutex             // Mutex for concurrent access to network configuration
}

// NewNetworkManager creates a new instance of NetworkManager.
// It loads the initial network configurations from the specified file and sets the default network configuration.
// Returns the initialized NetworkManager and any error encountered during initialization.
func NewNetworkManager() (*NetworkManager, error) {
	//nolint: exhaustruct
	networkManager := &NetworkManager{
		networkInfos:  NetworkInfos{},
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

	err = networkManager.SwitchNetwork(defaultNetwork)
	if err != nil {
		return nil, fmt.Errorf("set default network %s: %w", defaultNetwork, err)
	}

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

// SetCurrentNetwork sets the current network configuration for the NetworkManager.
func (n *NetworkManager) SetCurrentNetwork(config NetworkInfos) {
	n.networkMutex.Lock()
	defer n.networkMutex.Unlock()

	n.networkInfos = config
}

// Network returns the current network configuration.
// It returns the network configuration represented by an NetworkInfos struct.
func (n *NetworkManager) Network() *NetworkInfos {
	return &n.networkInfos
}

func versionFromStatus(status *node.State) (string, error) {
	pattern := `.+\.(\d+\.\d+)`

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(*status.Version)

	//nolint:gomnd
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to parse node version from: %s", *status.Version)
	}

	return matches[1], nil
}

// SwitchNetwork switches the current network configuration to the specified network.
// selectedNetworkStr: The name of the network configuration to switch to.
// Returns any error encountered during the switch operation.
func (n *NetworkManager) SwitchNetwork(selectedNetworkStr string) error {
	cfg, ok := n.knownNetworks[selectedNetworkStr]
	if !ok {
		return fmt.Errorf("unknown network option: %s", selectedNetworkStr)
	}

	url := cfg.URLs[0]
	client := node.NewClient(url)

	status, err := node.Status(client)
	if err != nil {
		return fmt.Errorf("failed to retrieve node status: %w", err)
	}

	version, err := versionFromStatus(status)
	if err != nil {
		return fmt.Errorf("getting network version: %w", err)
	}

	// compare chain id from node status with chain id from config
	nodeChainID := uint64(*status.ChainID)
	if nodeChainID != cfg.ChainID {
		logger.Errorf("chain id mismatch: %d != %d", nodeChainID, cfg.ChainID)

		return fmt.Errorf("chain id mismatch: %d != %d", nodeChainID, cfg.ChainID)
	}

	n.SetCurrentNetwork(NetworkInfos{
		NodeURL: url,
		Network: selectedNetworkStr,
		Version: version,
		ChainID: cfg.ChainID,
	})

	logger.Debugf("Set current network: %s (version %s)\n", selectedNetworkStr, version)
	logger.Debugf("Network config: %+v\n", n.Network())

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
