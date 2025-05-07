//go:build unix

package plugin

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/massalabs/station/pkg/logger"
)

// sendStopSignal sends an appropriate signal to gracefully stop a process on Unix systems.
func (p *Plugin) sendStopSignal() error {
	err := p.command.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to send SIGTERM to plugin %s: %w", p.ID, err)
	}

	logger.Infof("Sent SIGTERM to plugin %s.\n", p.ID)

	return nil
}

// setupProcess configures platform-specific process attributes (no-op on Unix).
func setupProcess(_ *exec.Cmd) {
	// No special setup needed for Unix systems
}
