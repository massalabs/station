package store

import (
	"os"
	"path/filepath"
)

// NSSDBPaths returns all the known NSS databases directories of mac operating system.
func NSSDBPaths() ([]string, error) {
	var nssDBPaths []string

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	nssDBPaths = []string{
		filepath.Join(home, "/Library/Application Support/Firefox/Profiles/*"),
	}

	return nssDBPaths, nil
}
