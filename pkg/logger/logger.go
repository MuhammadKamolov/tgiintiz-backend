package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

const EnvLocal = "local"

type Config struct {
	Env   string
	Level slog.Level
	File  string
}
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			_ = h.Handle(ctx, r)
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

func New(cfg Config) (*slog.Logger, func() error, error) {
	if cfg.Env != EnvLocal {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Level})
		return slog.New(handler), func() error { return nil }, nil
	}

	f, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      cfg.Level,
		TimeFormat: time.TimeOnly, // формат "15:04:05" вместо длинного ISO-8601
		NoColor:    false,
	})

	fileHandler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})

	combinedLogger := slog.New(&multiHandler{
		handlers: []slog.Handler{consoleHandler, fileHandler},
	})

	return combinedLogger, f.Close, nil
}
