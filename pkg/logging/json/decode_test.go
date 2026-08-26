package json

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pasataleo/go-logging/pkg/logging"
	"github.com/pasataleo/go-testingx/pkg/diff"
	"github.com/pasataleo/go-testingx/pkg/testingx"
)

// timeDiffer compares time.Time by its formatted value. The default walker
// cannot descend into a time.Time because its fields are unexported, and the
// formatted timestamp is both sufficient here and far more readable than a
// struct dump when a comparison fails.
type timeDiffer struct{}

func (timeDiffer) Diff(got, want time.Time, _ *diff.Opts) diff.Result {
	return diff.ValueResult(
		reflect.ValueOf(got.Format(time.RFC3339Nano)),
		reflect.ValueOf(want.Format(time.RFC3339Nano)),
	)
}

// withTime lets Equal compare an Entry, which carries a time.Time.
var withTime = diff.WithDiffer[time.Time](timeDiffer{})

// capture returns a logger writing into a buffer, plus the buffer, so a test
// can assert on exactly what the logger produced.
func capture(level logging.Level, options ...logging.Opt) (logging.Logger, *strings.Builder) {
	var buf strings.Builder
	return New(level, append([]logging.Opt{logging.WithWriter(&buf)}, options...)...), &buf
}

// timestamp pulls the raw time field out of a written line and parses it the
// way Decode would, so a test can predict the Time an entry will carry.
func timestamp(t *testing.T, line string, format string) time.Time {
	t.Helper()

	var raw struct {
		Time string `json:"time"`
	}
	testingx.Call(t, json.Unmarshal, []byte(line), &raw).Fatal().NoError()
	return testingx.CallAs[time.Time](t, time.ParseInLocation, format, raw.Time, time.UTC).
		Fatal().NoError().Capture()
}

// The guarantee this package makes: whatever Logger writes, Decode reads back.
func TestRoundTrip(t *testing.T) {
	logger, buf := capture(logging.Trace)
	logger.Named("db").WithFields(map[string]any{"rows": 3, "table": "users"}).Info("handled request")

	line := buf.String()

	// Field values arrive as their JSON types, so 3 comes back as a float64.
	testingx.Call(t, Decode, []byte(line)).NoError().Equal(&Entry{
		Time:    timestamp(t, line, time.RFC3339),
		Level:   logging.Info,
		Message: "handled request",
		Logger:  "db",
		Fields:  map[string]any{"rows": float64(3), "table": "users"},
	}, withTime)
}

func TestRoundTrip_EveryLevel(t *testing.T) {
	// Fatal exits and Panic panics on the producing side, so they cannot be
	// written here; they are covered by TestEmit_FatalIsDowngraded.
	for _, tc := range []struct {
		level logging.Level
		log   func(logging.Logger, string)
	}{
		{logging.Trace, func(l logging.Logger, m string) { l.Trace(m) }},
		{logging.Debug, func(l logging.Logger, m string) { l.Debug(m) }},
		{logging.Info, func(l logging.Logger, m string) { l.Info(m) }},
		{logging.Warn, func(l logging.Logger, m string) { l.Warn(m) }},
		{logging.Error, func(l logging.Logger, m string) { l.Error(m) }},
	} {
		logger, buf := capture(logging.Trace)
		tc.log(logger, "message")

		line := buf.String()
		testingx.Call(t, Decode, []byte(line)).NoError().Equal(&Entry{
			Time:    timestamp(t, line, time.RFC3339),
			Level:   tc.level,
			Message: "message",
		}, withTime)
	}
}

func TestRoundTrip_CustomTimeFormat(t *testing.T) {
	const format = "2006-01-02 15:04:05"

	logger, buf := capture(logging.Trace, logging.WithTimeFormat(format))
	logger.Info("hello")
	line := buf.String()

	// Reading with the producer's options parses the timestamp.
	testingx.Call(t, Decode, []byte(line), logging.WithTimeFormat(format)).NoError().Equal(&Entry{
		Time:    timestamp(t, line, format),
		Level:   logging.Info,
		Message: "hello",
	}, withTime)

	// Reading without them still yields the entry, just without a timestamp:
	// the level and message are worth more than the time.
	testingx.Call(t, Decode, []byte(line)).NoError().Equal(&Entry{
		Level:   logging.Info,
		Message: "hello",
	}, withTime)
}

