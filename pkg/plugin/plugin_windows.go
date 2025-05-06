//go:build windows
// +build windows

package plugin

import (
	"os/exec"
	"syscall"

	"github.com/massalabs/station/pkg/logger"
	"golang.org/x/sys/windows"
)

// sendStopSignal sends an appropriate signal to gracefully stop a process on Windows
func (p *Plugin) sendStopSignal() error {
	err := windows.GenerateConsoleCtrlEvent(windows.CTRL_C_EVENT, uint32(p.command.Process.Pid))
	if err != nil {
		logger.Warnf("Failed to send Ctrl+C event to plugin %s: %v\n", p.ID, err)
		return err
	}

	logger.Infof("Sent Ctrl+C event to plugin %s.\n", p.ID)
	return nil
}

// setupProcess configures Windows-specific process attributes
func setupProcess(cmd *exec.Cmd) {
	// Set CREATE_NEW_PROCESS_GROUP flag to make it respond to Ctrl+C
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
