package logging

// Logger defines a levelled, structured logging interface. Implementations
// must support seven severity levels (Trace through Panic), structured fields,
// sub-loggers, and level checking.
type Logger interface {
	// Trace and Tracef log a message at Trace level.
	Trace(msg string)
	Tracef(msg string, args ...any)

	// Debug and Debugf log a message at Debug level.
	Debug(msg string)
	Debugf(msg string, args ...any)

	// Info and Infof log a message at Info level.
	Info(msg string)
	Infof(msg string, args ...any)

	// Warn and Warnf log a message at Warn level.
	Warn(msg string)
	Warnf(msg string, args ...any)

	// Error and Errorf log a message at Error level.
	Error(msg string)
	Errorf(msg string, args ...any)

	// Fatal and Fatalf log a message at Fatal level. Implementations should
	// terminate the process after logging.
	Fatal(msg string)
	Fatalf(msg string, args ...any)

	// Panic and Panicf log a message at Panic level. Implementations should
	// panic after logging.
	Panic(msg string)
	Panicf(msg string, args ...any)

	// WithError returns a new Logger with the given error attached as a
	// structured field.
	WithError(err error) Logger

	// WithField and WithFields return a new Logger with additional key-value
	// pairs added to the structured context.
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger

	// Named returns a new Logger with the given name appended to identify
	// the sub-component producing the log output.
	Named(name string) Logger

	// ShouldLog reports whether this logger would emit a message at the
	// given level. Callers can use this to skip expensive message
	// construction when the level is disabled.
	ShouldLog(level Level) bool
}
