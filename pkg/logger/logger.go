package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/massalabs/station/pkg/dirs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//nolint:gochecknoglobals
var Logger *zap.SugaredLogger

func NewLogger() *zap.SugaredLogger {
	logDir, err := dirs.GetConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	logDirPath := filepath.Join(logDir, "logs")
	logFilePath := filepath.Join(logDirPath, "massastation.log")

	// Create the log directory if it doesn't exist
	if err := os.MkdirAll(logDirPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Configure log encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//nolint:gomnd
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	// Create the log core
	fileWriteSyncer := zapcore.AddSync(logFile)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer),
		zapcore.InfoLevel)

	// Create the logger
	logger := zap.New(core, zap.AddCaller())

	return logger.Sugar()
}
