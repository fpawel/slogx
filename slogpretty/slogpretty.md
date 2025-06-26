# slogpretty

`slogpretty` is a custom handler for Go's `log/slog` package that provides human-friendly, colorized, and structured log output. It is designed to improve the developer experience during local development by making logs easier to read and analyze.

## Features

- **Colorized log levels and messages** (if output is a terminal)
- **Customizable timestamp format**
- **Support for attribute groups and attribute rewriting**
- **Fast attribute formatting** (avoids JSON by default, but supports it)
- **Optional source code location in log output**
- **Efficient memory usage** via `sync.Pool` for temporary slices


## Usage

### Basic Example

```go
import (
    "log/slog"
    "os"
    "github.com/fpawel/slogx/slogpretty"
)

func main() {
    handler := slogpretty.NewPrettyHandler().
        WithColorEnabled(true).
        WithTimeLayout("15:04:05").
        WithWriter(os.Stdout)
    logger := slog.New(handler)
    logger.Info("Hello, world!", slog.String("user", "alice"))
}
```

**Output:**
```
15:04:05 INFO  Hello, world! {"user":"alice"}
```

### Attribute Groups

```go
logger.Info("Order created",
    slog.Group("order",
        slog.Int("id", 123),
        slog.String("status", "paid"),
    ),
    slog.String("user", "bob"),
)
```
**Output:**
```
INFO  Order created {"order":{"id":123,"status":"paid"},"user":"bob"}
```

### Custom Attribute Formatter

```go
handler := slogpretty.NewPrettyHandler().
    WithAttrFormatter(func(m map[string]any) string {
        return "ATTRS"
    })
logger := slog.New(handler)
logger.Info("Test", slog.String("foo", "bar"))
```
**Output:**
```
INFO  Test ATTRS
```

## Configuration

- `WithWriter(io.Writer)`: Set output destination.
- `WithTimeLayout(string)`: Set timestamp format.
- `WithLogLevel(slog.Leveler)`: Set minimum log level.
- `WithSourceInfo(bool)`: Enable/disable file and line info.
- `WithAttrRewriter(func)`: Rewrite or filter attributes.
- `WithColorEnabled(bool)`: Enable/disable color output.
- `WithAttrFormatter(func)`: Custom attribute formatting.

## License

[MIT](../LICENSE)
