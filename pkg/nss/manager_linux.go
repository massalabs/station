package nss

import (
	"fmt"
	"os"
	"path/filepath"
)

// defaultNSSDatabasePaths returns the known NSS databases directories of a Linux operating system.
func defaultNSSDatabasePaths() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	nssDBPaths := []string{
		"/etc/pki/nssdb/",
		filepath.Join(home, ".pki/nssdb/"),
		filepath.Join(home, "snap/chromium/current/.pki/nssdb/"),
		filepath.Join(home, ".mozilla/firefox/*"),
		filepath.Join(home, "snap/firefox/common/.mozilla/firefox/*"),
	}

	return nssDBPaths, nil
}
