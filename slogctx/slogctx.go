package slogctx

import (
	"context"
	"log/slog"
	"sync"
)

type (
	ctxFieldsKey string
	ctxLoggerKey string
)

const (
	keyFields ctxFieldsKey = "slog_fields"
	keyLogger ctxLoggerKey = "slog_logger"
)

func WithLog(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, logger)
}

func Log(ctx context.Context) *slog.Logger {
	if v, ok := ctx.Value(keyFields).(*slog.Logger); ok {
		return v
	}
	return slog.Default()
}

func WithValues(ctx context.Context, args ...interface{}) context.Context {
	if ctx == nil {
		panic("cannot create context from nil parent")
	}
	v, ok := ctx.Value(keyFields).(*sync.Map)
	if !ok {
		v = new(sync.Map)
	}
	for i := 0; i < len(args); i += 2 {
		v.Store(args[i], args[i+1])
	}
	return context.WithValue(ctx, keyFields, v)
}
