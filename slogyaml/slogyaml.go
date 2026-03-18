package slogyaml

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fpawel/slogx/internal"
	"gopkg.in/yaml.v3"
)

type YAMLHandler struct {
	Logger          *log.Logger     // Underlying Go logger for output
	TimestampFormat string          // Format for timestamps (default: "15:04:05")
	Level           slog.Leveler    // Minimum log level to output
	AddSource       bool            // If true, includes file and line number in logs
	AddLevel        bool            // If true, includes level in logs
	BaseAttrs       []slog.Attr     // Default attributes added to every log record
	AttrGroups      []string        // Attribute group nesting for structured logs
	RewriteAttrFunc RewriteAttrFunc // Optional function to rewrite attributes before output
}

// RewriteAttrFunc allows rewriting or filtering attributes before output.
type RewriteAttrFunc func(groups []string, a slog.Attr) slog.Attr

// FormatAttrsFunc formats attributes as a string for log output.
type FormatAttrsFunc func(map[string]any) string

// NewYAMLHandler creates a new YAMLHandler with default settings.
//
// Color output is enabled if the output is a terminal.
func NewYAMLHandler() *YAMLHandler {
	return &YAMLHandler{
		Logger:          log.New(os.Stderr, "", 0),
		TimestampFormat: internal.DefaultTimeLayout,
		Level:           slog.LevelDebug,
	}
}

// clone returns a copy of the handler with the same settings.
func (h *YAMLHandler) clone() *YAMLHandler {
	return &YAMLHandler{
		Logger:          h.Logger,
		TimestampFormat: h.TimestampFormat,
		Level:           h.Level,
		AddSource:       h.AddSource,
		RewriteAttrFunc: h.RewriteAttrFunc,
		BaseAttrs:       append([]slog.Attr(nil), h.BaseAttrs...),
		AttrGroups:      append([]string(nil), h.AttrGroups...),
	}
}

// WithWriter returns a copy of the handler with a new output writer.
func (h *YAMLHandler) WithWriter(w io.Writer) *YAMLHandler {
	clone := h.clone()
	clone.Logger = log.New(w, "", 0)
	return clone
}

// WithTimeLayout returns a copy of the handler with a new timestamp format.
func (h *YAMLHandler) WithTimeLayout(layout string) *YAMLHandler {
	clone := h.clone()
	clone.TimestampFormat = layout
	return clone
}

// WithLogLevel returns a copy of the handler with a new minimum log level.
func (h *YAMLHandler) WithLogLevel(l slog.Leveler) *YAMLHandler {
	clone := h.clone()
	clone.Level = l
	return clone
}

// WithSourceInfo returns a copy of the handler with source info enabled or disabled.
func (h *YAMLHandler) WithSourceInfo(v bool) *YAMLHandler {
	clone := h.clone()
	clone.AddSource = v
	return clone
}

// WithLevel returns a copy of the handler with level enabled or disabled.
func (h *YAMLHandler) WithLevel(v bool) *YAMLHandler {
	clone := h.clone()
	clone.AddLevel = v
	return clone
}

// WithAttrRewriter returns a copy of the handler with a new attribute rewriter.
func (h *YAMLHandler) WithAttrRewriter(f RewriteAttrFunc) *YAMLHandler {
	clone := h.clone()
	clone.RewriteAttrFunc = f
	return clone
}

// WithAttrs returns a copy of the handler with additional base attributes.
func (h *YAMLHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := h.clone()
	clone.BaseAttrs = append(clone.BaseAttrs, attrs...)
	return clone
}

// WithGroup returns a copy of the handler with an additional attribute group.
func (h *YAMLHandler) WithGroup(name string) slog.Handler {
	clone := h.clone()
	clone.AttrGroups = append(clone.AttrGroups, name)
	return clone
}

// Enabled reports whether the handler handles records at the given level.
func (h *YAMLHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l.Level() >= h.Level.Level()
}

// Handle formats and outputs a log record.
func (h *YAMLHandler) Handle(_ context.Context, r slog.Record) error {
	if h.Logger == nil {
		return fmt.Errorf("logger is not initialized")
	}

	var data []any

	if h.TimestampFormat != "" {
		data = append(data, r.Time.Format(h.TimestampFormat))
	}
	if h.AddLevel {
		data = append(data, internal.FormatLevelLabel(r, false))
	}
	if h.AddSource && r.PC != 0 {
		if src := formatSourceInfo(r); src != "" {
			data = append(data, src)
		}
	}
	data = append(data, h.collectAttrs(r)...)
	data = []any{map[string]any{r.Message: data}}
	raw, err := yaml.Marshal(data)
	if err != nil {
		raw, _ = yaml.Marshal(map[string]string{"error": fmt.Sprintf("failed to format message: %s", err)})
	}

	h.Logger.Printf("%s", raw)

	return nil
}

// renderAttrs collects and formats attributes for a log record.
func (h *YAMLHandler) collectAttrs(r slog.Record) []any {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	attrs = append(attrs, h.BaseAttrs...)
	m := flattenAttrs(attrs, h.AttrGroups, h.RewriteAttrFunc)
	for i := len(h.AttrGroups) - 1; i >= 0; i-- {
		m = []any{map[string]any{h.AttrGroups[i]: m}}
	}
	return m
}

// flattenAttrs flattens attributes and groups into a map for formatting.
func flattenAttrs(attrs []slog.Attr, groups []string, replace func([]string, slog.Attr) slog.Attr) []any {
	out := make([]any, 0, len(attrs))
	for _, a := range attrs {
		if replace != nil {
			a = replace(groups, a)
		}
		v := a.Value.Any()
		if a.Value.Kind() == slog.KindGroup {
			v = flattenAttrs(a.Value.Group(), append(groups, a.Key), replace)
		}
		out = append(out, map[string]any{a.Key: v})
	}
	return out
}

// formatSourceInfo returns a string with file, line, and function name for the log record.
func formatSourceInfo(r slog.Record) string {
	if r.PC == 0 {
		return ""
	}
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	funcName, file := f.Function, f.File
	if f.Function == "" {
		funcName = "unknown-function"
	}
	if file == "" {
		file = "unknown-file"
	}
	funcName = filepath.Base(funcName)
	file = filepath.Base(file)
	if i := internal.LastDot(funcName); i >= 0 {
		funcName = funcName[i:]
	}
	return fmt.Sprintf("%s:%d%s", file, f.Line, funcName)
}
