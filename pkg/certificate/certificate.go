package certificate

// This package aims to solve the impossibility of generating a valid certificate for an entire top-level domain (TLD).
//
// The following approaches are not possible:
// - using a *.massa certificate
// - adding several SANs to the same certificate (e.g. *.my.massa, *.thyra.massa, *.web.massa...) due to
//    * the technical limit to the number of SANs that can be included in a certificate as well as
//    * the delay between adding a website to the blockchain and adding it to the certificate.
//
// Therefore, here we use the SNI mechanism to:
// - get the server name and
// - generate a temporary certificate instead of just retrieving an existing one.
//
// To pass the security checks of the browser and the operating system, we also sign this certificate
// with a root certificate generated with mkcert.

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
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

const privateKeySizeInBits = 2048

// wrapAndPrint wraps error and print error message to std err.
// Printing to stderr is mandatory as this is not done by the http server later in the process.
func wrapAndPrint(msg string, err error) error {
	wrappingMsg := "while handling SNI Hello request"

	err = fmt.Errorf("%s: %w", msg, err)
	fmt.Fprintf(os.Stderr, "%s: %v", wrappingMsg, err)

	return fmt.Errorf("%s: %w", wrappingMsg, err)
}

// GenerateTLS processes an SNI Hello request by generating a dynamic certificate.
func GenerateTLS(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeySizeInBits)
	if err != nil {
		return nil, wrapAndPrint("while generating keypair", err)
	}

	caPath, err := mkcertCARootPath()
	if err != nil {
		return nil, wrapAndPrint("while getting mkcert CA root path", err)
	}

	caCertificate, err := LoadCertificate(caPath + "/rootCA.pem")
	if err != nil {
		return nil, wrapAndPrint("while loading CA certificate", err)
	}

	caPrivateKey, err := LoadPrivateKey(caPath + "/rootCA-key.pem")
	if err != nil {
		return nil, wrapAndPrint("while loading CA key", err)
	}

	certBytes, privateKey, err := GenerateSignedCertificate(hello.ServerName, privateKey, caCertificate, caPrivateKey)
	if err != nil {
		return nil, wrapAndPrint("while generating signed certificate", err)
	}

	var cert tls.Certificate

	cert.Certificate = append(cert.Certificate, certBytes)
	cert.PrivateKey = privateKey

	return &cert, nil
}

// GenerateSignedCertificate generates a certificate and signed it with the given AC.
func GenerateSignedCertificate(
	serverName string,
	privateKey *rsa.PrivateKey,
	caCertificate *x509.Certificate, caPrivateKey crypto.PrivateKey,
) ([]byte, *rsa.PrivateKey, error) {
	if len(serverName) == 0 {
		return nil, nil, errors.New("while generating certificate: server name is empty")
	}

	//nolint:exhaustruct
	template := &x509.Certificate{
		SerialNumber: &big.Int{}, // mandatory, but useless as it should be a unique id given to the certificate by the CA
		Subject: pkix.Name{ // not necessary, but cool to have if the user ask for certificate details.
			CommonName:   serverName,
			Organization: []string{"thyra dynamically generated"},
		},
		NotBefore: time.Now(),
		// one day of validity is enough since the certificate is generated dynamically each time
		NotAfter: time.Now().AddDate(0, 0, 1),
	}

	template.DNSNames = append(template.DNSNames, serverName)

	// Create the certificate.
	cert, err := x509.CreateCertificate(rand.Reader, template, caCertificate, privateKey.Public(), caPrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("while creating certificate: %w", err)
	}

	return cert, privateKey, nil
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
