package logger

import (
	"fmt"
	"os"
	"strings"
)

// Log levels constants.
const (
	LogDebug = "DEBUG"
	LogInfo  = "INFO"
	LogWarn  = "WARN"
	LogError = "ERROR"
	LogFatal = "FATAL"
)

// LogLevel retrieves the log level from the LOG_LEVEL environment variable.
// If the variable is not set, LogLevel returns INFO.
// LogLevel will return an error if the LOG_LEVEL string does not match one of the following:
// DEBUG, INFO, WARN, ERROR, FATAL.
func LogLevel() (string, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return LogInfo, nil
	}

	switch strings.ToUpper(logLevel) {
	case LogDebug, LogInfo, LogWarn, LogError, LogFatal:
		return logLevel, nil
	default:
		return "", fmt.Errorf("invalid LOG_LEVEL value %q; expecting one of DEBUG, INFO, WARN, ERROR, FATAL", logLevel)
	}
}
