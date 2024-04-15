package slogx

import (
	"log/slog"
	"time"
)

func Since(tm time.Time) slog.Attr {
	return slog.String("since", time.Since(tm).String())
}
