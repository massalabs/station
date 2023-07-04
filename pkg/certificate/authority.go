package certificate

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/massalabs/station/pkg/certificate/store"
)

const (
	CertificateFile = "rootCA.pem"
	KeyFile         = "rootCA-key.pem"
)

// CA is an certificate authority struct
type CA struct {
	privateKey crypto.PrivateKey
	cert       *x509.Certificate
}

// Load loads the CA certificate and private key from the default filesystem locations.
func (c *CA) Load() error {
	caRootPath, err := mkcertCARootPath()
	if err != nil {
		return err
	}

	certPath := filepath.Join(caRootPath, CertificateFile)

	c.cert, err = LoadCertificate(certPath)
	if err != nil {
		return fmt.Errorf("failed to parse the CA certificate file (%s): %w", certPath, err)
	}

	keyPath := filepath.Join(caRootPath, KeyFile)

	c.privateKey, err = LoadPrivateKey(keyPath)
	if err != nil {
		return fmt.Errorf("failed to parse the CA private key file (%s): %w", keyPath, err)
	}

	return nil
}

// IsKnownByOS checks if the CA is known by the operating system.
func (c *CA) IsKnownByOS() bool {
	_, err := c.cert.Verify(x509.VerifyOptions{})
	return err == nil
}

// AddToOS adds the CA to the operating system.
func (c *CA) AddToOS() error {
	return store.Add(c.cert)
}
