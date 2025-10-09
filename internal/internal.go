package internal

import (
	"fmt"
	"log/slog"

	"github.com/fatih/color"
)

const DefaultTimeLayout = "15:04:05"

var LevelsInfo = map[slog.Level]struct {
	Text      string
	ColorFunc func(string, ...interface{}) string
}{
	slog.LevelDebug: {"DEBUG", color.MagentaString},
	slog.LevelInfo:  {"INFO ", color.BlueString},
	slog.LevelWarn:  {"WARN ", color.YellowString},
	slog.LevelError: {"ERROR", color.RedString},
}

// LastDot returns the index of the last dot in a string, or -1 if not found.
func LastDot(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return i
		}
	}
	return -1
}

// FormatLevelLabel returns the formatted log level label, with color if enabled.
func FormatLevelLabel(r slog.Record, enableColor bool) string {
	info, ok := LevelsInfo[r.Level.Level()]
	label := r.Level.String()
	if ok && info.Text != "" {
		label = info.Text
	}
	if enableColor && info.ColorFunc != nil {
		return info.ColorFunc(label)
	}
	return fmt.Sprintf("%-5s", label)
}
