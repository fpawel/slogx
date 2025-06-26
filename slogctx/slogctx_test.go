package slogctx

import (
	"context"
	"github.com/fpawel/slogx/slogtest"
	"log/slog"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func attrsMap(attrs []slog.Attr) map[string]any {
	m := map[string]any{}
	for _, a := range attrs {
		m[a.Key] = a.Value.Any()
	}
	return m
}

func TestWithValuesAndHandle(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "foo", 1, "bar", "baz")
	logger.InfoContext(ctx, "msg1")

	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(1), attrs["foo"])
	require.Equal(t, "baz", attrs["bar"])
}

func TestWithValues_StructValue(t *testing.T) {
	type S struct{ X int }
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	val := S{X: 42}
	ctx := WithValues(context.Background(), "s", val)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, val, attrs["s"])
}

func TestWithValues_DuplicateKeys(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "a", 1)
	ctx = WithValues(ctx, "a", 2)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(2), attrs["a"])
}

func TestWithValues_InvalidKeys(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "", 1, 123, "bad", "ok", 42)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(42), attrs["ok"])
	require.NotContains(t, attrs, "")
	require.NotContains(t, attrs, "bad")
}

func TestWithValues_NilContext(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(nil, "a", 1)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(1), attrs["a"])
}

func TestWithValues_OddArgs(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "a", 1, "b")
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(1), attrs["a"])
	require.NotContains(t, attrs, "b")
}

func TestWithoutValues_RemoveKey(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "a", 1, "b", 2)
	ctx = WithoutKeys(ctx, "a")
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.NotContains(t, attrs, "a")
	require.Equal(t, int64(2), attrs["b"])
}

func TestWithoutValues_AllKeys(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "a", 1)
	ctx = WithoutKeys(ctx, "a")
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Empty(t, attrs)
}

func TestWithoutValues_EmptyKeys(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	ctx := WithValues(context.Background(), "a", 1)
	WithoutKeys(ctx)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(1), attrs["a"])
}

func TestWithoutValues_NilContext(t *testing.T) {
	WithoutKeys(nil, "a")
}

func TestHandler_WithAttrsAndGroup(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	h := NewHandler(handler)
	h2 := h.WithAttrs([]slog.Attr{slog.String("x", "y")})
	h3 := h2.WithGroup("g")
	logger := slog.New(h3)
	ctx := WithValues(context.Background(), "a", 1)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, int64(1), attrs["a"])
	require.Equal(t, "y", attrs["x"])
	require.Equal(t, []string{"g"}, logs[0].Groups)
}

func TestHandler_Enabled(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	h := NewHandler(handler)
	require.True(t, h.Enabled(context.Background(), slog.LevelInfo))
}

func TestWithValues_AnyTypeValue(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	logger := slog.New(NewHandler(handler))
	type custom struct{ V string }
	val := custom{"abc"}
	ctx := WithValues(context.Background(), "custom", val)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.True(t, reflect.DeepEqual(val, attrs["custom"]))
}

func TestHandler_NestedGroupsAndAttrs(t *testing.T) {
	handler := slogtest.NewObservedHandler()
	h := NewHandler(handler)
	h2 := h.WithGroup("outer").WithAttrs([]slog.Attr{slog.String("foo", "bar")})
	h3 := h2.WithGroup("inner").WithAttrs([]slog.Attr{slog.String("baz", "qux")})
	logger := slog.New(h3)
	ctx := WithValues(context.Background(), "ctxKey", 123)
	logger.InfoContext(ctx, "msg")
	logs := handler.Logs()
	require.Len(t, logs, 1)
	attrs := attrsMap(logs[0].Attrs)
	require.Equal(t, "bar", attrs["foo"])
	require.Equal(t, "qux", attrs["baz"])
	require.Equal(t, int64(123), attrs["ctxKey"])
	require.Equal(t, []string{"outer", "inner"}, logs[0].Groups)
}

func TestGetFirstValue(t *testing.T) {
	ctx := WithValues(context.Background(), "foo", 123, "bar", "baz")
	val, ok := GetFirstValue(ctx, "foo")
	require.True(t, ok)
	require.Equal(t, 123, val)

	val, ok = GetFirstValue(ctx, "bar")
	require.True(t, ok)
	require.Equal(t, "baz", val)

	val, ok = GetFirstValue(ctx, "missing")
	require.False(t, ok)
	require.Nil(t, val)

	val, ok = GetFirstValue(ctx, "")
	require.False(t, ok)
	require.Nil(t, val)
}

func TestHasKey(t *testing.T) {
	ctx := WithValues(context.Background(), "foo", 1)
	require.True(t, HasKey(ctx, "foo"))
	require.False(t, HasKey(ctx, "bar"))
	require.False(t, HasKey(ctx, ""))
}

func TestGetFirstValue_DuplicateKeys(t *testing.T) {
	ctx := WithValues(context.Background(), "a", 1)
	ctx = WithValues(ctx, "a", 2)
	val, ok := GetFirstValue(ctx, "a")
	require.True(t, ok)
	require.Equal(t, 1, val) // returns the first occurrence
}
