package config

import (
	"os"

	"github.com/massalabs/station/pkg/logger"
	"gopkg.in/yaml.v2"
)

//nolint:tagliatelle
type RPCConfigItem struct {
	URL     string `yaml:"URL"`
	Default *bool  `yaml:"Default,omitempty"`
}

type ConfigFile struct {
	StationFirstRunSetupDone bool                     `yaml:"StationFirstRunSetupDone"`
	Networks                 map[string]RPCConfigItem `yaml:"Networks"`
}

var DefaultConfig = ConfigFile{
	StationFirstRunSetupDone: false,
	Networks: map[string]RPCConfigItem{
		"mainnet": {
			URL:     "https://mainnet.massa.net/api/v2",
			Default: boolPtr(true),
		},
		"buildnet": {
			URL: "https://buildnet.massa.net/api/v2",
		},
	},
}

// boolPtr returns a pointer to the given bool value
func boolPtr(b bool) *bool {
	return &b
}

func writeDefaultConfig(filePath string) {
	ymlFile, err := yaml.Marshal(DefaultConfig)
	if err != nil {
		logger.Fatalf("failed to marshal default config: %v", err)
	}

	err = os.WriteFile(filePath, ymlFile, permissionUrwGrOr)
	if err != nil {
		logger.Fatalf("failed to write default config to file: %v", err)
	}
}
