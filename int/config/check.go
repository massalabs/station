package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/massalabs/station/int/configuration"
	"github.com/massalabs/station/int/sni"
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

	// Check certificates.
	caRootPath, err := configuration.CertPath()
	if err != nil {
		return caNonBlockingError("failed to get CA path", err)
	}

	certPath := filepath.Join(caRootPath, configuration.CertificateAuthorityFileName)

	// Ensure CA certificate exists and is valid
	isBrandNewCA, err := ensureValidCA(certPath, caRootPath)
	if err != nil {
		return caNonBlockingError("failed to ensure valid CA", err)
	}

	// Check certificate in system store
	err = checkCertificate(certPath, caRootPath)
	if err != nil {
		resultErr = caNonBlockingError("Error while checking certificate", err)
	}

	// Check NSS configuration
	err = checkNSS(certPath, isBrandNewCA)
	if err != nil {
		resultErr = caNonBlockingError("Error while checking NSS", err)
	}

	return resultErr
}

// ensureValidCA ensures that a valid CA certificate exists at the specified path.
// It returns true if a new CA was generated, false if existing CA is valid.
func ensureValidCA(certPath, caRootPath string) (bool, error) {
	// Check if certificate file exists
	_, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		logger.Infof("Certificate file does not exist, generating new CA...")
		return true, generateCA(caRootPath, "certificate file does not exist")
	}

	// Certificate file exists, validate it
	existingCert, loadErr := certificate.LoadCertificate(certPath)
	if loadErr != nil {
		logger.Warnf("Failed to load existing certificate: %v, generating new one", loadErr)
		return true, generateCA(caRootPath, "failed to load existing certificate")
	}

	// Check if certificate is expired
	if existingCert.NotAfter.Before(time.Now()) {
		logger.Warnf("Certificate is expired (NotAfter=%s), generating new one", existingCert.NotAfter.String())
		return true, generateCA(caRootPath, "certificate is expired")
	}

	// Certificate is valid
	logger.Debugf("Certificate is valid until %s", existingCert.NotAfter.String())
	return false, nil
}

// generateCA generates a new CA certificate with consistent parameters and logging.
func generateCA(caRootPath, reason string) error {
	logger.Infof("Generating new CA certificate (reason: %s)...", reason)

	err := certificate.GenerateCA(
		configuration.OrganizationName,
		configuration.CertificateAuthorityKeyFileName,
		configuration.CertificateAuthorityFileName,
		caRootPath,
	)
	if err != nil {
		return fmt.Errorf("failed to generate CA certificate: %w", err)
	}

	certPath := filepath.Join(caRootPath, configuration.CertificateAuthorityFileName)
	logger.Infof("New CA certificate generated successfully at %s", certPath)

	return nil
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
func checkCertificate(certPath, caRootPath string) error {
	// Load certificate
	certCA, err := loadCertificate(certPath)
	if err != nil {
		return err
	}

	// Check if certificate is trusted by the operating system
	return ensureCertificateInSystemStore(certCA, caRootPath)
}

// loadCertificate loads a certificate.
func loadCertificate(certPath string) (*x509.Certificate, error) {
	certCA, err := certificate.LoadCertificate(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load the CA: %w", err)
	}

	return certCA, nil
}

// ensureCertificateInSystemStore ensures the certificate is trusted by the operating system.
func ensureCertificateInSystemStore(certCA *x509.Certificate, caRootPath string) error {
	err := verifyServedCertificate(caRootPath)
	if err != nil {
		logger.Debugf("Certificate is not trusted by the operating system (%s), adding to store", err)
		return addCertificateToSystemStore(certCA)
	}

	logger.Debug("Certificate is already trusted by the operating system")
	return nil
}

// verifyServedCertificate verifies, as a browser would, a certificate issued by the CA for MassaStation.
// Verifying the CA certificate itself is not enough: it succeeds when the CA is trusted but not present in
// a keychain, while clients only receive the leaf certificate and need to find the CA to build the chain.
func verifyServedCertificate(caRootPath string) error {
	//nolint:exhaustruct // Only the server name is needed to generate the certificate.
	tlsCert, err := sni.GenerateTLS(&tls.ClientHelloInfo{ServerName: MassaStationURL}, caRootPath)
	if err != nil {
		return fmt.Errorf("failed to generate a certificate for %s: %w", MassaStationURL, err)
	}

	leaf, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse the certificate for %s: %w", MassaStationURL, err)
	}

	//nolint:exhaustruct // We don't care about checking specific attributes
	_, err = leaf.Verify(x509.VerifyOptions{DNSName: MassaStationURL})
	if err != nil {
		return fmt.Errorf("failed to verify the certificate for %s: %w", MassaStationURL, err)
	}

	return nil
}

// addCertificateToSystemStore adds a certificate to the operating system store.
func addCertificateToSystemStore(certCA *x509.Certificate) error {
	err := store.Add(certCA)
	if err != nil {
		logger.Errorf("Failed to add certificate to operating system: %v", err)
		return fmt.Errorf("failed to add the CA to the operating system: %w", err)
	}

	logger.Debug("Certificate successfully added to the operating system store")
	return nil
}

type NSSManagerLogger struct{}

func (m *NSSManagerLogger) Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

func (m *NSSManagerLogger) Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// checkNSS checks the NSS configuration.
func checkNSS(certPath string, isBrandNewCA bool) error {
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

	if isBrandNewCA {
		// If we have a brand new CA, we need to delete it from NSS
		// because a CA may already exist with the same name but
		// with a different certificate.
		err = manager.DeleteCA(configuration.CertificateAuthorityName)
		if err != nil {
			logger.Debug("failed to delete the CA from NSS", err)
		} else {
			logger.Debug("the CA was deleted from NSS.")
		}
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
