# v0.2.0

## FEATURES

- `logging/json` reads its own output: `Decode` parses a line into an `Entry`, `Entry.Emit` replays it onto a logger, and `Pipe` does both over an `io.Reader`.
- `Pipe` relays lines that are not entries verbatim at `Warn`.
- `Level` implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler`.

## IMPROVEMENTS

- `pkg/logging` no longer imports `go-colour`; colouring level names moved to `pkg/logging/text`, which keeps WASM binaries smaller.
- `Entry.Emit` replays `Fatal` and `Panic` entries at `Error`, so a relayed line cannot exit or panic the process relaying it.

## BUG FIXES

- `Level.SetString` is case-insensitive, as documented.

<!--
## SECURITY
Vulnerabilities or security-related changes addressed in this release.
-->

<!--
## DEPRECATIONS
Functionality that will be removed in a future release.
-->

## BREAKING CHANGES

- `Level.ColouredString()` is removed.

## UPGRADE NOTES

- Use `logging/text.New(...)` for coloured output, or map levels to colours in your own code.

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
