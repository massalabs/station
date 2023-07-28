package nss

import (
	"fmt"
	"os/exec"
)

// This file provides a CertUtilRunner struct which encapsulates certutil commands. Certutil is a command-line utility
// that can create and manage certificate and key database files, and can also create many kinds of certificates.
//
// The Run method provides a unified way to execute any certutil command with provided arguments.
//
// Future enhancements:
//  More functionalities can be added to the CertUtilRunner struct if required, based on the usage of the certutil tool.
//  Instead of using the certutil tool, the NSS databases can be managed directly using the NSS library.

var _ Runner = &CertUtilRunner{}

// CertUtilRunner encapsulates certutil commands.
type CertUtilRunner struct {
	binaryPath string
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
