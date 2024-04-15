package slogctx

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
)

type Struct struct {
	Number int64
	String string
}

func TestHandler(t *testing.T) {
	ctx := WithValues(context.Background(), "number", 12, "string", "data", "struct", Struct{
		Number: 42,
		String: "struct_data",
	})
	logger := slog.New(NewHandler(slog.NewJSONHandler(os.Stdout, nil)))

	logger.ErrorContext(ctx, "this is an error")
}

func TestHandlerConcurrent(t *testing.T) {
	ctx := WithValues(context.Background(), "number", 12, "string", "data", "struct", Struct{
		Number: 42,
		String: "struct_data",
	})
	logger := slog.New(NewHandler(slog.NewJSONHandler(io.Discard, nil)))

	for i := 0; i < 100; i++ {
		go logger.ErrorContext(ctx, "this is an error")
	}
}

func BenchmarkHandler(b *testing.B) {
	b.ReportAllocs()
	ctx := WithValues(context.Background(), "number", 12, "string", "data", "struct", Struct{
		Number: 42,
		String: "struct_data",
	})
	logger := slog.New(NewHandler(slog.NewJSONHandler(io.Discard, nil)))

	for i := 0; i < b.N; i++ {
		logger.ErrorContext(ctx, "this is an error")
	}
}
