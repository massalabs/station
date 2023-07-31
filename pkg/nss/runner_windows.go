//go:build windows
// +build windows

package nss

import (
	"fmt"
	"os"
	"path/filepath"
)

// NewCertUtilRunner returns a new CertUtilRunner.
// It returns an error if the certutil binary is not found.
func NewCertUtilRunner() (*CertUtilRunner, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("getting executable path: %w", err)
	}

	executableDirectoryPath := filepath.Dir(executablePath)
	certutilBinaryPath := filepath.Join(executableDirectoryPath, "mar-tools", "certutil.exe")

	_, err = os.Stat(certutilBinaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find certutil binary: %w", err)
	}

	return &CertUtilRunner{binaryPath: certutilBinaryPath}, nil
}
