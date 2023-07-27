package config

import (
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/massalabs/station/int/configuration"
	"github.com/massalabs/station/pkg/certificate"
	"github.com/massalabs/station/pkg/certificate/store"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/nss"
)

// failureConsequence is the consequence of a failure to add the CA to NSS.
const failureConsequence = "Station will only work using http, or you will have to add the CA to your browser manually."

// Check performs a check on the configuration.
func Check() error {
	var resultErr error

	caRootPath, err := configuration.CAPath()
	if err != nil {
		return caNonBlockingError("failed to get CA path", err)
	}

	certPath := filepath.Join(caRootPath, configuration.CertificateAuthorityFileName)

	err = checkCertificate(certPath)
	if err != nil {
		resultErr = caNonBlockingError("Error while checking certificate", err)
	}

	err = checkNSS(certPath)
	if err != nil {
		resultErr = caNonBlockingError("Error while checking NSS", err)
	}

	return resultErr
}

// nonBlockingError logs a non blocking error. It always returns a `nil` error to be used in `return`.
func nonBlockingError(context string, consequence string, err error) error {
	if err != nil {
		logger.Warnf("%s: %s.", context, err)
		logger.Warn(consequence)
	}

	return nil
}

// caNonBlockingError logs a non blocking error related to the CA loading failure.
func caNonBlockingError(context string, err error) error {
	return nonBlockingError(context, failureConsequence, err)
}

// checkCertificate checks the certificate configuration.
func checkCertificate(certPath string) error {
	certCA, err := certificate.LoadCertificate(certPath)
	if err != nil {
		return fmt.Errorf("failed to load the CA: %w", err)
	}

	// disable linting as we don't care about checking specific attributes
	//nolint:exhaustruct
	_, err = certCA.Verify(x509.VerifyOptions{})
	if err != nil {
		logger.Debug("the CA is not known by the operating system.")

		err := store.Add(certCA)
		if err != nil {
			return fmt.Errorf("failed to add the CA to the operating system: %w", err)
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
		return fmt.Errorf("failed to instantiate the certutil runner: %w", err)
	}

	service, err := nss.NewCertUtilService(runner)
	if err != nil {
		return fmt.Errorf("failed to instantiate the certutil service: %w", err)
	}

	loggerInstance := &NSSManagerLogger{}
	manager := nss.NewManager([]string{}, service, loggerInstance)

	err = manager.AppendDefaultNSSDatabasePaths()
	if err != nil {
		return fmt.Errorf("failed to append default NSS database paths: %w", err)
	}

	if !manager.HasCA(configuration.CertificateAuthorityName) {
		logger.Debug("the CA is not known by at least one local NSS database.")

		err := manager.AddCA(configuration.CertificateAuthorityName, certPath)
		if err != nil {
			//nolint:errcheck
			caNonBlockingError("failed to add the CA to NSS", err)
		} else {
			logger.Debug("the CA was added to NSS.")
		}
	} else {
		logger.Debug("the CA is known by all local NSS databases.")
	}

	return nil
}
