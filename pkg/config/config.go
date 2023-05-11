package config

import (
	"errors"
	"os"
	"path"
)

//nolint:gochecknoglobals
var Version string

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Unable to get user home dir: " + err.Error())
	}

	confDir := path.Join(homeDir, ".config", "thyra")

	_, err = os.Stat(confDir)
	if err != nil {
		return "", errors.New("Unable to read config dir: " + confDir + ": " + err.Error())
	}

	return confDir, nil
}
