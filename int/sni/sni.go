package sni

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/massalabs/station/int/configuration"
	"github.com/massalabs/station/pkg/certificate"
	"github.com/massalabs/station/pkg/logger"
)

var ErrInvalidArgument = errors.New("invalid argument")

const (
	privateKeySizeInBits   = 2048
	serialNumberSizeInBits = 128
	loggerPrefix           = "SNI -"
)

// randomSerialNumber returns a random 128-bit serial number.
// RFC 5280 limits serial numbers to 20 octets: longer ones are rejected by some verifiers (e.g. macOS).
func randomSerialNumber() (*big.Int, error) {
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), serialNumberSizeInBits))
	if err != nil {
		return nil, fmt.Errorf("unable to generate serial number: %w", err)
	}

	return serialNumber, nil
}

// createCertificateTemplate builds a x509 certificate template using the server name and a serial number.
// DNSNames is set to the given server name, NotBefore is set to the current time,
// NotAfter is set to 1 day after the current time, the organization is set to "station dynamically generated",
// and the certificate is restricted to TLS server authentication.
//
// The function returns an error if the server name is empty or the serial number is nil.
func createCertificateTemplate(serverName string, serialNumber *big.Int) (*x509.Certificate, error) {
	if len(serverName) == 0 {
		return nil, fmt.Errorf("%w: server name is empty", ErrInvalidArgument)
	}

	if serialNumber == nil {
		return nil, fmt.Errorf("%w: serial number is nil", ErrInvalidArgument)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   serverName,
			Organization: []string{"station dynamically generated"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 0, 1),
		DNSNames:    []string{serverName},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	return template, nil
}

// generateSignedCertificate creates a certificate and then signs it using the provided Certificate Authority (CA).
// It uses root certificate and private key from the default location.
// It uses randomSerialNumber and createCertificateTemplate to ensure uniqueness and proper formatting of the certificate.
func generateSignedCertificate(serverName, caPath string) ([]byte, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySizeInBits)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate private key: %w", err)
	}

	caCertificate, err := certificate.LoadCertificate(filepath.Join(caPath, configuration.CertificateAuthorityFileName))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load CA certificate: %w", err)
	}

	caPrivateKey, err := certificate.LoadPrivateKey(filepath.Join(caPath, configuration.CertificateAuthorityKeyFileName))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to load CA private key: %w", err)
	}

	serialNumber, err := randomSerialNumber()
	if err != nil {
		return nil, nil, err
	}

	template, err := createCertificateTemplate(serverName, serialNumber)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, caCertificate, privateKey.Public(), caPrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create certificate: %w", err)
	}

	return cert, privateKey, nil
}

// GenerateTLS creates a TLS certificate using the server name of the client hello info request.
// Implementation design: it logs information because log configuration at go-swagger level is not available.
func GenerateTLS(hello *tls.ClientHelloInfo, caPath string) (*tls.Certificate, error) {
	logger.Debugf("%s generating TLS certificate", loggerPrefix)

	if hello == nil {
		logger.Errorf("%s client hello info is nil", loggerPrefix)

		return nil, fmt.Errorf("%w: client hello info is nil", ErrInvalidArgument)
	}

	certBytes, privateKey, err := generateSignedCertificate(hello.ServerName, caPath)
	if err != nil {
		msg := fmt.Sprintf("%s generate signed certificate for %s failed: %s", loggerPrefix, hello.ServerName, err.Error())
		logger.Error(msg)

		return nil, errors.New(msg)
	}

	return &tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privateKey,
	}, nil
}
