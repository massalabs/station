//go:build linux
// +build linux

package store

import (
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/massalabs/station/pkg/su"
)

// inspired by: https://github.com/FiloSottile/mkcert/blob/master/main.go

// For Debian based systems.
const (
	OSCertificateDirectory = "/usr/local/share/ca-certificates"
	OSCertificateCommand   = "update-ca-certificates"
)

func AddToOS(cert *x509.Certificate) error {
	err := WriteCertificate(OSCertificateFilename(cert), cert)
	if err != nil {
		return fmt.Errorf("failed to write the CA certificate to the filesystem: %w", err)
	}

	command, err := su.Command(OSCertificateCommand)
	if err != nil {
		return fmt.Errorf("failed to create the command to update the CA certificates: %w", err)
	}

	err = command.Run()
	if err != nil {
		return fmt.Errorf("failed to update the CA certificates: %w", err)
	}

	return nil
}

func DeleteFromOS(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}

func OSCertificateFilename(cert *x509.Certificate) string {
	return filepath.Join(OSCertificateDirectory, CAUniqueFilename(cert))
}
