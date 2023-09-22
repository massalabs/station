package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	maxSizePerFile = 10 // megabytes
	maxBackups     = 5
	maxAge         = 60 // days
)

// Logger holds the logger instance.
type Logger struct {
	*zap.SugaredLogger
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

	rotator := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizePerFile,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   false,
	}

	// Create the log core.
	fileWriteSyncer := zapcore.AddSync(rotator)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer),
		logLevel,
	)

	// Create the logger.
	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{SugaredLogger: logger.Sugar()}, nil
}
