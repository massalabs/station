package logger

const notInitializedMessage = "global logger not initialized"

// Initialize the logger at package level.
var (
	//nolint:gochecknoglobals
	logger *Logger
)

// InitializeGlobal initializes the global logger at package level.
// Once initialized, the logger can be used by calling directly the classic logging functions (Debug, Info...).
func InitializeGlobal(path string) error {
	var err error

	logger, err = New(path)

	return err
}

// Debug logs a message at the debug level.
// This function takes a variable number of arguments and uses fmt.Sprint to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Debug(args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Debug(args...)
}

// Debugf logs a message at the debug level.
// This function takes a variable number of arguments and uses fmt.Sprintf to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Debugf(template string, args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Debugf(template, args...)
}

// Info logs a message at the info level.
// This function takes a variable number of arguments and uses fmt.Sprint to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Info(args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Info(args...)
}

// Infof logs a message at the info level.
// This function takes a variable number of arguments and uses fmt.Sprintf to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Infof(template string, args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Infof(template, args...)
}

// Warn logs a message at the warn level.
// This function takes a variable number of arguments and uses fmt.Sprint to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Warn(args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Warn(args...)
}

// Warnf logs a message at the warn level.
// This function takes a variable number of arguments and uses fmt.Sprintf to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Warnf(template string, args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Warnf(template, args...)
}

// Error logs a message at the error level.
// This function takes a variable number of arguments and uses fmt.Sprint to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Error(args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Error(args...)
}

// Errorf logs a message at the error level.
// This function takes a variable number of arguments and uses fmt.Sprintf to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Errorf(template string, args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Errorf(template, args...)
}

// Fatal logs a message at the fatal level and stops the program.
// This function takes a variable number of arguments and uses fmt.Sprint to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Fatal(args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Fatal(args...)
}

// Fatalf logs a message at the fatal level and stops the program.
// This function takes a variable number of arguments and uses fmt.Sprintf to construct the log message.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Fatalf(template string, args ...interface{}) {
	if logger == nil {
		panic(notInitializedMessage)
	}

	logger.Fatalf(template, args...)
}

// Close closes the logger.
// This function should be called before the program exits.
// Warning: Calling this function will cause all subsequent calls to the logger to panic.
// Note: The global logger must be initialized before calling this function, otherwise it will panic.
func Close() error {
	if logger == nil {
		panic(notInitializedMessage)
	}

	return logger.Close()
}
