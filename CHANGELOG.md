# v0.2.0

## FEATURES

- `logging/json` now reads its own output as well as writing it. `Decode` parses a written line into an `Entry`, `Entry.Emit` replays one onto a logger at the level the producer recorded, and `Pipe` composes the two over an `io.Reader` — so a process capturing another process's stderr can put those lines into its own log stream, at the original levels, with the original names and fields
- `Decode` takes the same `logging.Opt` values as `New`, because `WithTimeFormat` and `WithTimezone` change the bytes on the wire; a timestamp that will not parse leaves `Entry.Time` zero rather than failing the whole line
- Lines that are not entries — malformed JSON, or JSON without a `level` field — are relayed verbatim at `Warn` by `Pipe`, which covers panics, runtime traces, and producers that do not use this logger
- `Level.MarshalText` and `Level.UnmarshalText`, so a `Level` round-trips through `encoding/json` and marshals as `"info"` rather than `3` in any struct, not only in a log entry

## IMPROVEMENTS

- `pkg/logging` (the interface package) no longer imports `go-colour`. ANSI colour handling for log level names is now owned by `pkg/logging/text`, the package that actually produces coloured output. Concrete impact: consumers that only need the `Logger` interface or `Level` type — for example WASM/embedded applications — no longer transitively pull in `go-colour` (and through it, `regexp` and `regexp/syntax`). For a representative WASM build this drops ~75 KB raw / ~30 KB gzipped from the final binary.
- `Entry.Emit` replays entries written at `Fatal` or `Panic` at `Error` instead, preserving the original level as a `level` field. `Logger.Fatal` calls `os.Exit(1)` and `Logger.Panic` panics, which is right for the process that hit the error and wrong for one merely relaying someone else's line — without this, a single plugin logging at `Fatal` would terminate its host.
- `pkg/logging` gained its first tests, along with `pkg/logging/json`. The round-trip tests are the point: they pin the writer and the reader together, so a change to the entry format cannot silently break a consumer.
- `pkg/logging` still imports nothing outside `fmt`, `io`, `os`, and `time`. The new case-insensitive comparison in `SetString` is deliberately ASCII-only rather than reaching for `strings.ToLower` and the unicode case tables, for the same binary-size reason that moved colour handling into `logging/text`.

## BUG FIXES

- `Level.SetString` documented itself as case-insensitive but compared against lowercase literals, so `"INFO"` and `"Info"` were rejected as unknown levels. It now matches its documentation, which also fixes case-insensitivity for `LOG_LEVEL` environment variables and `--log-level` flags read through go-config.

## BREAKING CHANGES

- `Level.ColouredString()` has been removed from `pkg/logging`. The implementation has moved into `pkg/logging/text` as a private helper (`colouredLevel`), since the text logger was the only caller. The level name → ANSI colour mapping is unchanged.

## UPGRADE NOTES

- If you were calling `Level.ColouredString()` directly: either switch to using `logging/text.New(...)` (which still produces coloured output), or copy the level-to-colour switch into your own code. The mapping is the seven levels you'd expect: Trace=grey, Debug=cyan, Info=green, Warn=yellow, Error=red, Fatal=bright-red, Panic=bright-red.

<!--
## SECURITY
Vulnerabilities or security-related changes addressed in this release.
-->

<!--
## DEPRECATIONS
Functionality that will be removed in a future release.
-->

# v0.1.0

## FEATURES

- `logging.Logger` interface with seven severity levels: Trace, Debug, Info, Warn, Error, Fatal, Panic
- `logging.Level` type with `SetString` support for use as a flag or environment variable field with go-config
- `logging/json` backend — structured JSON lines with caller fields nested under `"fields"`
- `logging/text` backend — human-readable lines with ANSI-coloured level names
- `logging/test` backend — routes output through `testing.TB`; `Fatal` calls `tb.FailNow()` instead of `os.Exit`
- `logging/slog` backend — `slog.Handler` bridge that forwards to a `logging.Logger` with level mapping
- Functional options: `WithWriter`, `WithTimezone`, `WithTimeFormat`
- go-config inject modules for each backend; level configurable via `--log-level` flag or `LOG_LEVEL` env var
