package slogpretty

import (
	"bytes"
	"log/slog"
	"os"
)

// Example: базовое использование PrettyHandler
func ExampleNewPrettyHandler_basic() {
	h := NewPrettyHandler().
		WithColorEnabled(false).
		WithTimeLayout("").
		WithWriter(os.Stdout)
	logger := slog.New(h)
	logger.Info("Hello, world!", slog.String("user", "alice"))
	// Output:
	// INFO  Hello, world! {"user":"alice"}
}

// Example: вывод с атрибутами и группами
func ExampleNewPrettyHandler_withAttrsAndGroups() {
	h := NewPrettyHandler().
		WithColorEnabled(false).
		WithTimeLayout("").
		WithWriter(os.Stdout)

	logger := slog.New(h)
	logger.Info("Order created",
		slog.Group("order",
			slog.Int("id", 123),
			slog.String("status", "paid"),
		),
		slog.String("user", "bob"),
	)
	// Output:
	// INFO  Order created {"order":{"id":123,"status":"paid"},"user":"bob"}
}

// Example: кастомный форматтер атрибутов
func ExamplePrettyHandler_WithAttrFormatter() {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithAttrFormatter(func(m map[string]any) string {
			return "ATTRS"
		}).
		WithTimeLayout("").
		WithWriter(os.Stdout)
	logger := slog.New(h)
	logger.Info("Test", slog.String("foo", "bar"))
	os.Stdout.Write(buf.Bytes())
	// Output:
	// INFO  Test ATTRS
}
