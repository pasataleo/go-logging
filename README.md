# go-logging

A structured, levelled logging interface for Go with JSON, text, slog, and test backends.

## Installation

```sh
go get github.com/pasataleo/go-logging
```

## Log levels

Levels are ordered from least to most severe:

| Level   | Description                                           |
|---------|-------------------------------------------------------|
| `None`  | Disables all logging                                  |
| `Trace` | Fine-grained diagnostic output                        |
| `Debug` | Verbose output useful during development              |
| `Info`  | General operational messages                          |
| `Warn`  | Potential issues that do not prevent normal operation |
| `Error` | Failures that prevented an operation from completing  |
| `Fatal` | Unrecoverable errors — terminates the process         |
| `Panic` | Unrecoverable errors — panics                         |

`Level` implements `SetString(string) error`, so it can be used directly as a flag or environment variable field with [go-config](https://github.com/pasataleo/go-config).

`Level` also exposes two utility methods:

- `ShouldLog(messageLevel Level) bool` — returns true if `messageLevel` is at least as severe as the receiver (i.e. the receiver is the logger's configured level)
- `IsEnabledAt(loggerLevel Level) bool` — inverse of `ShouldLog`; returns true if the receiver (a message level) is at least as severe as `loggerLevel`

## Logger interface

All loggers implement the `Logger` interface:

```go
type Logger interface {
    Trace(msg string)
    Tracef(msg string, args ...any)
    Debug(msg string)
    Debugf(msg string, args ...any)
    Info(msg string)
    Infof(msg string, args ...any)
    Warn(msg string)
    Warnf(msg string, args ...any)
    Error(msg string)
    Errorf(msg string, args ...any)
    Fatal(msg string)
    Fatalf(msg string, args ...any)
    Panic(msg string)
    Panicf(msg string, args ...any)

    WithError(err error) Logger
    WithField(key string, value any) Logger
    WithFields(fields map[string]any) Logger
    Named(name string) Logger
    ShouldLog(level Level) bool
}
```

## Implementations

### `logging/json`

Writes structured JSON lines. Caller-provided fields are nested under a `"fields"` key to avoid clashing with `time`, `level`, `msg`, and `logger`.

```go
logger := json.New(logging.Info)
logger.WithField("request_id", "abc").Info("handled request")
```

Output:

```json
{"time":"2025-01-01T00:00:00Z","level":"info","msg":"handled request","fields":{"request_id":"abc"}}
```

### `logging/text`

Writes human-readable lines with ANSI-coloured level names.

```go
logger := text.New(logging.Debug)
logger.Named("http").WithField("latency", "2.3s").Warn("slow response")
```

Output:

```
2025-01-01T00:00:00Z  WARN [http] slow response latency=2.3s
```

`Named` appends to the logger name with a dot separator, so `logger.Named("http").Named("server")` produces a logger named `http.server`.

### `logging/test`

Wraps a `testing.TB` so log output appears in test logs and respects `go test -v`. `Fatal` calls `tb.FailNow()` instead of `os.Exit`.

```go
func TestSomething(t *testing.T) {
    logger := test.New(t, logging.Trace)
    logger.Info("test started")
}
```

### `logging/slog`

Bridges with the standard library's `log/slog` package. Provides an `slog.Handler` that forwards records to a `logging.Logger`, mapping `slog.Level` values to the corresponding `logging.Level`.

```go
handler := slog.NewHandler(logger)
slogger := slog.New(handler)
slogger.Info("hello", "key", "value")
```

Supports `WithAttrs` and `WithGroup` for structured context.

## Options

`json.New` and `text.New` accept functional options:

```go
json.New(logging.Info,
    logging.WithWriter(os.Stdout),
    logging.WithTimezone(time.Local),
    logging.WithTimeFormat(time.Kitchen),
)
```

| Option                             | Description             | Default        |
|------------------------------------|-------------------------|----------------|
| `WithWriter(w io.Writer)`          | Output destination      | `os.Stderr`    |
| `WithTimezone(loc *time.Location)` | Timezone for timestamps | `time.UTC`     |
| `WithTimeFormat(format string)`    | Timestamp format string | `time.RFC3339` |

## Modules

Each backend provides a module that integrates with [go-config](https://github.com/pasataleo/go-config)'s inject system. Modules are self-configuring — the framework populates the `Level` field from the `--log-level` flag or `LOG_LEVEL` environment variable (defaults to `info`), then binds a `logging.Logger` for injection into other components.

```go
injector := inject.New()

// Pick one:
injector.Bind(json.NewModule(), "logger")          // production JSON output
injector.Bind(text.NewModule(), "logger")          // human-readable output
injector.Bind(test.NewModule(t), "logger")         // test output via testing.TB

// JSON and text modules also accept options:
injector.Bind(json.NewModule(logging.WithWriter(os.Stdout)), "logger")
```

The slog module bridges to the standard library. It receives the `logging.Logger` via injection and binds an `*slog.Logger`:

```go
injector.Bind(json.NewModule(), "logger")
injector.Bind(slog.NewModule())
```

### Using modules in a config struct

```go
type App struct {
    Logging *json.Module `inject:"logger"`
    Slog    *slog.Module
    Handler *MyHandler   // receives logging.Logger and *slog.Logger via injection
}
```

## Decoding

The `logging/json` package reads its own output as well as writing it. A process that captures another process's stderr
— see [go-plugin](https://github.com/pasataleo/go-plugin) — can turn those lines back into entries and replay them
through its own logger, at the original level, with the original name and fields.

The point of putting this next to the writer is not just convenience: it is that a single round-trip test in this
package pins the two halves together, so a change to the entry format cannot silently break a consumer that had
reimplemented the parse by hand.

```go
type Entry struct {
    Time    time.Time
    Level   logging.Level
    Message string
    Logger  string
    Fields  map[string]any
}

func Decode(line []byte, options ...logging.Opt) (*Entry, error)
func (e *Entry) Emit(logger logging.Logger)
func Pipe(logger logging.Logger, r io.Reader, options ...logging.Opt) error
```

`Decode` parses one line. `Emit` replays an entry onto a logger — applying `Named` and `WithFields` first, then calling
the method matching the entry's level. `Pipe` is the two composed over a stream, and is the only call most consumers
need:

```go
// Replay everything the child process writes to stderr through our own logger.
go json.Pipe(logger.Named("plugin"), stderr)
```

Lines that do not decode are emitted verbatim at `Warn`, which covers panics, runtime traces, and any producer that
does not use this logger. Blank lines are skipped, and a final line with no trailing newline is still replayed.

`Decode` treats a line as an entry only if it carries a `level` field, so arbitrary JSON a process happens to print is
relayed verbatim rather than surfacing as an empty entry.

### Why decoding needs the options

`Decode` takes the same `logging.Opt` values as `New`, because two of them change the bytes on the wire:

- **`WithTimeFormat`** — timestamps are written with `Opts.TimeFormat` (default `time.RFC3339`), so a decoder that
  assumes RFC3339 cannot read a logger configured otherwise. Pass the producer's options to read its output. A
  timestamp that will not parse leaves `Entry.Time` zero rather than failing the whole line.
- **`WithTimezone`** — the location the timestamp was rendered in.

`WithWriter` is ignored; `Pipe` reads from the `io.Reader` it is given.

### Level

Decoding requires `Level` to round-trip through JSON, so it gains `MarshalText` and `UnmarshalText`:

```go
func (l Level) MarshalText() ([]byte, error)
func (l *Level) UnmarshalText(text []byte) error
```

These make `Level` work with `encoding/json` generally, not only here — a `Level` field in any struct now marshals as
 `"info"` rather than `3`. `UnmarshalText` defers to `SetString`, so both are case-insensitive and both reject an
unknown name.

`SetString` previously compared against lowercase literals despite documenting itself as case-insensitive, so `"INFO"`
was rejected; it now matches its documentation. The comparison is deliberately ASCII-only rather than using
`strings.ToLower`, because `pkg/logging` is kept free of dependencies that cost binary size in WASM builds — the same
reasoning that moved ANSI colour handling out into `logging/text`.

### Replayed Fatal and Panic

`Emit` never terminates the calling process. `Logger.Fatal` calls `os.Exit(1)` and `Logger.Panic` panics, which is
correct for the process that hit the error and badly wrong for a process merely relaying someone else's log line — one
plugin logging at `Fatal` would take the host down with it.

Entries at `Fatal` or `Panic` are therefore replayed at `Error`, with the original level preserved as a `level` field
so nothing is lost. `Entry.Level` still reports what the producer wrote.
