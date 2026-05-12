package logging

import (
	"fmt"

	"github.com/pasataleo/go-colour/pkg/colour"
)

// Level represents a log severity level. Levels are ordered from least to most
// severe: Trace, Debug, Info, Warn, Error, Fatal, Panic.
type Level int8

const (
	// None disables all logging.
	None Level = iota
	// Trace is the least severe level, typically used for fine-grained
	// diagnostic output.
	Trace
	// Debug is used for verbose output useful during development.
	Debug
	// Info is used for general operational messages.
	Info
	// Warn indicates a potential issue that does not prevent normal operation.
	Warn
	// Error indicates a failure that prevented an operation from completing.
	Error
	// Fatal indicates an unrecoverable error. Implementations should
	// terminate the process after logging.
	Fatal
	// Panic indicates an unrecoverable error. Implementations should panic
	// after logging.
	Panic
)

// SetString parses a level name (case-insensitive) and sets the receiver.
// Implements the reflectx.StringSettable interface so Level can be used as a
// config/flag field.
func (l *Level) SetString(v string) error {
	switch v {
	case "none":
		*l = None
	case "trace":
		*l = Trace
	case "debug":
		*l = Debug
	case "info":
		*l = Info
	case "warn":
		*l = Warn
	case "error":
		*l = Error
	case "fatal":
		*l = Fatal
	case "panic":
		*l = Panic
	default:
		return fmt.Errorf("unknown log level: %q", v)
	}
	return nil
}

// String returns the lowercase name of the level (e.g. "info", "warn").
func (l Level) String() string {
	switch l {
	case Trace:
		return "trace"
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	case Panic:
		return "panic"
	default:
		return "none"
	}
}

// ColouredString returns the level name wrapped in ANSI colour codes for
// terminal output.
func (l Level) ColouredString() string {
	switch l {
	case Trace:
		return colour.Grey + "trace" + colour.Reset
	case Debug:
		return colour.Cyan + "debug" + colour.Reset
	case Info:
		return colour.Green + "info" + colour.Reset
	case Warn:
		return colour.Yellow + "warn" + colour.Reset
	case Error:
		return colour.Red + "error" + colour.Reset
	case Fatal:
		return colour.BrightRed + "fatal" + colour.Reset
	case Panic:
		return colour.BrightRed + "panic" + colour.Reset
	default:
		return "none"
	}
}

// ShouldLog reports whether a logger at this level should log a message at the
// given level. For example, Info.ShouldLog(Debug) returns false, and
// Info.ShouldLog(Warn) returns true.
func (l Level) ShouldLog(messageLevel Level) bool {
	return messageLevel >= l
}

// IsEnabledAt reports whether a message at this level would be logged by a
// logger at the given level. For example, Debug.IsEnabledAt(Info) returns false,
// and Warn.IsEnabledAt(Info) returns true.
func (l Level) IsEnabledAt(loggerLevel Level) bool {
	return l >= loggerLevel
}
