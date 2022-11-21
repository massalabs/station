package config

import (
	"errors"
	"os"
	"path"
)

func GetConfigDir() (string, error) {
	confDir := path.Join(os.Getenv("HOME"), ".config", "thyra")
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		return "", errors.New("Unable to find config dir: " + confDir)
	}

	return confDir, nil
}
