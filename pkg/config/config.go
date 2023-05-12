package config

import (
	"fmt"
	"os"
)

//nolint:gochecknoglobals
var Version = "dev"

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
