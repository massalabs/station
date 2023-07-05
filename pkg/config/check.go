package config

import (
	"github.com/massalabs/station/pkg/certificate"
	"github.com/massalabs/station/pkg/logger"
)

func Check() error {
	return checkCertificate()
}

func checkCertificate() error {
	certCa := certificate.CA{}

	err := certCa.Load()
	if err != nil {
		// non blocking error
		logger.Logger.Warnf("failed to load the CA: %s.", err)
		logger.Logger.Warn("Station will only work using http, or you will have to add the CA to your browser manually.")

		return nil
	}

	if !certCa.IsKnownByOS() {
		err := certCa.AddToOS()
		if err != nil {
			// non blocking error
			logger.Logger.Warnf("failed to add the CA to the operating system: %s.", err)
		}
	}

	return nil
}
