package slogctx

import (
	"context"
	"log/slog"
	"slices"
)

type (
	// Handler implements slog.Handler and supports attributes from context.Context.
	// It allows adding key-value pairs to the context, which are then included in log records.
	// The Handler wraps an existing slog.Handler and adds functionality to handle attributes stored in context.Context.
	//
	// Design notes and caveats:
	//   1. Each call to WithValues, WithUniqueValues or WithoutKeys clones the slice of fields to avoid data races.
	//      This is safe, but may be inefficient if there are many attributes in the context.
	//   2. Context is not intended for storing large amounts of data. Storing many attributes may increase memory usage and reduce performance.
	//   3. If an odd number of arguments is passed to WithValues or WithUniqueValues, the last argument is ignored without warning.
	//      This is for simplicity, but may hide mistakes in argument lists.
	//   4. NewHandler does not check if the handler is nil. Passing nil will cause a panic on first use.
	//      This is for API simplicity, but users must ensure a valid handler is provided.
	//   5. Duplicate keys are allowed and not filtered out in WithValues. Multiple values for the same key may appear in logs.
	//      This is for simplicity, but may complicate log analysis.
	//   6. Removing all keys via WithoutAllKeys sets the context value to nil, which may be non-obvious.
	//
	// Example usage:
	//   logger := slog.New(slogctx.NewHandler(existingHandler))
	//   ctx := slogctx.WithValues(context.Background(), "key1", "value1", "key2", "value2")
	//   logger.InfoContext(ctx, "Log message with context attributes")
	// The attributes will be added to the log record when Handle is called.
	Handler struct {
		slog.Handler
	}

	// fieldsData stores a slice of key-value pairs for logging.
	fieldsData struct {
		slice []field
	}

	// field represents a key-value pair for logging.
	field struct {
		key string
		val any
	}
)

// contextKeyFields is used as a unique key for storing attributes in context.Context.
// This key is unexported, so only this package can access the stored values.
var contextKeyFields = &struct{}{}

// Ensure Handler implements slog.Handler.
var _ slog.Handler = Handler{}

// NewHandler creates a new Handler wrapping the provided slog.Handler.
// Note: For simplicity, handler is not checked for nil. Passing nil will cause a panic on use (see design note 4).
func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{Handler: handler}
}

// Enabled checks if the given log level is enabled for this context.
// It uses the underlying slog.Handler's Enabled method.
func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

// Handle adds attributes from context to the log record.
// It retrieves attributes from context.Context using the contextKeyFields key.
// If multiple values for the same key are present, all are added to the record (see design note 5).
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if p, ok := ctx.Value(contextKeyFields).(*fieldsData); ok {
		for _, f := range p.slice {
			record.AddAttrs(slog.Any(f.key, f.val))
		}
	}
	return h.Handler.Handle(ctx, record)
}

// WithAttrs returns a new Handler with additional attributes.
// The attributes are added to the log record when Handle is called.
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{
		h.Handler.WithAttrs(attrs),
	}
}

// WithGroup returns a new Handler with an attribute group.
// The group name is used to group attributes in the log output.
func (h Handler) WithGroup(name string) slog.Handler {
	return Handler{
		Handler: h.Handler.WithGroup(name),
	}
}

// WithValues returns a new context with the provided key-value pairs added for logging.
// Only string keys are accepted and empty keys are ignored.
// If an odd number of arguments is passed, the last argument is ignored without warning (see design note 3).
// Duplicate keys are allowed and not filtered out for simplicity (see design note 5).
// The slice of fields is cloned on each call to avoid data races, but this may be inefficient with many attributes (see design note 1).
func WithValues(ctx context.Context, args ...any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	n := len(args)
	if n%2 != 0 {
		n-- // Odd argument is ignored (see design note 3)
	}
	if n == 0 {
		return ctx
	}

	newFields := make([]field, 0, n/2)
	for i := 0; i+1 < n; i += 2 {
		key, ok := args[i].(string)
		if !ok || key == "" {
			continue
		}
		newFields = append(newFields, field{key: key, val: args[i+1]})
	}
	if len(newFields) == 0 {
		return ctx
	}

	p, _ := ctx.Value(contextKeyFields).(*fieldsData)
	if p == nil || len(p.slice) == 0 {
		return context.WithValue(ctx, contextKeyFields, &fieldsData{slice: newFields})
	}

	combined := append(slices.Clone(p.slice), newFields...)
	return context.WithValue(ctx, contextKeyFields, &fieldsData{slice: combined})
}

