package config

import (
	"crypto/x509"
	"path/filepath"

	"github.com/massalabs/station/int/configuration"
	"github.com/massalabs/station/pkg/certificate"
	"github.com/massalabs/station/pkg/certificate/store"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/nss"
)

// caFailureConsequence is the consequence of a failure to load the CA.
//
//nolint:lll
const caFailureConsequence = "Station will only work using http, or you will have to add the CA to your browser manually (OS store)."

// nssFailureConsequence is the consequence of a failure to add the CA to NSS.
//
//nolint:lll
const nssFailureConsequence = "Station will only work using http, or you will have to add the CA to your browser manually (NSS store)."

// Check performs a check on the configuration.
func Check() error {
	caRootPath, err := configuration.CAPath()
	if err != nil {
		return caCheckNonBlockingError("failed to get CA path", err)
	}

	certPath := filepath.Join(caRootPath, configuration.CertificateAuthorityFileName)

	keyPath := filepath.Join(caRootPath, configuration.CertificateAuthorityKeyFileName)

	err = checkCertificate(certPath, keyPath)
	if err != nil {
		return err
	}

	err = checkNSS(certPath)
	if err != nil {
		return err
	}

	return nil
}

// nonBlockingError logs a non blocking error.
func nonBlockingError(context string, consequence string, err error) error {
	if err != nil {
		logger.Warnf("%s: %s.", context, err)
		logger.Warn(consequence)
	}

	return nil
}

// caCheckNonBlockingError logs a non blocking error related to the CA loading failure.
func caCheckNonBlockingError(context string, err error) error {
	return nonBlockingError(context, caFailureConsequence, err)
}

// nssCheckNonBlockingError logs a non blocking error related to the CA loading failure.
func nssCheckNonBlockingError(context string, err error) error {
	return nonBlockingError(context, nssFailureConsequence, err)
}

// checkCertificate checks the certificate configuration.
func checkCertificate(certPath string, keyPath string) error {
	certCA, err := certificate.LoadCertificate(filepath.Join(certPath, keyPath))
	if err != nil {
		return caCheckNonBlockingError("failed to load the CA", err)
	}

	// disable linting as we don't care about checking specific attributes
	//nolint:exhaustruct
	_, err = certCA.Verify(x509.VerifyOptions{})
	if err != nil {
		logger.Debug("the CA is not known by the operating system.")

		err := store.Add(certCA)
		if err != nil {
			return caCheckNonBlockingError("failed to add the CA to the operating system", err)
		}

		logger.Debug("the CA was added to the operating system.")
	} else {
		logger.Debug("the CA is known by the operating system.")
	}

	return nil
}

type NSSManagerLogger struct{}

func (m *NSSManagerLogger) Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args)
}

func (m *NSSManagerLogger) Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args)
}

// checkNSS checks the NSS configuration.
func checkNSS(certPath string) error {
	runner, err := nss.NewCertUtilRunner()
	if err != nil {
		return nssCheckNonBlockingError("failed to instantiate the certutil runner", err)
	}

	service, err := nss.NewCertUtilService(runner)
	if err != nil {
		return nssCheckNonBlockingError("failed to instantiate the certutil service", err)
	}

	loggerInstance := &NSSManagerLogger{}
	manager := nss.NewManager([]string{}, service, loggerInstance)

	if !manager.HasCA(configuration.CertificateAuthorityName) {
		logger.Debug("the CA is not known by at least one local NSS database.")

		err := manager.AddCA(configuration.CertificateAuthorityName, certPath)
		if err != nil {
			// non blocking error
			logger.Warnf("failed to add the CA to NSS: %s.", err)
			logger.Warn(caFailureConsequence)
		} else {
			logger.Debug("the CA was added to NSS.")
		}
	} else {
		logger.Debug("the CA is known by all local NSS databases.")
	}

	return nil
}
