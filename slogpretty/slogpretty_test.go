package slogpretty

import (
	"bytes"
	"context"
	"log/slog"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPrettyHandler_Levels(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithTimeLayout("")

	logger := slog.New(h.WithLogLevel(slog.LevelWarn))
	logger.Info("info should not appear")
	logger.Warn("warn should appear")
	logger.Error("error should appear")

	out := buf.String()
	require.NotContains(t, out, "info should not appear")
	require.Contains(t, out, "warn should appear")
	require.Contains(t, out, "ERROR")
}

func TestPrettyHandler_WithSourceInfo(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithSourceInfo(true).
		WithTimeLayout("")

	logger := slog.New(h)
	logger.Info("msg")

	out := buf.String()
	require.Contains(t, out, "msg")
	require.Regexp(t, `\.go:\d+`, out)
	t.Log(out)
}

func TestPrettyHandler_WithAttrFormatter(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithAttrFormatter(func(m map[string]any) string {
			return "ATTRS"
		}).
		WithTimeLayout("")

	logger := slog.New(h)
	logger.Info("test", slog.String("foo", "bar"))

	out := buf.String()
	require.Contains(t, out, "ATTRS")
}

func TestPrettyHandler_WithAttrRewriter(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithAttrRewriter(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "secret" {
				return slog.Attr{}
			}
			return a
		}).
		WithTimeLayout("")

	logger := slog.New(h)
	logger.Info("msg", slog.String("secret", "xxx"), slog.String("visible", "ok"))

	out := buf.String()
	require.NotContains(t, out, "secret")
	require.Contains(t, out, "visible")
}

func TestPrettyHandler_WithAttrsAndGroups(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithTimeLayout("")

	logger := slog.New(h)
	logger = logger.With(slog.String("base", "b"))
	logger = logger.WithGroup("grp")
	logger.Info("msg", slog.Int("id", 1))

	out := buf.String()

	require.Contains(t, out, `"grp":`)
	require.Contains(t, out, `"id":1`)
	require.Contains(t, out, `"base":"b"`)
}

func TestPrettyHandler_WithWriter(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithTimeLayout("")

	logger := slog.New(h)
	logger.Info("msg")
	require.Contains(t, buf.String(), "msg")
}

func TestPrettyHandler_LoggerNil(t *testing.T) {
	h := NewPrettyHandler()
	h.Logger = nil
	err := h.Handle(context.Background(), slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0))
	require.Error(t, err)
}

func TestPrettyHandler_Colorize(t *testing.T) {
	h := NewPrettyHandler().WithColorEnabled(true)
	s := h.colorize("test", nil)
	require.Equal(t, "test", s)
}

func TestPrettyHandler_WithTimeLayout(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithTimeLayout("2006")

	logger := slog.New(h)
	logger.Info("msg")
	require.Contains(t, buf.String(), time.Now().Format("2006"))
}

func getCallerPC() (uintptr, string, int, bool) {
	return getCallerPCDepth(1)
}

func getCallerPCDepth(depth int) (uintptr, string, int, bool) {
	pc, file, line, ok := getCaller(depth + 1)
	return pc, file, line, ok
}

func getCaller(depth int) (uintptr, string, int, bool) {
	var pcs [1]uintptr
	n := runtime.Callers(depth+2, pcs[:])
	if n == 0 {
		return 0, "", 0, false
	}
	frames := runtime.CallersFrames(pcs[:])
	f, _ := frames.Next()
	return pcs[0], f.File, f.Line, true
}

func Test_formatSourceInfo(t *testing.T) {
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	require.Equal(t, "", formatSourceInfo(r))

	// Сымитировать PC
	var pc uintptr
	pc, _, _, _ = getCallerPC()
	r = slog.NewRecord(time.Now(), slog.LevelInfo, "msg", pc)
	info := formatSourceInfo(r)
	require.Regexp(t, `\.go:\d+`, info)
}

func TestPrettyHandler_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	h := NewPrettyHandler().
		WithWriter(&buf).
		WithColorEnabled(false).
		WithTimeLayout("").
		WithLogLevel(slog.LevelWarn)
	logger := slog.New(h)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	out := buf.String()
	require.NotContains(t, out, "debug")
	require.NotContains(t, out, "info")
	require.Contains(t, out, "WARN")
	require.Contains(t, out, "ERROR")
}
