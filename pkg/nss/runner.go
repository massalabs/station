package nss

import (
	"github.com/massalabs/station/pkg/runner"
)

// This file provides a CertUtilRunner struct which encapsulates certutil commands. Certutil is a command-line utility
// that can create and manage certificate and key database files, and can also create many kinds of certificates.
//
// The Run method provides a unified way to execute any certutil command with provided arguments.
//
// Future enhancements:
//  More functionalities can be added to the CertUtilRunner struct if required, based on the usage of the certutil tool.
//  Instead of using the certutil tool, the NSS databases can be managed directly using the NSS library.

var _ runner.Runner = &CertUtilRunner{}

// CertUtilRunner encapsulates certutil commands.
type CertUtilRunner struct {
	runner.CommandRunner
}
