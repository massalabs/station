//go:build darwin

package store

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/massalabs/station/pkg/runner"
)

var _ runner.Runner = &SecurityRunner{}

// SecurityRunner encapsulates security commands.
type SecurityRunner struct {
	runner.CommandRunner
}

// NewSecurityRunner returns a new SecurityRunner.
// It returns an error if the security binary is not found.
func NewSecurityRunner() (*SecurityRunner, error) {
	bin, err := exec.LookPath("security")
	if err != nil {
		return nil, fmt.Errorf("failed to find security binary: %w", err)
	}

	return &SecurityRunner{runner.CommandRunner{BinaryPath: bin}}, nil
}

func Add(cert *x509.Certificate) error {
	security, err := NewSecurityRunner()
	if err != nil {
		return fmt.Errorf("failed to instantiate the certutil runner: %w", err)
	}

	tempFile, err := os.CreateTemp("", "cert")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	permissionUrwGrOr := 0o644

	err = os.WriteFile(tempFile.Name(), pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), fs.FileMode(permissionUrwGrOr))
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	err = security.Run(
		"add-trusted-cert", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", tempFile.Name(),
	)
	if err != nil {
		return fmt.Errorf("failed to add the certificate to the system keychain: %w", err)
	}

	return nil
}

func Delete(_ *x509.Certificate) error {
	return fmt.Errorf("not implemented")
}
