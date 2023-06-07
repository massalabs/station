package config

type AppConfig struct {
	Network    string
	NodeURL    string
	DNSAddress string
}

const (
	testnetNodeURL = "https://test.massa.net/api/v2"
	// testnet20.2.
	testnetDNSAddress = "AS12YMz7NjyP3aeEWcSsiC58Hba8UxHapfGv7i4PmNMS2eKfmaqqC"

	labnetNodeURL    = "https://labnet.massa.net/api/v2"
	labnetDNSAddress = "AS1PV17jWkbUs7mfXsn8Xfs9AK6tHiJoxuGu7RySFMV8GYdMeUSh"

	buildnetNodeURL    = "https://buildernet.massa.net/api/v2"
	buildnetDNSAddress = "AS1HqqZF5nFiZEzC7A19t7pUBRtvNfAq5c5PsESSPRE9eFYrGxhg"

	MassaStationURL = "station.massa"
)

func GetNetwork(network string) string {
	//nolint:goconst
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