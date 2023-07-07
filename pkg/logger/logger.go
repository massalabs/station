package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FileModeReadWriteReadRead is the file mode used to create the log file.
// It is -rw-r--r--.
const FileModeReadWriteReadRead = 0o644

// Logger holds the logger instance.
type Logger struct {
	*zap.SugaredLogger
	file *os.File
}

// Close closes the logger.
// It must be called before the program exits.
func (l *Logger) Close() error {
	if l.SugaredLogger == nil {
		return nil
	}

	err := l.SugaredLogger.Sync()
	if err != nil {
		return fmt.Errorf("sync logger: %w", err)
	}

	l.SugaredLogger = nil

	err = l.file.Close()
	if err != nil {
		return fmt.Errorf("close log file: %w", err)
	}

	return nil
}

// New creates and returns a new Logger instance.
func New(path string) (*Logger, error) {
	// Get the log level.
	logLevelStr, err := LogLevel()
	if err != nil {
		return nil, fmt.Errorf("get log level: %w", err)
	}

	var logLevel zapcore.Level

	err = logLevel.Set(logLevelStr)
	if err != nil {
		return nil, fmt.Errorf("set log level: %w", err)
	}

	// Configure log encoder.
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FileModeReadWriteReadRead)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	// Create the log core.
	fileWriteSyncer := zapcore.AddSync(logFile)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer),
		logLevel,
	)

	// Create the logger.
	logger := zap.New(core, zap.AddCaller())

	return &Logger{SugaredLogger: logger.Sugar(), file: logFile}, nil
}
