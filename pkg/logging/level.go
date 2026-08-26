package logging

import "fmt"

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
	switch lowerASCII(v) {
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

// MarshalText renders the level as its lowercase name, so a Level field
// marshals as "info" rather than as the integer 3. Implements
// encoding.TextMarshaler.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText parses a level name, case-insensitively. Implements
// encoding.TextUnmarshaler, which is what lets a Level round-trip through
// encoding/json — see logging/json.Decode.
func (l *Level) UnmarshalText(text []byte) error {
	return l.SetString(string(text))
}

// lowerASCII lowercases an ASCII string, returning v unchanged when there is
// nothing to do.
//
// This exists instead of strings.ToLower because this package is deliberately
// kept light: consumers that need only the Logger interface or the Level type
// include it in WASM builds, where every dependency costs binary size, and
// strings.ToLower reaches the unicode case tables.
func lowerASCII(v string) string {
	for i := 0; i < len(v); i++ {
		if c := v[i]; c >= 'A' && c <= 'Z' {
			lowered := []byte(v)
			for j := i; j < len(lowered); j++ {
				if c := lowered[j]; c >= 'A' && c <= 'Z' {
					lowered[j] = c + ('a' - 'A')
				}
			}
			return string(lowered)
		}
	}
	return v
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
