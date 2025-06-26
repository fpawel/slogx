# slogtest

`slogtest` is a Go package providing utilities for testing code that uses the standard `log/slog` package. It is designed to help you capture, inspect, and assert log output in your unit tests.

## Features

- **In-memory log handler**: Capture all log records for inspection.
- **Convenient test logger creation**: Easily create loggers for use in tests.
- **Helpers for deterministic log output**: Print logs to stdout without timestamps for stable test output.

## Usage

### Basic Example

```go
import (
    "testing"
    "log/slog"
    "github.com/fpawel/slogx/slogtest"
)

func TestLogging(t *testing.T) {
    logger, observed := slogtest.NewTestLogger(t)
    logger.Info("hello", slog.String("key", "value"))
    logs := observed.Logs()
    if len(logs) != 1 || logs[0].Message != "hello" {
        t.Errorf("unexpected logs: %+v", logs)
    }
}
```

### Print logs without timestamps

```go
import (
    "log/slog"
    "github.com/fpawel/slogx/slogtest"
)

func main() {
    handler := slogtest.NewStdoutTextHandlerWithoutTimestamp()
    logger := slog.New(handler)
    logger.Info("no timestamp here")
}
```

## API

- `NewTestLogger(t *testing.T) (*slog.Logger, *ObservedHandler)`: Returns a logger and an in-memory handler for assertions.
- `NewStdoutTextHandlerWithoutTimestamp() slog.Handler`: Returns a handler that prints logs to stdout without timestamps.
- `ObservedHandler.Logs() []ObservedLog`: Returns all captured logs for inspection.

## License

[MIT](../LICENSE)