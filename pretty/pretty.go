package pretty

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

type (
	Handler struct {
		SlogOpts
		Logger     *log.Logger
		TimeLayout string // by default, do not display the time locally
		Attrs      []Attr
		Groups     []string
	}
	Record      = slog.Record
	Attr        = slog.Attr
	SlogOpts    = slog.HandlerOptions
	SlogHandler = slog.Handler
	marshalFunc func(interface{}) ([]byte, error)
	Level       = slog.Level
	levelInfo   struct {
		text      string
		colorFunc func(format string, a ...interface{}) string
	}
)

var (
	levelsInfo = map[Level]levelInfo{
		slog.LevelDebug: {
			"DEBUG", color.MagentaString},
		slog.LevelInfo: {
			"INFO ", color.BlueString},
		slog.LevelWarn: {
			"WARN ", color.YellowString},
		slog.LevelError: {
			"ERROR", color.RedString},
	}
)

// SetAsSlogDefault sets default global slog Logger with Handler with default settings for local development
func SetAsSlogDefault() {
	slog.SetDefault(slog.New(NewHandler()))
}

// NewHandler creates default Handler with default settings for local development
func NewHandler() Handler {
	return Handler{
		Logger:     log.New(os.Stderr, "", 0),
		TimeLayout: "15:04:05",
		SlogOpts: SlogOpts{
			Level:     slog.LevelDebug,
			AddSource: false,
		},
	}
}

func (h Handler) WithOutput(output io.Writer) Handler {
	h.Logger = log.New(output, "", 0)
	return h
}

func (h Handler) WithTimeLayout(layout string) Handler {
	h.TimeLayout = layout
	return h
}

func (h Handler) WithLevel(l Level) Handler {
	h.SlogOpts.Level = l
	return h
}

func (h Handler) WithSlogOpts(o SlogOpts) Handler {
	h.SlogOpts = o
	return h
}

func (h Handler) WithAddSource(v bool) Handler {
	h.SlogOpts.AddSource = v
	return h
}

func (h Handler) WithReplaceAttr(replaceAttr func(groups []string, a Attr) Attr) Handler {
	h.SlogOpts.ReplaceAttr = replaceAttr
	return h
}

func (h Handler) Handle(_ context.Context, r Record) error {
	var outputParts []interface{}
	if h.TimeLayout != "" {
		outputParts = append(outputParts, color.WhiteString(r.Time.Format(h.TimeLayout)))
	}

	outputParts = append(outputParts, h.recordLevel(r), color.CyanString(r.Message))

	strAttrs, err := h.recordAttrs(r)
	if err != nil {
		return err
	}
	if strAttrs != "" {
		outputParts = append(outputParts, strAttrs)
	}

	if h.SlogOpts.AddSource {
		outputParts = append(outputParts, color.GreenString(recordFormatSource(r)))
	}

	h.Logger.Println(outputParts...)

	return nil
}

func (h Handler) Enabled(_ context.Context, l Level) bool {
	x, y := l.Level(), h.SlogOpts.Level.Level()
	f := x >= y
	return f
}

func (h Handler) WithAttrs(attrs []Attr) SlogHandler {
	h.Attrs = append(h.Attrs, attrs...)
	return h
}

func (h Handler) WithGroup(name string) SlogHandler {
	h.Groups = append(h.Groups, name)
	return h
}

func (h Handler) recordAttrs(r Record) (string, error) {
	xs := attrsValues(append(recordAttrs(r), h.Attrs...)...)
	if len(xs) == 0 {
		return "", nil
	}
	for i := len(h.Groups) - 1; i >= 0; i-- {
		xs = map[string]interface{}{
			h.Groups[i]: xs,
		}
	}
	s, err := json.Marshal(xs)
	if err != nil {
		return "", err
	}
	return color.WhiteString(string(s)), nil
}

func (h Handler) recordLevel(r Record) string {
	l := levelsInfo[r.Level.Level()]
	level := l.text
	if level == "" {
		level = r.Level.String()
	}
	if l.colorFunc != nil {
		level = l.colorFunc(level)
	}
	return level
}

// formats a Source for the log event.
func recordFormatSource(r Record) string {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	function := filepath.Base(f.Function)
	for i, ch := range function {
		if string(ch) == "." {
			function = function[i:]
			break
		}
	}
	return fmt.Sprintf("%s:%d%s", filepath.Base(f.File), f.Line, function)
}

func recordAttrs(r Record) []Attr {
	xs := make([]Attr, 0, r.NumAttrs())
	r.Attrs(func(a Attr) bool {
		xs = append(xs, a)
		return true
	})
	return xs
}

func attrsValues(attrs ...Attr) map[string]interface{} {
	fields := make(map[string]interface{}, len(attrs))
	for _, a := range attrs {
		if a.Value.Kind() == slog.KindGroup {
			fields[a.Key] = attrsValues(a.Value.Group()...)
		} else {
			fields[a.Key] = a.Value.Any()
		}
	}
	return fields
}
