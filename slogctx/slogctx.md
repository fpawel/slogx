Here is an English version of the `README.md` for your package:

```markdown
# slogctx

`slogctx` is a Go package that extends the standard [`log/slog`](https://pkg.go.dev/log/slog) logger with the ability to add attributes (key-value pairs) from `context.Context` to log records. This makes it easy to pass additional data through context and automatically include it in your logs.

## Features

- **Add attributes to context:**  
  Use `WithValues` and `WithUniqueValues` to add key-value pairs to a context.
    - `WithValues` appends new attributes, allowing duplicate keys.
    - `WithUniqueValues` replaces values for existing keys.

- **Remove attributes from context:**  
    - `WithoutKeys` removes specified keys from the context attributes.
    - `WithoutAllKeys` clears all attributes.

- **Extract attributes:**  
    - `GetFirstValue` returns the first value for a key from the context.
    - `HasKey` checks if a key exists in the context attributes.

- **Integration with slog.Handler:**  
    - `NewHandler` wraps an existing `slog.Handler` and adds support for context attributes.
    - All context attributes are automatically added to each log record.

## Usage Example

```go
import (
    "context"
    "log/slog"
    "github.com/your/module/slogctx"
)

func main() {
    baseHandler := slog.NewJSONHandler(os.Stdout, nil)
    logger := slog.New(slogctx.NewHandler(baseHandler))

    ctx := slogctx.WithValues(context.Background(), "userID", 42, "requestID", "abc123")
    logger.InfoContext(ctx, "User request received")
}
```

## API Reference

- **NewHandler(handler slog.Handler) slog.Handler**  
  Wraps the given handler and adds support for context attributes.

- **WithValues(ctx context.Context, args ...any) context.Context**  
  Adds key-value pairs to the context. Keys must be strings.

- **WithUniqueValues(ctx context.Context, args ...any) context.Context**  
  Adds key-value pairs, replacing values for duplicate keys.

- **WithoutKeys(ctx context.Context, keys ...string) context.Context**  
  Removes the specified keys from the context attributes.

- **WithoutAllKeys(ctx context.Context) context.Context**  
  Removes all attributes from the context.

- **GetFirstValue(ctx context.Context, key string) (any, bool)**  
  Returns the first value for the key and a boolean indicating if it was found.

- **HasKey(ctx context.Context, key string) bool**  
  Checks if the key exists in the context attributes.

## Notes and Limitations

- Context is not intended for storing large amounts of data.
- Each attribute modification creates a new slice, which is safe for concurrent use but may be inefficient with many attributes.
- If an odd number of arguments is passed to `WithValues` or `WithUniqueValues`, the last argument is ignored.
- Duplicate keys are allowed in `WithValues` and are not filtered.

## License

[MIT](../LICENSE)
```