package sni

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
	privateKeySizeInBits = 2048
	loggerPrefix         = "SNI -"
)

// hashServerName returns a unique identifier for the server, hashing the server's name with the current time stamp.
// The function returns an error if the server name is empty.
func hashServerName(serverName string) ([]byte, error) {
	if len(serverName) == 0 {
		return nil, fmt.Errorf("%w: server name is empty", ErrInvalidArgument)
	}

	now, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal current time: %w", err)
	}

	return sha256.New().Sum(append([]byte(serverName), now...)), nil
}

// createCertificateTemplate builds a x509 certificate template using the server name and a unique ID.
// DNSNames is set to the given server name, the serial number is set to the unique site name id,
// NotBefore is set to the current time, NotAfter is set to 1 day after the current time,
// and the organization is set to "station dynamically generated".
//
// The function returns an error if the server name or the unique site name id are empty.
func createCertificateTemplate(serverName string, uniqueSiteNameID []byte) (*x509.Certificate, error) {
	if len(serverName) == 0 {
		return nil, fmt.Errorf("%w: server name is empty", ErrInvalidArgument)
	}

	if len(uniqueSiteNameID) == 0 {
		return nil, fmt.Errorf("%w: unique site name id is empty", ErrInvalidArgument)
	}

	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetBytes(uniqueSiteNameID),
		Subject: pkix.Name{
			CommonName:   serverName,
			Organization: []string{"station dynamically generated"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
		DNSNames:  []string{serverName},
	}

	return template, nil
}

// generateSignedCertificate creates a certificate and then signs it using the provided Certificate Authority (CA).
// It uses root certificate and private key from the default location.
// It uses hashServerName and createCertificateTemplate to ensure uniqueness and proper formatting of the certificate.
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

	uniqueSiteNameID, err := hashServerName(serverName)
	if err != nil {
		return nil, nil, err
	}

	template, err := createCertificateTemplate(serverName, uniqueSiteNameID)
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

	certBytes, privateKey, err := generateSignedCertificate(caPath, hello.ServerName)
	if err != nil {
		logger.Errorf("%s generate signed certificate for %s failed: %w", loggerPrefix, hello.ServerName, err)

		return nil, fmt.Errorf("%s generate signed certificate for %s failed: %w", loggerPrefix, hello.ServerName, err)
	}

	return &tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privateKey,
	}, nil
}
