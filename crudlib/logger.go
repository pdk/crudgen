package crudlib

// globalLogger is a private, global where we keep a logger. We can
// conditionally log what's happening with the database by setting/unsetting
// this var.
var globalLogger Logger

// Logger is a thing that can write to a log.
type Logger interface {
	Printf(format string, v ...interface{})
}

// SetLogger sets a logger so we'll start logging database actions.
func SetLogger(logger Logger) {
	globalLogger = logger
}

// UnsetLogger disables logging of SQL actions.
func UnsetLogger() {
	globalLogger = nil
}

// Log will write a line to the logger, IFF it has been set to a non-nil value.
func Log(format string, v ...interface{}) {
	if globalLogger != nil {
		globalLogger.Printf(format, v...)
	}
}
