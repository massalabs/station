package cache

// The cache mecanism use the filesystem to read and store file in the websitesCache
// inside the Config dir.

import (
	"fmt"
	"os"
	"path"
)

func fullPath(fileName, configDir string) (string, error) {
	cacheDir, err := fsDirectory(configDir)
	if err != nil {
		return "", fmt.Errorf("while reading cached file %s: %w", fileName, err)
	}

	return path.Join(cacheDir, fileName), nil
}

// IsPresent checks if the file is present in the local cache.
func IsPresent(file, configDir string) bool {
	fp, _ := fullPath(file, configDir)

	_, err := os.Stat(fp)

	return !os.IsNotExist(err)
}

// fsDirectory returns the cache directory on the file system.
// If the directory doesn't exist, it is created before being returned.
func fsDirectory(configDir string) (string, error) {
	cacheDir := path.Join(configDir, "websitesCache")
	_, err := os.Stat(cacheDir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error creating folder: %w", err)
		}
	}

	return cacheDir, nil
}

// Read returns the cached file content corresponding to the given name.
func Read(file, configDir string) ([]byte, error) {
	fullPath, err := fullPath(file, configDir)
	if err != nil {
		return nil, fmt.Errorf("while reading cached file %s: %w", file, err)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("while reading cached file %s: %w", file, err)
	}

	return content, nil
}

// Save adds a new file to the cache.
func Save(fileName string, content []byte, configDir string) error {
	fullPath, err := fullPath(fileName, configDir)
	if err != nil {
		return fmt.Errorf("while reading cached file %s: %w", fileName, err)
	}

	//nolint:gomnd
	err = os.WriteFile(fullPath, content, 0o600)
	if err != nil {
		return fmt.Errorf("while saving file %s to cache: %w", fullPath, err)
	}

	return nil
}
