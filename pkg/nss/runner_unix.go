//go:build unix

package nss

import (
	"fmt"
	"os/exec"

	"github.com/massalabs/station/pkg/runner"
)

// NewCertUtilRunner returns a new CertUtilRunner.
// It returns an error if the certutil binary is not found.
func NewCertUtilRunner() (*CertUtilRunner, error) {
	bin, err := exec.LookPath("certutil")
	if err != nil {
		return nil, fmt.Errorf("failed to find certutil binary: %w", err)
	}

	return &CertUtilRunner{runner.CommandRunner{BinaryPath: bin}}, nil
}
