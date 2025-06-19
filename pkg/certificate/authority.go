package certificate

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// randomSerialNumber generates a random serial number for the certificate.
func randomSerialNumber() (*big.Int, error) {
	bitSize := 128
	// Calculate the limit as 2^bitSize.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), uint(bitSize))

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	return serialNumber, nil
}

// generatePrivateKey generates a new RSA private key.
func generatePrivateKey() (crypto.PrivateKey, error) {
	bits := 3072

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	return privateKey, nil
}

// derivePublicKey derives the public key from the private key.
func derivePublicKey(privateKey crypto.PrivateKey) (crypto.PublicKey, error) {
	signer, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("privateKey does not implement crypto.Signer interface")
	}

	publicKey := signer.Public()

	return publicKey, nil
}

// hashPublicKey computes SHA-256 hash of the public key.
func hashPublicKey(publicKeyBytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(publicKeyBytes)

	return hasher.Sum(nil)
}

// createCATemplate creates a template for the Certificate Authority (CA) certificate.
func createCATemplate(
	serialNumber *big.Int,
	publicKey crypto.PublicKey,
	organization string,
) (*x509.Certificate, error) {
	spkiASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}

	_, err = asn1.Unmarshal(spkiASN1, &spki)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %w", err)
	}

	skid := hashPublicKey(spki.SubjectPublicKey.Bytes)

	years := 2

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   organization,
		},
		SubjectKeyId: skid,

		NotAfter:  time.Now().AddDate(years, 0, 0),
		NotBefore: time.Now(),

		KeyUsage: x509.KeyUsageCertSign,

		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	return template, nil
}

// writeCA writes the private key and certificate to the specified path.
func writeCA(
	certificateAuthorityKeyFileName,
	certificateAuthorityFileName,
	path string,
	privateKey crypto.PrivateKey,
	cert []byte,
) error {
	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	permissionUrwGrOr := 0o644

	err = os.WriteFile(filepath.Join(path, certificateAuthorityKeyFileName), pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: privDER}), fs.FileMode(permissionUrwGrOr))
	if err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	err = os.WriteFile(filepath.Join(path, certificateAuthorityFileName), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert}), fs.FileMode(permissionUrwGrOr))
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	return nil
}

// GenerateCA generates a new Certificate Authority (CA) certificate.
func GenerateCA(
	organizationName,
	certificateAuthorityKeyFileName,
	certificateAuthorityFileName,
	path string,
) error {
	serialNumber, err := randomSerialNumber()
	if err != nil {
		return err
	}

	privateKey, err := generatePrivateKey()
	if err != nil {
		return err
	}

	publicKey, err := derivePublicKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to derive public key: %w", err)
	}

	template, err := createCATemplate(serialNumber, publicKey, organizationName)
	if err != nil {
		return fmt.Errorf("failed to create CA template: %w", err)
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	err = writeCA(
		certificateAuthorityKeyFileName,
		certificateAuthorityFileName,
		path,
		privateKey,
		cert,
	)
	if err != nil {
		return fmt.Errorf("failed to write CA: %w", err)
	}

	return nil
}
