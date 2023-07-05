package config

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/massalabs/station/pkg/certificate"
	"github.com/massalabs/station/pkg/certificate/store"
	"github.com/massalabs/station/pkg/logger"
)

const certificateName = "massaStation"

func Check() error {
	return checkCertificate()
}

func checkCertificate() error {
	certCA, err := Load()
	if err != nil {
		// non blocking error
		logger.Warnf("failed to load the CA: %s.", err)
		logger.Warn("Station will only work using http, or you will have to add the CA to your browser manually.")

		return nil
	}

	// check that ca certificate is known by OS
	//nolint:exhaustruct
	_, err = certCA.Verify(x509.VerifyOptions{})
	if err != nil {
		err := AddToOS(certCA)
		if err != nil {
			// non blocking error
			logger.Warnf("failed to add the CA to the operating system: %s.", err)
		}
	}

	if !certCa.IsKnownByNSSDatabases(certificateName) {
		err := certCa.AddToNSSDatabases(certificateName)
		if err != nil {
			// non blocking error
			logger.Logger.Warnf("failed to add the CA to NSS: %s.", err)
		}
	}

	return nil
}

const (
	CertificateFile = "rootCA.pem"
	KeyFile         = "rootCA-key.pem"
)

// Load loads the CA certificate and private key from the default filesystem locations.
func Load() (*x509.Certificate, error) {
	caRootPath, err := mkcertCARootPath()
	if err != nil {
		return nil, err
	}

	certPath := filepath.Join(caRootPath, CertificateFile)

	cert, err := certificate.LoadCertificate(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the CA certificate file (%s): %w", certPath, err)
	}

	return cert, nil
}

// AddToOS adds the CA to the operating system.
func AddToOS(cert *x509.Certificate) error {
	err := store.Add(cert)
	if err != nil {
		return fmt.Errorf("failed to add the CA to the operating system: %w", err)
	}

	return nil
}

// Get the mkcert CA root path depending on the OS.
func mkcertCARootPath() (string, error) {
	if env := os.Getenv("CAROOT"); env != "" {
		return env, nil
	}

	var dir string

	switch {
	case runtime.GOOS == "windows":
		dir = os.Getenv("LocalAppData")

	case os.Getenv("XDG_DATA_HOME") != "":
		dir = os.Getenv("XDG_DATA_HOME")

	case runtime.GOOS == "darwin" && os.Getenv("HOME") != "":
		dir = os.Getenv("HOME")
		dir = filepath.Join(dir, "Library", "Application Support")

	case runtime.GOOS == "linux" && os.Getenv("HOME") != "":
		dir = os.Getenv("HOME")
		dir = filepath.Join(dir, ".local", "share")
	default:
		msg := "automatic Certificate Authority detection is not supported by your operating system. "
		msg += "Use the CAROOT environment variable to specify its location."

		return "", errors.New(msg)
	}

	return filepath.Join(dir, "mkcert"), nil
}
