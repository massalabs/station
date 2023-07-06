package store

import (
	"os"
	"path/filepath"
)

// NSSDBPaths returns all the known NSS databases directories of linux operating system.
func NSSDBPaths() ([]string, error) {
	var nssDBPaths []string
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	nssDBPaths = []string{
		"/etc/pki/nssdb/",
		filepath.Join(home, ".pki/nssdb/"),
		filepath.Join(home, "snap/chromium/current/.pki/nssdb/"),
		filepath.Join(home, ".mozilla/firefox/*"),
		filepath.Join(home, "snap/firefox/common/.mozilla/firefox/*"),
	}

	return nssDBPaths, nil
}
