package slogtest

import (
	"context"
	"log/slog"
	"sync"
)

// ObservedLog represents a single captured log entry.
// It includes the log level, message, attributes, and any active groups.
type ObservedLog struct {
	Level   slog.Level  // The log level (e.g., InfoLevel, ErrorLevel)
	Message string      // The log message
	Attrs   []slog.Attr // Structured attributes associated with the log entry
	Groups  []string    // Group hierarchy applied via WithGroup
}

// observedState holds the shared state (log entries and mutex) for all handlers.
// It is shared among handler clones to avoid copying sync.Mutex.
type observedState struct {
	mu   sync.Mutex
	logs []ObservedLog
}

// ObservedHandler is a custom slog.Handler implementation used for testing.
// It records all log entries in memory and allows inspection after the test.
type ObservedHandler struct {
	state  *observedState // Shared state between cloned handlers
	attrs  []slog.Attr    // Scoped attributes (via WithAttrs)
	groups []string       // Current group hierarchy (via WithGroup)
}

// NewObservedHandler returns a new instance of ObservedHandler.
// Use this handler to capture logs for testing purposes.
func NewObservedHandler() *ObservedHandler {
	return &ObservedHandler{
		state: &observedState{},
	}
}

// Enabled always returns true, allowing all log levels.
func (h *ObservedHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle captures a log record and appends it to the in-memory log buffer.
func (h *ObservedHandler) Handle(_ context.Context, r slog.Record) error {
	var recordAttrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		recordAttrs = append(recordAttrs, a)
		return true
	})

	combinedAttrs := append(append([]slog.Attr{}, h.attrs...), recordAttrs...)

	h.state.mu.Lock()
	defer h.state.mu.Unlock()
	h.state.logs = append(h.state.logs, ObservedLog{
		Level:   r.Level,
		Message: r.Message,
		Attrs:   combinedAttrs,
		Groups:  append([]string{}, h.groups...),
	})

	return nil
}

// WithAttrs returns a new handler with additional attributes applied to every record.
func (h *ObservedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	combined := append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &ObservedHandler{
		state:  h.state,
		attrs:  combined,
		groups: h.groups,
	}
}

// WithGroup returns a new handler with an additional group name added.
// Group names are used to namespace attributes in structured logs.
func (h *ObservedHandler) WithGroup(name string) slog.Handler {
	newGroups := append([]string{}, h.groups...)
	newGroups = append(newGroups, name)
	return &ObservedHandler{
		state:  h.state,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

// Logs returns a copy of all captured logs.
func (h *ObservedHandler) Logs() []ObservedLog {
	h.state.mu.Lock()
	defer h.state.mu.Unlock()
	return append([]ObservedLog(nil), h.state.logs...)
}
