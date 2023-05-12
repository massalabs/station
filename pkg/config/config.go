package config

import (
	"errors"
	"os"
)

//nolint:gochecknoglobals
var Version = "dev"

func GetConfigDir() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(confDir)
	if err != nil {
		return "", errors.New("Unable to read config dir: " + confDir + ": " + err.Error())
	}

	return confDir, nil
}

func GetCertDir() (string, error) {
	certDir, err := getCertDir()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(certDir)
	if err != nil {
		return "", errors.New("Unable to read cert dir: " + certDir + ": " + err.Error())
	}

	return certDir, nil
}
