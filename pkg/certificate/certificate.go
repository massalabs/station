// Package certificate provides operations for handling certificates and private keys used in TLS
// (Transport Layer Security).
// It is especially geared towards dealing with entities encoded in the PEM (Privacy Enhanced Mail) format.
// These operations include loading and parsing PEM-encoded certificates and private keys from file paths.
//
// The package exposes:
//   - The LoadCertificate function that retrieves a PEM-encoded certificate from a specified
//     file path, decodes the PEM content, and parses the certificate.
//   - The LoadPrivateKey function that fetches a PEM-encoded private key from a given
//     file path, decodes the PEM content, and parses the private key.
//
// Design considerations:
//   - This package adheres to exporting sentinel errors best practice.
//
// Future directions:
//   - As the use of this package evolves, consider adding more functionality, such as writing
//     PEM-encoded certificates and private keys to file paths.
//   - Before transitioning this package to a standalone GitHub repository, ascertain the need
//     for similar functionalities in other applications.
package certificate

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

var (

	// ErrFailedToReadFile is returned when the provided file cannot be read.
	ErrFailedToReadFile = errors.New("unable to read file")
	// ErrFailedToDecodeFile is returned when the provided file cannot be decoded.
	ErrFailedToDecodeFile = errors.New("unable to decode file content")
	// ErrFailedToParseContent is returned when the provided content cannot be parsed.
	ErrFailedToParseContent = errors.New("unable to parse content")
)

// readFile reads the content of a file at the given path.
func readFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("%w: %s", ErrFailedToReadFile, "empty path")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w at %s: %w", ErrFailedToReadFile, path, err)
	}

	return raw, nil
}

// LoadCertificate retrieves a PEM encoded certificate from a specified file path.
// It decodes the PEM file and parses the certificate, returning any errors encountered during the process.
func LoadCertificate(path string) (*x509.Certificate, error) {
	raw, err := readFile(path)
	if err != nil {
		return nil, err
	}

	der, err := DecodeExpectedPEM(raw, Certificate)
	if err != nil {
		return nil, fmt.Errorf("%w in file %s: %w", ErrFailedToDecodeFile, path, err)
	}

	cert, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w in file %s: %w", ErrFailedToParseContent, path, err)
	}

	return cert, nil
}

// LoadPrivateKey retrieves a PEM encoded private key from a specified file path.
// It decodes the PEM file and parses the private key, returning any errors encountered during the process.
func LoadPrivateKey(path string) (crypto.PrivateKey, error) {
	raw, err := readFile(path)
	if err != nil {
		return nil, err
	}

	der, err := DecodeExpectedPEM(raw, PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("%w in file %s: %w", ErrFailedToDecodeFile, path, err)
	}

	key, err := x509.ParsePKCS8PrivateKey(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w in file %s: %w", ErrFailedToParseContent, path, err)
	}

	return key, nil
}
