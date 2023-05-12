package config

import (
	"errors"
	"os"
)

// getConfigDir returns the config directory for the current OS.
// On Windows, it is at the same location as the executable, by default `C:\Program Files(x86)\MassaStation`.
func getConfigDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("getting executable path: %w", err)
	}

	dir := filepath.Dir(ex)

	return dir, nil
}

// getCertDir returns the cert directory for the current OS.
// On Windows, it is at the same location as the executable, by default `C:\Program Files(x86)\MassaStation\certs`.
func getCertDir() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	certDir := path.Join(confDir, "certs")

	return certDir, nil
}
