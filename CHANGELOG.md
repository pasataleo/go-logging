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

<!--
## IMPROVEMENTS
Enhancements to existing functionality.
-->

<!--
## BUG FIXES
Issues that have been resolved.
-->

<!--
## SECURITY
Vulnerabilities or security-related changes addressed in this release.
-->

<!--
## DEPRECATIONS
Functionality that will be removed in a future release.
-->

<!--
## BREAKING CHANGES
Changes that are not backwards compatible and require updates from consumers.
-->

<!--
## UPGRADE NOTES
Steps required when upgrading from a previous version.
-->
