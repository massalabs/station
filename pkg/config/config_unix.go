//go:build unix

package config

import (
	"fmt"
	"os"
	"path"
)

const (
	directoryPermissions = 0o755
)

// getConfigDir returns the config directory for the current OS.
// On Unix, it is placed in /usr/local/share/massa.
func getConfigDir() (string, error) {
	path := path.Join("/", "usr", "local", "share", "massastation")

	_, err := os.Stat(path)
	if err != nil {
		// Try to create the directory
		err = os.MkdirAll(path, directoryPermissions)
		if err != nil {
			return "", fmt.Errorf("creating config directory: %w", err)
		}
	}

	return path, nil
}

// getCertDir returns the cert directory for the current OS.
// On Unix, it is placed in /etc/massa/certs.
func getCertDir() (string, error) {
	path := path.Join("/", "etc", "massastation", "certs")

	return path, nil
}
