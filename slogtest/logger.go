package slogtest

import (
	"log/slog"
	"testing"
)

// NewTestLogger creates a test logger with an in-memory handler.
// It returns a *slog.Logger and the associated *ObservedHandler for inspection.
// Typically used in unit tests for verifying log output.
//
// Example:
//
//	logger, observed := slogtest.NewTestLogger(t)
//	logger.Info("hello", slog.String("key", "value"))
//	logs := observed.Logs()
func NewTestLogger(t *testing.T) (*slog.Logger, *ObservedHandler) {
	t.Helper()
	h := NewObservedHandler()
	return slog.New(h), h
}
