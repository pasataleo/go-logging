package test

import (
	"fmt"
	"maps"
	"strings"
	"testing"

	"github.com/pasataleo/go-logging/pkg/logging"
)

// Logger is a Logger implementation that forwards log messages to a
// [testing.TB], so output appears in test logs and respects -v filtering.
type Logger struct {
	tb     testing.TB
	level  logging.Level
	name   string
	fields map[string]any
}

// New creates a new test Logger that forwards messages at or above the given
// level to tb.
func New(tb testing.TB, level logging.Level) logging.Logger {
	return &Logger{
		tb:     tb,
		level:  level,
		fields: make(map[string]any),
	}
}

func (l *Logger) log(level logging.Level, msg string) {
	l.tb.Helper()
	if !l.level.ShouldLog(level) {
		return
	}

	var b strings.Builder
	b.WriteString("[")
	b.WriteString(level.String())
	b.WriteString("]")
	if l.name != "" {
		b.WriteString(" [")
		b.WriteString(l.name)
		b.WriteString("]")
	}
	b.WriteString(" ")
	b.WriteString(msg)
	for k, v := range l.fields {
		b.WriteString(" ")
		b.WriteString(k)
		b.WriteString("=")
		fmt.Fprintf(&b, "%v", v)
	}

	l.tb.Log(b.String())
}

func (l *Logger) Trace(msg string) { l.tb.Helper(); l.log(logging.Trace, msg) }
func (l *Logger) Tracef(msg string, args ...any) {
	l.tb.Helper()
	l.log(logging.Trace, fmt.Sprintf(msg, args...))
}
func (l *Logger) Debug(msg string) { l.tb.Helper(); l.log(logging.Debug, msg) }
func (l *Logger) Debugf(msg string, args ...any) {
	l.tb.Helper()
	l.log(logging.Debug, fmt.Sprintf(msg, args...))
}
func (l *Logger) Info(msg string) { l.tb.Helper(); l.log(logging.Info, msg) }
func (l *Logger) Infof(msg string, args ...any) {
	l.tb.Helper()
	l.log(logging.Info, fmt.Sprintf(msg, args...))
}
func (l *Logger) Warn(msg string) { l.tb.Helper(); l.log(logging.Warn, msg) }
func (l *Logger) Warnf(msg string, args ...any) {
	l.tb.Helper()
	l.log(logging.Warn, fmt.Sprintf(msg, args...))
}
func (l *Logger) Error(msg string) { l.tb.Helper(); l.log(logging.Error, msg) }
func (l *Logger) Errorf(msg string, args ...any) {
	l.tb.Helper()
	l.log(logging.Error, fmt.Sprintf(msg, args...))
}

func (l *Logger) Fatal(msg string) {
	l.tb.Helper()
	l.log(logging.Fatal, msg)
	l.tb.FailNow()
}

func (l *Logger) Fatalf(msg string, args ...any) {
	l.tb.Helper()
	l.log(logging.Fatal, fmt.Sprintf(msg, args...))
	l.tb.FailNow()
}

func (l *Logger) Panic(msg string) {
	l.tb.Helper()
	l.log(logging.Panic, msg)
	panic(msg)
}

func (l *Logger) Panicf(msg string, args ...any) {
	l.tb.Helper()
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
		tb:     l.tb,
		level:  l.level,
		name:   l.name,
		fields: fields,
	}
}

func (l *Logger) WithFields(fields map[string]any) logging.Logger {
	merged := make(map[string]any, len(l.fields)+len(fields))
	maps.Copy(merged, l.fields)
	maps.Copy(merged, fields)
	return &Logger{
		tb:     l.tb,
		level:  l.level,
		name:   l.name,
		fields: merged,
	}
}

func (l *Logger) Named(name string) logging.Logger {
	n := name
	if l.name != "" {
		n = l.name + "." + name
	}
	return &Logger{
		tb:     l.tb,
		level:  l.level,
		name:   n,
		fields: l.fields,
	}
}

func (l *Logger) ShouldLog(level logging.Level) bool {
	return l.level.ShouldLog(level)
}
