package config

import "github.com/massalabs/station/pkg/certificate"

func Check() error {
	return checkCertificate()
}

func checkCertificate() error {
	certCa := certificate.CA{}

	err := certCa.Load()
	if err != nil {
		// non blocking error
		Logger.Warnf("failed to load the CA: %s.", err)
		Logger.Warn("Station will only work using http, or you will have to add the CA to your browser manually.")
	}

	if !certCa.IsKnownByOS() {
		err := certCa.AddToOS()
		if err != nil {
			// non blocking error
			Logger.Warnf("failed to add the CA to the operating system: %s.", err)
		}
	}

	return nil
}
