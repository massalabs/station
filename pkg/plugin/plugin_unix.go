//go:build !windows
// +build !windows

package plugin

import (
	"os/exec"
	"syscall"

	"github.com/massalabs/station/pkg/logger"
)

// sendStopSignal sends an appropriate signal to gracefully stop a process on Unix systems
func (p *Plugin) sendStopSignal() error {
	err := p.command.Process.Signal(syscall.SIGTERM)
	if err != nil {
		logger.Warnf("Failed to send SIGTERM to plugin %s: %s\n", p.ID, err)
		return err
	}

	logger.Infof("Sent SIGTERM to plugin %s.\n", p.ID)
	return nil
}

// setupProcess configures platform-specific process attributes (no-op on Unix)
func setupProcess(cmd *exec.Cmd) {
	// No special setup needed for Unix systems
}
