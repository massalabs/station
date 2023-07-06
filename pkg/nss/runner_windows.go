package nss

import (
	"fmt"
	"os/exec"
)

// CertUtilRunner encapsulates certutil commands.
type CertUtilRunner struct{}

// NewCertUtilRunner returns a new CertUtilRunner.
// It returns an error if the certutil binary is not found.
func NewCertUtilRunner() (*CertUtilRunner, error) {
	return nil, fmt.Errorf("not implemented")
}

// Run runs the certutil command with the given arguments.
// It returns the combined output of stdout and stderr.
func (r *CertUtilRunner) Run(args ...string) error {
	return fmt.Errorf("not implemented")
}
