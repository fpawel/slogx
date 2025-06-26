package slogctx

import (
	"fmt"
	"github.com/fpawel/slogx/slogtest"
	"log/slog"
)

// ExampleWith demonstrates adding and removing logging attributes via context.
// First, two "req" attributes and one "user" attribute are added to the context.
// The logger outputs a message with all these attributes.
// Then, all "req" attributes are removed from the context, and the logger outputs another message.
func ExampleWith() {
	// Create a log handler that outputs logs in text format without timestamps.
	handler := slogtest.NewStdoutTextHandlerWithoutTimestamp()
	// Wrap the handler with slogctx to support context-based attributes.
	logger := slog.New(NewHandler(handler))

	// Add multiple attributes to the context: two "req" and one "user".
	ctx := WithValues(nil, "req", "abc", "req", "def", "user", "Jan")
	// Log a message; all attributes from the context will be included.
	logger.InfoContext(ctx, "Hello")
	// Remove all "req" attributes from the context.
	ctx = WithoutKeys(ctx, "req")
	// Log another message; only the "user" attribute remains.
	logger.InfoContext(ctx, "Hello")

	// Output:
	// level=INFO msg=Hello req=abc req=def user=Jan
	// level=INFO msg=Hello user=Jan
}

// ExampleWithUniqueValues demonstrates how to add unique attributes to the context.
// If a key already exists, its value is replaced with the new one.
func ExampleWithUniqueValues() {
	handler := slogtest.NewStdoutTextHandlerWithoutTimestamp()
	logger := slog.New(NewHandler(handler))

	// Add "user" and "role" attributes.
	ctx := WithValues(nil, "user", "Alice", "role", "admin")
	logger.InfoContext(ctx, "First")

	// Replace "user" value and add a new "session" attribute.
	ctx = WithUniqueValues(ctx, "user", "Bob", "session", "xyz")
	logger.InfoContext(ctx, "Second")

	// Output:
	// level=INFO msg=First user=Alice role=admin
	// level=INFO msg=Second role=admin user=Bob session=xyz
}

// ExampleWithoutAllKeys demonstrates how to remove all attributes from the context.
func ExampleWithoutAllKeys() {
	handler := slogtest.NewStdoutTextHandlerWithoutTimestamp()
	logger := slog.New(NewHandler(handler))

	// Add "id" and "status" attributes.
	ctx := WithValues(nil, "id", 123, "status", "ok")
	logger.InfoContext(ctx, "Before clear")

	// Remove all attributes from the context.
	ctx = WithoutAllKeys(ctx)
	logger.InfoContext(ctx, "After clear")

	// Output:
	// level=INFO msg="Before clear" id=123 status=ok
	// level=INFO msg="After clear"
}

// ExampleGetFirstValue demonstrates how to get the first value for a key from the context.
func ExampleGetFirstValue() {
	// Add "user" and "id" attributes.
	ctx := WithValues(nil, "user", "Eve", "id", 42)
	val, ok := GetFirstValue(ctx, "user")
	if ok {
		fmt.Println("user:", val.(string))
	}
	// Output:
	// user: Eve
}

// ExampleHasKey demonstrates how to check if a key exists in the context.
func ExampleHasKey() {
	// Add "token" attribute.
	ctx := WithValues(nil, "token", "abc123")
	if HasKey(ctx, "token") {
		fmt.Println("token exists")
	}
	// Output:
	// token exists
}

// ExampleWithValues demonstrates chaining context modifications and logging at different levels.
func ExampleWithValues() {
	handler := slogtest.NewStdoutTextHandlerWithoutTimestamp()
	logger := slog.New(NewHandler(handler))

	// Add initial attributes.
	ctx := WithValues(nil, "service", "auth", "env", "prod")
	logger.InfoContext(ctx, "Service started")

	// Add a request ID and log a warning.
	ctx = WithValues(ctx, "request_id", "req-001")
	logger.WarnContext(ctx, "Request failed")

	// Remove "env" attribute and log an error.
	ctx = WithoutKeys(ctx, "env")
	logger.ErrorContext(ctx, "Critical error")

	// Output:
	// level=INFO msg="Service started" service=auth env=prod
	// level=WARN msg="Request failed" service=auth env=prod request_id=req-001
	// level=ERROR msg="Critical error" service=auth request_id=req-001
}
