package runner

import (
	"fmt"
	"os/exec"
)

type CommandRunner struct {
	BinaryPath string
}

// Run runs the command with the given arguments.
// It returns the combined output of stdout and stderr in the error, if any.
func (r *CommandRunner) Run(args ...string) error {
	cmd := exec.Command(r.BinaryPath, args...) // #nosec G204

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}

	return nil
}
