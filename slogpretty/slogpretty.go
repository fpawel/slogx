package slogpretty

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// PrettyHandler is a custom slog.Handler that provides human-friendly,
// colorized, and structured log output, optimized for local development.
//
// Features:
//   - Colorized log levels and messages (if output is a terminal)
//   - Customizable timestamp format
//   - Support for log attribute groups and attribute rewriting
//   - Fast attribute formatting (avoids JSON by default)
//   - Optional source code location in log output
//   - Efficient memory usage via sync.Pool for temporary slices
type PrettyHandler struct {
	Logger          *log.Logger     // Underlying Go logger for output
	TimestampFormat string          // Format for timestamps (default: "15:04:05")
	Level           slog.Leveler    // Minimum log level to output
	AddSource       bool            // If true, includes file and line number in logs
	BaseAttrs       []slog.Attr     // Default attributes added to every log record
	AttrGroups      []string        // Attribute group nesting for structured logs
	EnableColor     bool            // Enables color output if true
	RewriteAttrFunc RewriteAttrFunc // Optional function to rewrite attributes before output
	FormatAttrsFunc FormatAttrsFunc // Function to format attributes as a string
}

// RewriteAttrFunc allows rewriting or filtering attributes before output.
type RewriteAttrFunc func(groups []string, a slog.Attr) slog.Attr

// FormatAttrsFunc formats attributes as a string for log output.
type FormatAttrsFunc func(map[string]any) string

const DefaultTimeLayout = "15:04:05"

var (
	levelsInfo = map[slog.Level]struct {
		text      string
		colorFunc func(string, ...interface{}) string
	}{
		slog.LevelDebug: {"DEBUG", color.MagentaString},
		slog.LevelInfo:  {"INFO ", color.BlueString},
		slog.LevelWarn:  {"WARN ", color.YellowString},
		slog.LevelError: {"ERROR", color.RedString},
	}

	// partsPool is used to reuse temporary slices for log message parts.
	partsPool = sync.Pool{
		New: func() any { return make([]interface{}, 0, 8) },
	}
)

// NewPrettyHandler creates a new PrettyHandler with default settings.
//
// Color output is enabled if the output is a terminal.
func NewPrettyHandler() *PrettyHandler {
	return &PrettyHandler{
		Logger:          log.New(os.Stderr, "", 0),
		TimestampFormat: DefaultTimeLayout,
		Level:           slog.LevelDebug,
		EnableColor:     isatty.IsTerminal(os.Stderr.Fd()),
		FormatAttrsFunc: jsonAttrFormatter,
	}
}

// SetPrettyHandlerAsDefault sets PrettyHandler as the default slog handler.
func SetPrettyHandlerAsDefault() {
	slog.SetDefault(slog.New(NewPrettyHandler()))
}

// jsonAttrFormatter formats attributes as "{key1:value1 key2:value2}".
func jsonAttrFormatter(m map[string]any) string {
	if len(m) == 0 {
		return ""
	}
	b, err := json.Marshal(m)
	if err != nil {
		b, _ = json.Marshal(map[string]string{"error": fmt.Sprintf("failed to format attributes: %s", err)})
	}
	return string(b)
}

// clone returns a copy of the handler with the same settings.
func (h *PrettyHandler) clone() *PrettyHandler {
	return &PrettyHandler{
		Logger:          h.Logger,
		TimestampFormat: h.TimestampFormat,
		Level:           h.Level,
		AddSource:       h.AddSource,
		RewriteAttrFunc: h.RewriteAttrFunc,
		BaseAttrs:       append([]slog.Attr(nil), h.BaseAttrs...),
		AttrGroups:      append([]string(nil), h.AttrGroups...),
		EnableColor:     h.EnableColor,
		FormatAttrsFunc: h.FormatAttrsFunc,
	}
}

// colorize applies a color function if color is enabled.
func (h *PrettyHandler) colorize(s string, f func(string, ...interface{}) string) string {
	if h.EnableColor && f != nil {
		return f(s)
	}
	return s
}

// WithWriter returns a copy of the handler with a new output writer.
func (h *PrettyHandler) WithWriter(w io.Writer) *PrettyHandler {
	clone := h.clone()
	clone.Logger = log.New(w, "", 0)
	return clone
}

// WithTimeLayout returns a copy of the handler with a new timestamp format.
func (h *PrettyHandler) WithTimeLayout(layout string) *PrettyHandler {
	clone := h.clone()
	clone.TimestampFormat = layout
	return clone
}

