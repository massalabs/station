package initialize

import (
	"fmt"
	"path/filepath"

	"github.com/massalabs/station/pkg/dirs"
	"github.com/massalabs/station/pkg/logger"
)

const (
	LogFileName   = "massastation.log"
	LogSubDirName = "logs"
)

// Logger sets up the global logger.
func Logger() error {
	logDir, err := dirs.GetConfigDir()
	if err != nil {
		return fmt.Errorf("get config dir: %w", err)
	}

	logDirPath := filepath.Join(logDir, LogSubDirName)
	logFilePath := filepath.Join(logDirPath, LogFileName)

	if err := logger.InitializeGlobal(logFilePath); err != nil {
		return fmt.Errorf("initialize global logger: %w", err)
	}

	return nil
}
