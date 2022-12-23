package certificate

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	rootName    = "rootCA.pem"
	rootKeyName = "rootCA-key.pem"
)

func GenerateSignedCertificate(serverName string, priv *rsa.PrivateKey) ([]byte, *rsa.PrivateKey, error) {
	if len(serverName) == 0 {
		return nil, nil, errors.New("while generating certificate: server name is empty")
	}

	caPath, err := mkcertCARootPath()
	if err != nil {
		return nil, nil, fmt.Errorf("while generating certificate: %w", err)
	}

	caCert, err := LoadCertificate(caPath + "/" + rootName)
	if err != nil {
		return nil, nil, fmt.Errorf("while loading CA certificate: %w", err)
	}

	caPrivateKey, err := LoadPrivateKey(caPath + "/" + rootKeyName)
	if err != nil {
		return nil, nil, fmt.Errorf("while loading CA key: %w", err)
	}

	//nolint:exhaustruct
	tpl := &x509.Certificate{
		SerialNumber: &big.Int{}, // mandatory, but useless as it should be a unique id given to the certificate by the CA
		Subject: pkix.Name{ // not necessary, but cool to have if the user ask for certificate details.
			CommonName:   serverName,
			Organization: []string{"thyra dynamically generated"},
		},
		NotBefore: time.Now(),
		// one day of validity is enough since the certificate is generated dynamically each time
		NotAfter: time.Now().AddDate(0, 0, 1),
	}

	tpl.DNSNames = append(tpl.DNSNames, serverName)

	// Create the certificate.
	cert, err := x509.CreateCertificate(rand.Reader, tpl, caCert, priv.Public(), caPrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("while creating certificate: %w", err)
	}

	return cert, priv, nil
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

// Loads a PEM encoded certificate.
func LoadCertificate(path string) (*x509.Certificate, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("while loading certificate file from %s: %w", path, err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("while reading certificate file from %s: %w", path, err)
	}

	der, _ := pem.Decode(raw) // rest argument, additional information if any, is ignored
	if der == nil || der.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("file %s is not a certificate", path)
	}

	cert, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("while parsing DER bytes: %w", err)
	}

	return cert, nil
}

// Loads a PEM encoded private key.
func LoadPrivateKey(path string) (crypto.PrivateKey, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("while loading key file from %s: %w", path, err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("while reading key file from %s: %w", path, err)
	}

	der, _ := pem.Decode(raw)
	if der == nil || der.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("file %s is not a private key", path)
	}

	key, err := x509.ParsePKCS8PrivateKey(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("while parsing DER bytes: %w", err)
	}

	return key, nil
}