// WithLogLevel returns a copy of the handler with a new minimum log level.
func (h *PrettyHandler) WithLogLevel(l slog.Leveler) *PrettyHandler {
	clone := h.clone()
	clone.Level = l
	return clone
}

// WithSourceInfo returns a copy of the handler with source info enabled or disabled.
func (h *PrettyHandler) WithSourceInfo(v bool) *PrettyHandler {
	clone := h.clone()
	clone.AddSource = v
	return clone
}

// WithAttrRewriter returns a copy of the handler with a new attribute rewriter.
func (h *PrettyHandler) WithAttrRewriter(f RewriteAttrFunc) *PrettyHandler {
	clone := h.clone()
	clone.RewriteAttrFunc = f
	return clone
}

// WithColorEnabled returns a copy of the handler with color enabled or disabled.
func (h *PrettyHandler) WithColorEnabled(enabled bool) *PrettyHandler {
	clone := h.clone()
	clone.EnableColor = enabled
	return clone
}

// WithAttrFormatter returns a copy of the handler with a new attribute formatter.
func (h *PrettyHandler) WithAttrFormatter(f FormatAttrsFunc) *PrettyHandler {
	clone := h.clone()
	clone.FormatAttrsFunc = f
	return clone
}

// WithAttrs returns a copy of the handler with additional base attributes.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := h.clone()
	clone.BaseAttrs = append(clone.BaseAttrs, attrs...)
	return clone
}

// WithGroup returns a copy of the handler with an additional attribute group.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	clone := h.clone()
	clone.AttrGroups = append(clone.AttrGroups, name)
	return clone
}

// Enabled reports whether the handler handles records at the given level.
func (h *PrettyHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l.Level() >= h.Level.Level()
}

// Handle formats and outputs a log record.
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	if h.Logger == nil {
		return fmt.Errorf("logger is not initialized")
	}

	parts := partsPool.Get().([]interface{})
	parts = parts[:0]
	defer partsPool.Put(parts[:0])

	if h.TimestampFormat != "" {
		parts = append(parts, h.colorize(r.Time.Format(h.TimestampFormat), color.WhiteString))
	}

	parts = append(parts, h.formatLevelLabel(r), h.colorize(r.Message, color.CyanString))

	if attrStr, err := h.renderAttrs(r); err == nil && attrStr != "" {
		parts = append(parts, attrStr)
	} else if err != nil {
		return err
	}

	if h.AddSource && r.PC != 0 {
		if src := formatSourceInfo(r); src != "" {
			parts = append(parts, h.colorize(src, color.GreenString))
		}
	}

	h.Logger.Println(parts...)
	return nil
}

// formatLevelLabel returns the formatted log level label, with color if enabled.
func (h *PrettyHandler) formatLevelLabel(r slog.Record) string {
	info, ok := levelsInfo[r.Level.Level()]
	label := r.Level.String()
	if ok && info.text != "" {
		label = info.text
	}
	if h.EnableColor && info.colorFunc != nil {
		return info.colorFunc(label)
	}
	return fmt.Sprintf("%-5s", label)
}

// renderAttrs collects and formats attributes for a log record.
func (h *PrettyHandler) renderAttrs(r slog.Record) (string, error) {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	attrs = append(attrs, h.BaseAttrs...)

	m := flattenAttrs(attrs, h.AttrGroups, h.RewriteAttrFunc)
	if len(m) == 0 {
		return "", nil
	}

	for i := len(h.AttrGroups) - 1; i >= 0; i-- {
		m = map[string]any{h.AttrGroups[i]: m}
	}

	return h.colorize(h.FormatAttrsFunc(m), color.WhiteString), nil
}

// flattenAttrs flattens attributes and groups into a map for formatting.
func flattenAttrs(attrs []slog.Attr, groups []string, replace func([]string, slog.Attr) slog.Attr) map[string]any {
	out := make(map[string]any, len(attrs))
	for _, a := range attrs {
		if replace != nil {
			a = replace(groups, a)
		}
		if a.Value.Kind() == slog.KindGroup {
			out[a.Key] = flattenAttrs(a.Value.Group(), append(groups, a.Key), replace)
		} else {
			out[a.Key] = a.Value.Any()
		}
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
	if i := lastDot(funcName); i >= 0 {
		funcName = funcName[i:]
	}
	return fmt.Sprintf("%s:%d%s", file, f.Line, funcName)
}

// lastDot returns the index of the last dot in a string, or -1 if not found.
func lastDot(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}
