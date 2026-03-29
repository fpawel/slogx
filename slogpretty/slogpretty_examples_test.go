package slogpretty

import (
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
	logger.Info("Hello, world!", slog.String("user", "alice"), "user", "bob")
	// Output:
	// INFO  Hello, world! {"user":["alice","bob"]}
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
