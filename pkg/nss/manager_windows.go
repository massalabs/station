package nss

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getFirefoxProfilePath returns the path to the Firefox default profile directory.
// It will search for the first directory that ends with ".default-release" which is the directory that contains the NSS database.
func getFirefoxProfilePath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("failed to get APPDATA environment variable")
	}

	firefoxProfilesDir := filepath.Join(appData, "Mozilla", "Firefox", "Profiles")
	_, err := os.Stat(firefoxProfilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("failed to find Firefox profiles directory: %w", err)
		}
		return "", fmt.Errorf("failed to check Firefox profiles directory: %w", err)
	}

	profiles, err := os.ReadDir(firefoxProfilesDir)
	if err != nil {
		return "", fmt.Errorf("failed to read Firefox profiles directory: %w", err)
	}

	var defaultProfileDir string
	for _, profile := range profiles {
		if profile.IsDir() && strings.HasSuffix(profile.Name(), ".default-release") {
			defaultProfileDir = profile.Name()
			break
		}
	}

	if defaultProfileDir == "" {
		return "", fmt.Errorf("failed to find default Firefox profile directory")
	}

	nssDBPath := filepath.Join(firefoxProfilesDir, defaultProfileDir)

	return nssDBPath, nil
}

// defaultNSSDatabasePaths returns the known NSS databases directories of a Windows operating system.
func defaultNSSDatabasePaths() ([]string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return nil, fmt.Errorf("failed to get APPDATA environment variable")
	}

	nssDBPaths := []string{}

	firefoxDBPath, err := getFirefoxProfilePath()
	if err == nil {
		nssDBPaths = append(nssDBPaths, firefoxDBPath)
	}

	return nssDBPaths, nil
}
