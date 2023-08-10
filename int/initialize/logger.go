package initialize

import (
	"fmt"
	"path/filepath"

	"github.com/massalabs/station/pkg/logger"
)

const (
	LogFileName       = "massastation.log"
	RepairLogFileName = "massastation-repair.log"
	LogSubDirName     = "logs"
)

// Logger sets up the global logger.
func Logger(repairMode bool, configDir string) error {
	logDirPath := filepath.Join(configDir, LogSubDirName)

	var logFilePath string

	if repairMode {
		logFilePath = filepath.Join(logDirPath, RepairLogFileName)
	} else {
		logFilePath = filepath.Join(logDirPath, LogFileName)
	}

	if err := logger.InitializeGlobal(logFilePath); err != nil {
		return fmt.Errorf("initialize global logger: %w", err)
	}

	return nil
}
