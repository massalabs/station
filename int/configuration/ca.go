package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Constants for CA-related files and names.
const (
	CertificateAuthorityName = "massaStation"
)

const errMsgCAPath = "CAPath detection is not supported, use the CAROOT environment variable to specify its location"

// CAPath returns the default Certificate Authority storage directory for the current OS.
// It checks the "CAROOT" environment variable first. If "CAROOT" is not set, it will
// use OS-specific default paths.
func CAPath() (string, error) {
	// If the CAROOT environment variable is set, use that
	if caroot := os.Getenv("CAROOT"); caroot != "" {
		return caroot, nil
	}

	var baseDir string

	switch {
	case runtime.GOOS == "windows":
		baseDir = os.Getenv("LocalAppData")

	// Mostly common for Unix-like operating systems
	case os.Getenv("XDG_DATA_HOME") != "":
		baseDir = os.Getenv("XDG_DATA_HOME")

	case runtime.GOOS == "darwin" && os.Getenv("HOME") != "":
		homeDir := os.Getenv("HOME")
		baseDir = filepath.Join(homeDir, "Library", "Application Support")

	case runtime.GOOS == "linux" && os.Getenv("HOME") != "":
		homeDir := os.Getenv("HOME")
		baseDir = filepath.Join(homeDir, ".local", "share")

	default:
		return "", fmt.Errorf(errMsgCAPath)
	}

	// Append "mkcert" to the base directory and return
	return filepath.Join(baseDir, "mkcert"), nil
}
