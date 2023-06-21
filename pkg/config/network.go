package config

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	mu            sync.RWMutex
	Network       string
	NodeURL       string
	DNSAddress    string
	SupportedNets map[string]NetworkConfig
}

type NetworkConfig struct {
	DNS  string   `yaml:"dns"`
	URLs []string `yaml:"urls"`
}

func loadConfig(filename string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg map[string]NetworkConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		SupportedNets: make(map[string]NetworkConfig),
	}

	for net, netConfig := range cfg {
		appConfig.SupportedNets[net] = netConfig
	}

	return appConfig, nil
}

var (
	testnetNodeURL    = "https://test.massa.net/api/v2"
	labnetNodeURL     = "https://labnet.massa.net/api/v2"
	buildnetNodeURL   = "https://buildernet.massa.net/api/v2"
	MassaStationURL   = "station.massa"
)

var appConfig *AppConfig // Declare appConfig variable

func InitConfig(filename string) error {
	cfg, err := loadConfig(filename)
	if err != nil {
		return err
	}

	appConfig = cfg
	return nil
}

func GetNetwork(network string) string {
	// nolint:goconst
	if network == "TESTNET" || network == "LABNET" || network == "BUILDNET" {
		return network
	}

	return "UNKNOWN"
}

func GetNodeURL(urlOrNetwork string) string {
	switch urlOrNetwork {
	case "TESTNET":
		return testnetNodeURL

	case "LABNET":
		return labnetNodeURL

	case "BUILDNET":
		return buildnetNodeURL

	case "LOCALHOST":
		return "http://127.0.0.1:33035"

	default:
		return urlOrNetwork
	}
}

func GetDNSAddress(urlOrNetwork string, dnsFlag string) string {
	if dnsFlag != "" {
		return dnsFlag
	}

	appConfig.mu.RLock()
	defer appConfig.mu.RUnlock()

	switch urlOrNetwork {
	case "TESTNET":
		return appConfig.SupportedNets["testnet"].DNS

	case "LABNET":
		return appConfig.SupportedNets["labnet"].DNS

	case "BUILDNET":
		return appConfig.SupportedNets["buildnet"].DNS

	case "LOCALHOST":
		return ""

	default:
		return ""
	}
}