func TestDecode_Rejects(t *testing.T) {
	for _, tc := range []struct{ name, line string }{
		{"not json", "panic: runtime error: index out of range"},
		{"no level", `{"msg":"hello"}`},
		{"unknown level", `{"level":"nonsense","msg":"hello"}`},
		{"json but not an object", `["hello"]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testingx.Call(t, Decode, []byte(tc.line)).Error()
		})
	}
}

func TestEmit_PreservesNameLevelAndFields(t *testing.T) {
	logger, buf := capture(logging.Trace)

	(&Entry{
		Level:   logging.Warn,
		Message: "slow response",
		Logger:  "http",
		Fields:  map[string]any{"latency": "2.3s"},
	}).Emit(logger.Named("plugin"))

	line := buf.String()

	// Names nest, so a relayed entry sits under the host's logger name.
	testingx.Call(t, Decode, []byte(line)).NoError().Equal(&Entry{
		Time:    timestamp(t, line, time.RFC3339),
		Level:   logging.Warn,
		Message: "slow response",
		Logger:  "plugin.http",
		Fields:  map[string]any{"latency": "2.3s"},
	}, withTime)
}

// Emit must never terminate the relaying process. If this regresses, the test
// binary exits and the whole package fails rather than just this test.
func TestEmit_FatalIsDowngraded(t *testing.T) {
	for _, level := range []logging.Level{logging.Fatal, logging.Panic} {
		logger, buf := capture(logging.Trace)

		entry := &Entry{Level: level, Message: "unrecoverable"}
		entry.Emit(logger)

		line := buf.String()
		testingx.Call(t, Decode, []byte(line)).NoError().Equal(&Entry{
			Time:    timestamp(t, line, time.RFC3339),
			Level:   logging.Error,
			Message: "unrecoverable",
			Fields:  map[string]any{"level": level.String()},
		}, withTime)

		// The entry itself still reports what the producer wrote.
		testingx.Capture(t, entry.Level).Equal(level)
	}
}

func TestEmit_NoneIsDropped(t *testing.T) {
	logger, buf := capture(logging.Trace)
	(&Entry{Level: logging.None, Message: "ignored"}).Emit(logger)

	testingx.Capture(t, buf.String()).Empty()
}

func TestEmit_NilEntry(t *testing.T) {
	logger, buf := capture(logging.Trace)

	var entry *Entry
	entry.Emit(logger)

	testingx.Capture(t, buf.String()).Empty()
}

func TestPipe(t *testing.T) {
	// Two real entries with a panic trace and a blank line between them.
	source, sourceBuf := capture(logging.Trace)
	source.Info("first")
	source.Named("db").Error("second")

	lines := strings.Split(strings.TrimRight(sourceBuf.String(), "\n"), "\n")
	testingx.Capture(t, lines).Fatal().Len(2)

	stream := lines[0] + "\n" +
		"panic: runtime error: index out of range\n" +
		"\n" +
		lines[1] + "\n"

	dest, destBuf := capture(logging.Trace)
	testingx.Call(t, Pipe, dest.Named("plugin"), strings.NewReader(stream)).Fatal().NoError()

	out := strings.Split(strings.TrimRight(destBuf.String(), "\n"), "\n")
	testingx.Capture(t, out).Fatal().Len(3)

	testingx.Call(t, Decode, []byte(out[0])).NoError().Equal(&Entry{
		Time:    timestamp(t, out[0], time.RFC3339),
		Level:   logging.Info,
		Message: "first",
		Logger:  "plugin",
	}, withTime)

	// The undecodable line is relayed verbatim at Warn.
	testingx.Call(t, Decode, []byte(out[1])).NoError().Equal(&Entry{
		Time:    timestamp(t, out[1], time.RFC3339),
		Level:   logging.Warn,
		Message: "panic: runtime error: index out of range",
		Logger:  "plugin",
	}, withTime)

	testingx.Call(t, Decode, []byte(out[2])).NoError().Equal(&Entry{
		Time:    timestamp(t, out[2], time.RFC3339),
		Level:   logging.Error,
		Message: "second",
		Logger:  "plugin.db",
	}, withTime)
}

// A final line with no trailing newline must still be replayed.
func TestPipe_UnterminatedFinalLine(t *testing.T) {
	source, sourceBuf := capture(logging.Trace)
	source.Info("last")

	dest, destBuf := capture(logging.Trace)
	testingx.Call(t, Pipe, dest, strings.NewReader(strings.TrimRight(sourceBuf.String(), "\n"))).
		Fatal().NoError()

	line := destBuf.String()
	testingx.Call(t, Decode, []byte(line)).NoError().Equal(&Entry{
		Time:    timestamp(t, line, time.RFC3339),
		Level:   logging.Info,
		Message: "last",
	}, withTime)
}

func TestPipe_Empty(t *testing.T) {
	dest, destBuf := capture(logging.Trace)
	testingx.Call(t, Pipe, dest, strings.NewReader("")).NoError()
	testingx.Capture(t, destBuf.String()).Empty()
}

// The destination logger's level still applies to replayed entries.
func TestPipe_RespectsDestinationLevel(t *testing.T) {
	source, sourceBuf := capture(logging.Trace)
	source.Debug("chatty")

	dest, destBuf := capture(logging.Warn)
	testingx.Call(t, Pipe, dest, strings.NewReader(sourceBuf.String())).NoError()
	testingx.Capture(t, destBuf.String()).Empty()
}

// The wire format itself, so that changing it has to be deliberate.
func TestLogger_WireFormat(t *testing.T) {
	logger, buf := capture(logging.Trace)
	logger.Named("db").WithField("rows", 3).Info("handled request")

	var entry map[string]any
	testingx.Call(t, json.Unmarshal, []byte(buf.String()), &entry).Fatal().NoError()

	testingx.Capture(t, entry).Equal(map[string]any{
		"time":   entry["time"],
		"level":  "info",
		"msg":    "handled request",
		"logger": "db",
		"fields": map[string]any{"rows": float64(3)},
	})
	testingx.Call(t, time.Parse, time.RFC3339, entry["time"]).NoError()
}
