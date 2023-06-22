package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Network    string
	NodeURL    string
	DNSAddress string
}

type NetworkConfig struct {
	DNS  string   `yaml:"dns"`
	URLs []string `yaml:"urls"`
}


const (
	testnetNodeURL = "https://test.massa.net/api/v2"
	// testnet20.2.
	testnetDNSAddress = "AS12YMz7NjyP3aeEWcSsiC58Hba8UxHapfGv7i4PmNMS2eKfmaqqC"

	labnetNodeURL    = "https://labnet.massa.net/api/v2"
	labnetDNSAddress = "AS1PV17jWkbUs7mfXsn8Xfs9AK6tHiJoxuGu7RySFMV8GYdMeUSh"

	buildnetNodeURL    = "https://buildernet.massa.net/api/v2"
	buildnetDNSAddress = "AS12aGWkBorEM2EpKeNyigSkoCqwdxm872g5KkzUiox3v4VosFW3F"

	MassaStationURL = "station.massa"
)

func GetNetwork(network string) string {
	//nolint:goconst
	if network == "TESTNET" || network == "LABNET" || network == "BUILDNET" {
		return network
	}

	return "UNKNOWN"
}
type Config struct {
	Networks map[string]NetworkConfig `yaml:"network"`
}




// LoadConfig reads the YAML configuration file and returns a map of network configurations.
// The keys of the map represent the network names, and the values contain the corresponding network configuration.
// If the configuration file is successfully loaded and parsed, the map of network configurations is returned along with nil error.
// If there is an error reading the file or parsing the YAML data, the function returns nil map and the encountered error.
func LoadConfig(filename string) (map[string]NetworkConfig, error) {
	var config map[string]NetworkConfig

	// Read the configuration file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Error reading configuration file:", err)
		return nil, err
	}

	// Parse the YAML data
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("Error parsing YAML data:", err)
		return nil, err
	}
	testvalue :=map[string]NetworkConfig{
		"testnet": {
		  DNS: "AS12YMz7NjyP3aeEWcSsiC58Hba8UxHapfGv7i4PmNMS2eKfmaqqC",
		  URLs: []string{"https://test.massa.net/api/v2"},
		},
		"buildnet": {
		  DNS: "AS12aGWkBorEM2EpKeNyigSkoCqwdxm872g5KkzUiox3v4VosFW3F",
		  URLs: []string{"https://buildernet.massa.net/api/v2"},
		},
		"labnet": {
		  DNS: "AS1PV17jWkbUs7mfXsn8Xfs9AK6tHiJoxuGu7RySFMV8GYdMeUSh",
		  URLs: []string{"192.168.2.1", "192.168.2.2"},
		},
	  }
	  

	log.Println("Configuration loaded successfully:", config)

	return testvalue, nil
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

	switch urlOrNetwork {
	case "TESTNET":
		return testnetDNSAddress

	case "LABNET":
		return labnetDNSAddress

	case "BUILDNET":
		return buildnetDNSAddress

	case "LOCALHOST":
		return ""

	default:
		return ""
	}
}
