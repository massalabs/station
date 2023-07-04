package config

import "github.com/massalabs/station/pkg/certificate"

func Check() error {
	return checkCertificate()
}

func checkCertificate() error {
	ca := certificate.CA{}
	ca.Load()

	if !ca.IsKnownByOS() {
		err := ca.AddToOS()
		if err != nil {
			//non blocking error
			Logger.Warnf("failed to add the CA to the operating system: %s", err)
		}
	}

	return nil
}
