package certificate

// This file provides a collection of operations to handle PEM (Privacy Enhanced Mail) encoded entities.
// PEM, commonly used for encoding TLS keys and certificates, employs a base64 encoding scheme,
// surrounded by header and footer lines.
// For a comprehensive understanding of the PEM format, refer to RFC 1421.
// For an in-depth examination of Textual Encodings of Public-Key Infrastructure X.509 (PKIX),
// Public-Key Cryptography Standards (PKCS), and Cryptographic Message Syntax (CMS) Structures, see RFC 7468.
//
// The package provides:
//  - The PemType enumeration, encapsulating various recognized PEM types.
//  - The DecodePEM function, used for decoding arbitrary PEM data without a specific type requirement.
//  - The DecodeExpectedPEM function, designed to decode PEM data of a predetermined type.
//
// Design considerations:
//  - This package adheres to exporting sentinel errors best practice.
//
// Future enhancements:
//  Prior to transitioning this package into a standalone GitHub repository, please consider:
//  - Adding encoding functionality.
//  - Confirming whether similar requirements exist in other applications, validating the broad utility of this package.

import (
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	// ErrInvalidPemType is returned when the provided PemType is not valid.
	ErrInvalidPemType = errors.New("invalid PEM type")

	// ErrInvalidData is returned when the provided data cannot be decoded.
	ErrInvalidData = errors.New("invalid data")

	// ErrDataIsNotExpectedPemType is returned when the decoded PEM data is not of the expected type.
	ErrDataIsNotExpectedPemType = errors.New("data is not matching expected PEM type")
)

// PemType represents the PEM types that can be used with the DecodePEM function.
// Use the Certificate and PrivateKey constants to specify the PEM type.
type PemType uint8

const (
	Certificate PemType = iota + 1
	CertificateRequest
	X509CRL
	PrivateKey
)

// IsValid checks if the PemType value is a valid one.
// It returns true for Certificate, CertificateRequest, X509CRL or PrivateKey,
// and false for any other value.
func (p PemType) IsValid() bool {
	return p >= Certificate && p <= PrivateKey
}

// String returns the string representation of the PemType.
// It returns "CERTIFICATE" for Certificate, "PRIVATE KEY" for PrivateKey,
// "CERTIFICATE REQUEST" for certificateRequest, "X509 CRL" for X509CRL,
// and an empty string for any other value.
func (p PemType) String() string {
	switch p {
	case Certificate:
		return "CERTIFICATE"
	case CertificateRequest:
		return "CERTIFICATE REQUEST"
	case X509CRL:
		return "X509 CRL"
	case PrivateKey:
		return "PRIVATE KEY"
	default:
		return ""
	}
}

// PemTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func NewPemType(s string) (PemType, error) {
	switch s {
	case "CERTIFICATE":
		return Certificate, nil
	case "PRIVATE KEY":
		return PrivateKey, nil
	default:
		return 0, fmt.Errorf("%w: %s", ErrInvalidPemType, s)
	}
}

// DecodePEM decodes the given PEM encoded data without validating its type.
// It returns the decoded data as a pem.Block, or an error if the data cannot be decoded.
func DecodePEM(data []byte) (*pem.Block, error) {
	decodedBlock, _ := pem.Decode(data)
	if decodedBlock == nil {
		return nil, ErrInvalidData
	}

	return decodedBlock, nil
}

// DecodeExpectedPEM decodes the given PEM encoded data and checks that its type is the expected one.
// It returns the decoded data as a pem.Block, or an error if the data cannot be decoded
// or the type of the data does not match the expected one.
func DecodeExpectedPEM(data []byte, expectedType PemType) (*pem.Block, error) {
	if !expectedType.IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPemType, expectedType)
	}

	decodedBlock, err := DecodePEM(data)
	if err != nil {
		return nil, err
	}

	if decodedBlock.Type != expectedType.String() {
		return nil, fmt.Errorf("%w: %s", ErrDataIsNotExpectedPemType, decodedBlock.Type)
	}

	return decodedBlock, nil
}
