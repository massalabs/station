package store

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

var ErrCertificateNotFound = fmt.Errorf("certificate not found")

func CAUniqueFilename(cert *x509.Certificate) string {
	CAUniqueName := strings.ReplaceAll("MassaLabs CA "+cert.SerialNumber.String(), " ", "_")

	return fmt.Sprintf("%s.cert", CAUniqueName)
}

// TODO: I would like to put this function in pkg/certificate/certificate.go but there is an import cycle.
// WriteCertificate writes a certificate to a file.
func WriteCertificate(path string, cert *x509.Certificate) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("while creating file: %w", err)
	}
	defer file.Close()

	// Create a PEM block using the certificate
	certBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw, // Assuming "cert" is the *x509.Certificate
	})

	// Write the PEM block to the file
	_, err = file.Write(certBytes)
	if err != nil {
		return fmt.Errorf("while writing certificate to file: %w", err)
	}

	return nil
}
