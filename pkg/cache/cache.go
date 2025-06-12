package cache

// The cache mechanism use the filesystem to read and store file in the websitesCache
// inside the Config dir.

import (
	"fmt"
	"os"
	"path"
)

type Cache struct {
	ConfigDir string
}

func fullPath(fileName, configDir string) (string, error) {
	cacheDir, err := fsDirectory(configDir)
	if err != nil {
		return "", fmt.Errorf("while reading cached file %s: %w", fileName, err)
	}

	return path.Join(cacheDir, fileName), nil
}

// IsPresent checks if the file is present in the local cache.
func (c *Cache) IsPresent(file string) bool {
	fp, _ := fullPath(file, c.ConfigDir)

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
func (c *Cache) Read(file string) ([]byte, error) {
	fullPath, err := fullPath(file, c.ConfigDir)
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
func (c *Cache) Save(fileName string, content []byte) error {
	fullPath, err := fullPath(fileName, c.ConfigDir)
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

// Delete removes a file from the cache.
func (c *Cache) Delete(fileName string) error {
	fullPath, err := fullPath(fileName, c.ConfigDir)
	if err != nil {
		return fmt.Errorf("while reading cached file %s: %w", fileName, err)
	}

	err = os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("while deleting from cache: %w", err)
	}

	return nil
}
