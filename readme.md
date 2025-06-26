# slogx

`slogx` is a set of extensions and utilities for the standard Go logging package [`log/slog`](https://pkg.go.dev/log/slog). The repository is designed to make logging in Go projects more convenient: it improves log readability, simplifies testing, and adds context support.

### [`slogpretty`](slogpretty/slogpretty.md)

Provides human-readable, colored, and structured log output. Suitable for local development and debugging, supports time format customization, attribute grouping, and formatting customization.

### [`slogtest`](slogtest/slogtest.md)

Makes it easy to test code that uses `log/slog`. Includes a handler that stores logs in memory for later inspection, and a handler for outputting logs without timestamps (for deterministic tests).

### [`slogctx`](slogctx/slogctx.md)

Adds support for passing structured attributes via `context.Context`, allowing you to automatically add information about the user, request, and other context-related parameters to logs.

## Installation

```sh
go get github.com/fpawel/slogx
```

## License

[MIT](LICENSE)