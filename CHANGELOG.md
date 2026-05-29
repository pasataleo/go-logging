# v0.2.0

## IMPROVEMENTS

- `pkg/logging` (the interface package) no longer imports `go-colour`. ANSI colour handling for log level names is now owned by `pkg/logging/text`, the package that actually produces coloured output. Concrete impact: consumers that only need the `Logger` interface or `Level` type — for example WASM/embedded applications — no longer transitively pull in `go-colour` (and through it, `regexp` and `regexp/syntax`). For a representative WASM build this drops ~75 KB raw / ~30 KB gzipped from the final binary.

## BREAKING CHANGES

- `Level.ColouredString()` has been removed from `pkg/logging`. The implementation has moved into `pkg/logging/text` as a private helper (`colouredLevel`), since the text logger was the only caller. The level name → ANSI colour mapping is unchanged.

## UPGRADE NOTES

- If you were calling `Level.ColouredString()` directly: either switch to using `logging/text.New(...)` (which still produces coloured output), or copy the level-to-colour switch into your own code. The mapping is the seven levels you'd expect: Trace=grey, Debug=cyan, Info=green, Warn=yellow, Error=red, Fatal=bright-red, Panic=bright-red.

<!--
## FEATURES
New functionality added in this release.
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
