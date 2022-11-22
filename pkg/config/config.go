package config

import (
	"errors"
	"os"
	"path"
)

func GetConfigDir() (string, error) {
	confDir := path.Join(os.Getenv("HOME"), ".config", "thyra")

	_, err := os.Stat(confDir)
	if err != nil {
		return "", errors.New("Unable to read config dir: " + confDir)
	}

	return confDir, nil
}