// WithUniqueValues returns a new context with the provided key-value pairs added for logging,
// replacing any existing values for the same keys.
// Only string keys are accepted and empty keys are ignored.
// If an odd number of arguments is passed, the last argument is ignored without warning (see design note 3).
// The slice of fields is cloned on each call to avoid data races, but this may be inefficient with many attributes (see design note 1).
func WithUniqueValues(ctx context.Context, args ...any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	n := len(args)
	if n%2 != 0 {
		n-- // odd value is ignored
	}
	if n == 0 {
		return ctx
	}

	// Собираем новые ключи и значения
	newFields := make([]field, 0, n/2)
	replaceSet := make(map[string]struct{}, n/2)
	for i := 0; i+1 < n; i += 2 {
		key, ok := args[i].(string)
		if !ok || key == "" {
			continue
		}
		replaceSet[key] = struct{}{}
		newFields = append(newFields, field{key: key, val: args[i+1]})
	}
	if len(newFields) == 0 {
		return ctx
	}

	// Определим базовые поля (без тех, что будут заменены)
	var existing []field
	if p, ok := ctx.Value(contextKeyFields).(*fieldsData); ok && len(p.slice) > 0 {
		existing = p.slice
	}

	// Предвыделим память под итоговый слайс
	result := make([]field, 0, len(existing)+len(newFields))

	// Добавим только старые поля с уникальными ключами
	for _, f := range existing {
		if _, replace := replaceSet[f.key]; !replace {
			result = append(result, f)
		}
	}

	// Добавим новые поля
	combined := append(result, newFields...)

	return context.WithValue(ctx, contextKeyFields, &fieldsData{slice: combined})
}

// WithoutAllKeys returns a new context with all logging attributes removed.
// This is done by setting the context value to nil for the internal key (see design note 6).
func WithoutAllKeys(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyFields, nil)
}

// WithoutKeys returns a new context with the specified keys removed from the logging attributes.
// The slice of fields is cloned on each call to avoid data races, but this may be inefficient with many attributes (see design note 1).
func WithoutKeys(ctx context.Context, keys ...string) context.Context {
	if ctx == nil || len(keys) == 0 {
		return ctx
	}

	p, ok := ctx.Value(contextKeyFields).(*fieldsData)
	if !ok || p == nil || len(p.slice) == 0 {
		return ctx
	}

	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	// Предвыделяем память для потенциально полного слайса
	result := make([]field, 0, len(p.slice))
	for _, f := range p.slice {
		if _, skip := keySet[f.key]; !skip {
			result = append(result, f)
		}
	}

	if len(result) == 0 {
		return context.WithValue(ctx, contextKeyFields, nil)
	}
	return context.WithValue(ctx, contextKeyFields, &fieldsData{slice: result})
}

// GetFirstValue retrieves the first value associated with the given key from the context attributes.
// Returns the value and true if found, or (nil, false) if the key is empty or not present.
//
// Example:
//
//	val, ok := GetFirstValue(ctx, "userID")
//	if ok { /* use val */ }
func GetFirstValue(ctx context.Context, key string) (any, bool) {
	if key == "" {
		return nil, false
	}
	p, _ := ctx.Value(contextKeyFields).(*fieldsData)
	if p == nil {
		return nil, false
	}
	for _, f := range p.slice {
		if f.key == key {
			return f.val, true
		}
	}
	return nil, false
}

// HasKey checks if the given key exists in the context attributes.
// Returns true if the key is present and not empty, otherwise false.
//
// Example:
//
//	if HasKey(ctx, "requestID") { /* ... */ }
func HasKey(ctx context.Context, key string) bool {
	if key == "" {
		return false
	}
	p, _ := ctx.Value(contextKeyFields).(*fieldsData)
	if p == nil {
		return false
	}
	for _, f := range p.slice {
		if f.key == key {
			return true
		}
	}
	return false
}
