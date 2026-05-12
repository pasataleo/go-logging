package logging

import (
	"io"
	"os"
	"time"
)

// Opt is a functional option for configuring a logger.
type Opt func(*Opts)

// Opts holds resolved configuration for a logger.
type Opts struct {
	Writer     io.Writer
	Location   *time.Location
	TimeFormat string
}

// DefaultOpts returns the default options with any overrides applied.
func DefaultOpts(options []Opt) Opts {
	o := Opts{
		Writer:     os.Stderr,
		Location:   time.UTC,
		TimeFormat: time.RFC3339,
	}
	for _, opt := range options {
		opt(&o)
	}
	return o
}

// WithWriter sets the output destination for log entries.
func WithWriter(w io.Writer) Opt {
	return func(o *Opts) {
		o.Writer = w
	}
}

// WithTimezone sets the timezone used when formatting log timestamps.
func WithTimezone(loc *time.Location) Opt {
	return func(o *Opts) {
		o.Location = loc
	}
}

// WithTimeFormat sets the format string used for log timestamps. The format
// follows the conventions of [time.Time.Format].
func WithTimeFormat(format string) Opt {
	return func(o *Opts) {
		o.TimeFormat = format
	}
}
