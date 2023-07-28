package nss

import (
	"fmt"
	"os"
	"path/filepath"
)

// defaultNSSDatabasePaths returns the known NSS databases directories of a Windows operating system.
func defaultNSSDatabasePaths() ([]string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return nil, fmt.Errorf("failed to get APPDATA environment variable")
	}

	nssDBPaths := []string{
		filepath.Join(appData, "Mozilla", "Firefox", "Profiles", "*"),
	}

	return nssDBPaths, nil
}
