package config

import (
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//nolint:gochecknoglobals
var Logger *zap.Logger

func NewLogger() *zap.Logger {
	logDir, err := GetConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	logFilePath := filepath.Join(logDir, "massastation.log")
	// Create the log directory if it doesn't exist
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
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
	return zap.New(core, zap.AddCaller())
}
