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
	// Windows doesn't have an exact equivalent to SIGTERM
	// For CREATE_NEW_PROCESS_GROUP processes, CTRL_BREAK_EVENT is the closest equivalent
	// The process group ID is the same as the process ID of the group leader
	err := windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(p.command.Process.Pid))
	if err != nil {
		logger.Warnf("Failed to send Ctrl+Break event to plugin %s: %v\n", p.ID, err)
		return err
	}

	logger.Infof("Sent Ctrl+Break event to plugin %s.\n", p.ID)
	return nil
}

// setupProcess configures Windows-specific process attributes
func setupProcess(cmd *exec.Cmd) {
	// CREATE_NEW_PROCESS_GROUP makes the process the leader of its own process group
	// This allows it to receive CTRL_BREAK_EVENT signals
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
