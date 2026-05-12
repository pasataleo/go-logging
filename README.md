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
