package json

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"time"

	"github.com/pasataleo/go-logging/pkg/logging"
)

// Logger is a Logger implementation that writes structured JSON lines to an
// io.Writer.
type Logger struct {
	level      logging.Level
	writer     io.Writer
	location   *time.Location
	timeFormat string
	name       string
	fields     map[string]any
}

// New creates a new JSON Logger at the given level. By default it writes to
// os.Stderr, uses UTC, and formats timestamps as RFC3339. Use functional
// options to override these defaults.
func New(level logging.Level, options ...logging.Opt) logging.Logger {
	o := logging.DefaultOpts(options)
	return &Logger{
		level:      level,
		writer:     o.Writer,
		location:   o.Location,
		timeFormat: o.TimeFormat,
		fields:     make(map[string]any),
	}
}

func (l *Logger) log(level logging.Level, msg string) {
	if !l.level.ShouldLog(level) {
		return
	}

	entry := map[string]any{
		"time":  time.Now().In(l.location).Format(l.timeFormat),
		"level": level.String(),
		"msg":   msg,
	}
	if l.name != "" {
		entry["logger"] = l.name
	}
	if len(l.fields) > 0 {
		entry["fields"] = l.fields
	}

	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(l.writer, `{"level":"error","msg":"failed to marshal log entry: %s"}`+"\n", err)
		return
	}
	_, _ = l.writer.Write(append(data, '\n'))
}

func (l *Logger) Trace(msg string)               { l.log(logging.Trace, msg) }
func (l *Logger) Tracef(msg string, args ...any) { l.log(logging.Trace, fmt.Sprintf(msg, args...)) }
func (l *Logger) Debug(msg string)               { l.log(logging.Debug, msg) }
func (l *Logger) Debugf(msg string, args ...any) { l.log(logging.Debug, fmt.Sprintf(msg, args...)) }
func (l *Logger) Info(msg string)                { l.log(logging.Info, msg) }
func (l *Logger) Infof(msg string, args ...any)  { l.log(logging.Info, fmt.Sprintf(msg, args...)) }
func (l *Logger) Warn(msg string)                { l.log(logging.Warn, msg) }
func (l *Logger) Warnf(msg string, args ...any)  { l.log(logging.Warn, fmt.Sprintf(msg, args...)) }
func (l *Logger) Error(msg string)               { l.log(logging.Error, msg) }
func (l *Logger) Errorf(msg string, args ...any) { l.log(logging.Error, fmt.Sprintf(msg, args...)) }

func (l *Logger) Fatal(msg string) {
	l.log(logging.Fatal, msg)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...any) {
	l.log(logging.Fatal, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (l *Logger) Panic(msg string) {
	l.log(logging.Panic, msg)
	panic(msg)
}

func (l *Logger) Panicf(msg string, args ...any) {
	formatted := fmt.Sprintf(msg, args...)
	l.log(logging.Panic, formatted)
	panic(formatted)
}

func (l *Logger) WithError(err error) logging.Logger {
	return l.WithField("error", err.Error())
}

func (l *Logger) WithField(key string, value any) logging.Logger {
	fields := make(map[string]any, len(l.fields)+1)
	maps.Copy(fields, l.fields)
	fields[key] = value
	return &Logger{
		level:      l.level,
		writer:     l.writer,
		location:   l.location,
		timeFormat: l.timeFormat,
		name:       l.name,
		fields:     fields,
	}
}

func (l *Logger) WithFields(fields map[string]any) logging.Logger {
	merged := make(map[string]any, len(l.fields)+len(fields))
	maps.Copy(merged, l.fields)
	maps.Copy(merged, fields)
	return &Logger{
		level:      l.level,
		writer:     l.writer,
		location:   l.location,
		timeFormat: l.timeFormat,
		name:       l.name,
		fields:     merged,
	}
}

func (l *Logger) Named(name string) logging.Logger {
	n := name
	if l.name != "" {
		n = l.name + "." + name
	}
	return &Logger{
		level:      l.level,
		writer:     l.writer,
		location:   l.location,
		timeFormat: l.timeFormat,
		name:       n,
		fields:     l.fields,
	}
}

func (l *Logger) ShouldLog(level logging.Level) bool {
	return l.level.ShouldLog(level)
}
