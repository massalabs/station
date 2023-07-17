//go:build unix

package nss

import (
	"fmt"
	"os/exec"
)

// CertUtilRunner encapsulates certutil commands.
type CertUtilRunner struct {
	binaryPath string
}

// NewCertUtilRunner returns a new CertUtilRunner.
// It returns an error if the certutil binary is not found.
func NewCertUtilRunner() (*CertUtilRunner, error) {
	bin, err := exec.LookPath("certutil")
	if err != nil {
		return nil, fmt.Errorf("failed to find certutil binary: %w", err)
	}

	return &CertUtilRunner{binaryPath: bin}, nil
}

// Run runs the certutil command with the given arguments.
// It returns the combined output of stdout and stderr in the error, if any.
func (r *CertUtilRunner) Run(args ...string) error {
	cmd := exec.Command(r.binaryPath, args...) // #nosec G204

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}

	return nil
}
