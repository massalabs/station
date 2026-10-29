package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	maxSizePerFile = 10 // megabytes
	maxBackups     = 5
	maxAge         = 60 // days
)

// unsyncedWriter hides the Sync method of the wrapped writer from zapcore.AddSync.
// Syncing stdout is useless (it is not buffered) and fails when stdout is a pipe or a terminal
// (e.g. "sync /dev/stdout: bad file descriptor" for plugins whose output is piped to MassaStation).
type unsyncedWriter struct {
	io.Writer
}

// Logger holds the logger instance.
type Logger struct {
	*zap.SugaredLogger
	path      string
	startTime time.Time
	file      io.Closer
}

// Path returns the path of the file the logger writes to.
func (l *Logger) Path() string {
	return l.path
}

// StartTime returns when the logger was created, i.e. the start of the current logging session.
func (l *Logger) StartTime() time.Time {
	return l.startTime
}

// Close closes the logger.
// It must be called before the program exits.
func (l *Logger) Close() error {
	if l.SugaredLogger == nil {
		return nil
	}

	err := l.Sync()
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

	rotator := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizePerFile,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   false,
	}

	// Create the log core.
	fileWriteSyncer := zapcore.AddSync(rotator)
	consoleWriteSyncer := zapcore.AddSync(unsyncedWriter{os.Stdout})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer),
		logLevel,
	)

	// Create the logger.
	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{SugaredLogger: logger.Sugar(), path: path, startTime: time.Now(), file: rotator}, nil
}
