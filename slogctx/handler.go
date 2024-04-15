package slogctx

import (
	"context"
	"log/slog"
	"sync"
)

var _ slog.Handler = Handler{}

type Handler struct {
	slog.Handler
}

func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{
		Handler: handler,
	}
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if v, ok := ctx.Value(keyFields).(*sync.Map); ok {
		v.Range(func(key, val any) bool {
			if keyString, ok := key.(string); ok {
				record.AddAttrs(slog.Any(keyString, val))
			}
			return true
		})
	}
	return h.Handler.Handle(ctx, record)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.Handler.WithAttrs(attrs)}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.Handler.WithGroup(name)
}
