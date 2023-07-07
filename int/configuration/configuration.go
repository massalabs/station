package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

//nolint:lll
const errMsgPath = "configuration path detection is not supported, use the MASSA_HOME environment variable to specify its location"

// Path returns the default configuration directory for the current OS.
// It checks the "MASSA_HOME" environment variable first. If "MASSA_HOME" is not set, it will
// use OS-specific default paths.
func Path() (string, error) {
	// If the MASSA_HOME environment variable is set, use that
	if massaHome := os.Getenv("MASSA_HOME"); massaHome != "" {
		return massaHome, nil
	}

	switch {
	case runtime.GOOS == "windows":
		ex, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("getting executable path: %w", err)
		}

		return filepath.Dir(ex), nil

	case runtime.GOOS == "linux", runtime.GOOS == "darwin":
		return path.Join("/", "usr", "local", "share", "massastation"), nil

	default:
		return "", fmt.Errorf(errMsgPath)
	}
}
