package config

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

// NetworkManagerInterface defines the interface for managing networks.
type NetworkManagerInterface interface {
	// SetNetworks sets the known networks for the NetworkManager.
	SetNetworks(networks map[string]NetworkConfig)

	// Networks retrieves a pointer to a slice of known networks from the NetworkManager.
	// It returns a pointer to a slice containing the names of the known networks.
	// The slice will be updated if the known networks are modified.
	Networks() *[]string

	// NetworkFromString retrieves the network configuration corresponding to the given network name.
	// It returns the network configuration represented by an AppConfig struct.
	// An error is returned if the network configuration is not found or if the provided network name is invalid.
	NetworkFromString(networkName string) (*AppConfig, error)

	// SetNetwork sets the current network configuration for the NetworkManager.
	SetNetwork(config AppConfig)

	// Network returns the current network configuration.
	// It returns the network configuration represented by an AppConfig struct.
	Network() *AppConfig

	// SwitchNetwork switches the current network configuration to the specified network.
	// selectedNetworkStr: The name of the network configuration to switch to.
	// Returns any error encountered during the switch operation.
	SwitchNetwork(selectedNetworkStr string) error
}
