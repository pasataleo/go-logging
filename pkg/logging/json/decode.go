package json

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pasataleo/go-logging/pkg/logging"
)

// Entry is the decoded form of a line written by Logger.
//
// Decoding lives next to the writer on purpose. A single round-trip test in
// this package pins the two halves together, so a change to the entry format
// cannot silently break a consumer that had reimplemented the parse by hand.
type Entry struct {
	Time    time.Time
	Level   logging.Level
	Message string
	Logger  string
	Fields  map[string]any
}

// rawEntry mirrors the shape Logger writes. Time stays a string here because
// the format is configurable, so it can only be interpreted once the options
// are known. Level is a pointer so that a line without one is distinguishable
// from a line reading "none".
type rawEntry struct {
	Time    string         `json:"time"`
	Level   *logging.Level `json:"level"`
	Message string         `json:"msg"`
	Logger  string         `json:"logger"`
	Fields  map[string]any `json:"fields"`
}

// Decode parses a single line written by Logger.
//
// It takes the same options as New because two of them change the bytes on the
// wire: WithTimeFormat and WithTimezone. Pass the producer's options to read
// its output. WithWriter is ignored.
//
// A line that is not an entry — malformed JSON, or JSON without a level —
// returns an error. A timestamp that will not parse leaves Entry.Time zero
// rather than failing the line, since the level and message are still worth
// having.
func Decode(line []byte, options ...logging.Opt) (*Entry, error) {
	opts := logging.DefaultOpts(options)

	var raw rawEntry
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, err
	}
	if raw.Level == nil {
		return nil, fmt.Errorf("not a log entry: no level field")
	}

	entry := &Entry{
		Level:   *raw.Level,
		Message: raw.Message,
		Logger:  raw.Logger,
		Fields:  raw.Fields,
	}
	if raw.Time != "" {
		if parsed, err := time.ParseInLocation(opts.TimeFormat, raw.Time, opts.Location); err == nil {
			entry.Time = parsed
		}
	}
	return entry, nil
}

// Emit replays an entry onto logger at the level the producer recorded,
// applying the entry's name and fields first. Names nest, so replaying an
// entry named "db" onto a logger named "plugin" logs under "plugin.db".
//
// Entries at Fatal or Panic are replayed at Error. Logger.Fatal exits the
// process and Logger.Panic panics, which is right for the process that hit the
// error and wrong for one merely relaying someone else's line — otherwise a
// single plugin logging at Fatal would take its host down. The level the
// producer wrote is preserved as a "level" field, and Entry.Level still
// reports it.
func (e *Entry) Emit(logger logging.Logger) {
	if e == nil {
		return
	}

	if e.Logger != "" {
		logger = logger.Named(e.Logger)
	}
	if len(e.Fields) > 0 {
		logger = logger.WithFields(e.Fields)
	}

	switch e.Level {
	case logging.Trace:
		logger.Trace(e.Message)
	case logging.Debug:
		logger.Debug(e.Message)
	case logging.Info:
		logger.Info(e.Message)
	case logging.Warn:
		logger.Warn(e.Message)
	case logging.Error:
		logger.Error(e.Message)
	case logging.Fatal, logging.Panic:
		logger.WithField("level", e.Level.String()).Error(e.Message)
	case logging.None:
		// A producer that wrote "none" meant do not log; honour that.
	default:
		logger.Info(e.Message)
	}
}

// Pipe reads log lines from r and replays each one through logger until r is
// exhausted. It is Decode and Emit composed over a stream, and is the only
// call most consumers need — a process capturing another process's stderr can
// put that output into its own log stream, at the original levels, with the
// original names and fields:
//
//	go json.Pipe(logger.Named("plugin"), stderr)
//
// Lines that do not decode are emitted verbatim at Warn, which covers panics,
// runtime traces, and any producer that does not use this logger. Blank lines
// are skipped.
//
// Pipe blocks until r reports EOF, for which it returns nil, or fails, for
// which it returns that error.
func Pipe(logger logging.Logger, r io.Reader, options ...logging.Opt) error {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')

		// ReadString returns what it read before an error, so a final line
		// with no trailing newline still gets replayed.
		if trimmed := strings.TrimRight(line, "\r\n"); trimmed != "" {
			replay(logger, trimmed, options)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

func replay(logger logging.Logger, line string, options []logging.Opt) {
	entry, err := Decode([]byte(line), options...)
	if err != nil {
		logger.Warn(line)
		return
	}
	entry.Emit(logger)
}
