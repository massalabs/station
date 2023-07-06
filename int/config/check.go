package config

import (
	"path/filepath"

	"github.com/massalabs/station/pkg/certificate"
	"github.com/massalabs/station/pkg/logger"
)

// caFailureConsequence is the consequence of a failure to load the CA.
const caFailureConsequence = "Station will only work using http, or you will have to add the CA to your browser manually."

// Check performs a check on the configuration.
func Check() error {
	return checkCertificate()
}

// nonBlockingError logs a non blocking error.
func nonBlockingError(context string, consequence string, err error) error {
	if err != nil {
		logger.Logger.Warnf("%s: %s.", context, err)
		logger.Logger.Warn(consequence)
	}

	return nil
}

// caCheckNonBlockingError logs a non blocking error related to the CA loading failure.
func caCheckNonBlockingError(context string, err error) error {
	return nonBlockingError(context, caFailureConsequence, err)
}

// checkCertificate checks the certificate configuration.
func checkCertificate() error {
	certCa, err := certificate.NewCA(CertificateAuthorityName)
	if err != nil {
		return caCheckNonBlockingError("failed to instantiate the CA", err)
	}

	caRootPath, err := CAPath()
	if err != nil {
		return caCheckNonBlockingError("failed to get CA path", err)
	}

	certPath := filepath.Join(caRootPath, CertificateAuthorityFileName)

	keyPath := filepath.Join(caRootPath, CertificateAuthorityKeyFileName)

	err = certCa.Load(certPath, keyPath)
	if err != nil {
		return caCheckNonBlockingError("failed to load the CA", err)
	}

	if !certCa.IsKnownByOS() {
		logger.Logger.Debug("the CA is not known by the operating system.")

		err := certCa.AddToOS()
		if err != nil {
			// non blocking error
			logger.Logger.Warnf("failed to add the CA to the operating system: %s.", err)
			logger.Logger.Warn(caFailureConsequence)
		}
	} else {
		logger.Logger.Debug("the CA is known by the operating system.")
	}

	if !certCa.IsKnownByNSSDatabases() {
		logger.Logger.Debug("the CA is not known by at least one local NSS database.")

		err := certCa.AddToNSSDatabases()
		if err != nil {
			// non blocking error
			logger.Logger.Warnf("failed to add the CA to NSS: %s.", err)
			logger.Logger.Warn(caFailureConsequence)
		}
	} else {
		logger.Logger.Debug("the CA is known by all local NSS databases.")
	}

	return nil
}
