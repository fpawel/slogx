package slogtest

import (
	"log/slog"
	"testing"
)

func TestLogger_WithAttrsAndGroups(t *testing.T) {
	logger, observed := NewTestLogger(t)

	logger = logger.With(slog.String("scope", "unit"))
	logger = logger.WithGroup("grp")

	logger.Info("something happened", slog.Int("code", 123))

	logs := observed.Logs()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	got := logs[0]
	if got.Message != "something happened" {
		t.Errorf("wrong message: %q", got.Message)
	}

	if len(got.Attrs) != 2 {
		t.Errorf("expected 2 attributes, got: %+v", got.Attrs)
	}

	if len(got.Groups) != 1 || got.Groups[0] != "grp" {
		t.Errorf("expected group 'grp', got: %+v", got.Groups)
	}
}
