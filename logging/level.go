package logging

import (
	"flag"
	"fmt"

	"github.com/pasataleo/go-errors/errors"
)

type Level int

var (
	// Let's make sure we can parse Level as a flag.
	_ flag.Value = (*Level)(nil)
)

const (
	None Level = iota
	Error
	Warn
	Info
	Trace
)

func (level Level) String() string {
	switch level {
	case None:
		return "none"
	case Error:
		return "error"
	case Warn:
		return "warn"
	case Info:
		return "info"
	case Trace:
		return "trace"
	default:
		panic(fmt.Sprintf("unrecognized log level: %d", level))
	}
}

func (level *Level) Set(value string) error {
	switch value {
	case "none":
		*level = None
	case "error":
		*level = Error
	case "warn":
		*level = Warn
	case "info":
		*level = Info
	case "trace":
		*level = Trace
	default:
		return errors.Newf(nil, errors.ErrorCodeUnknown, "invalid log level %s, must be one of [\"none\", \"error\", \"warn\", \"info\", \"trace\"]", value)
	}

	return nil
}

func (level Level) includes(other Level) bool {
	return other <= level
}
