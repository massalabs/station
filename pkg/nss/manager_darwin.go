package nss

import (
	"fmt"
	"os"
	"path/filepath"
)

// nssDBUsualPath returns the known NSS databases directories of a Darwin operating system.
func nssDBUsualPath() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	nssDBPaths := []string{
		filepath.Join(home, "/Library/Application Support/Firefox/Profiles/*"),
	}

	return nssDBPaths, nil
}
