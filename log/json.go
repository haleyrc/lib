package log

import (
	"context"
	"log/slog"
	"os"
)

// JSONHandler writes log lines to stdout as one-line-per-object JSON.
type JSONHandler struct {
	slog.Handler
}

// NewJSONHandler returns a slog.Handler that writes log lines to stdout as
// one-line-per-object JSON.
func NewJSONHandler(opts ...Option) JSONHandler {
	cfg := config{
		w:          os.Stdout,
		debug:      false,
		freezeTime: false,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	level := slog.LevelInfo
	if cfg.debug {
		level = slog.LevelDebug
	}

	handler := JSONHandler{
		Handler: slog.NewJSONHandler(cfg.w, &slog.HandlerOptions{
			Level: level,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if cfg.freezeTime && a.Key == slog.TimeKey {
					a.Value = slog.StringValue("2024-10-03T12:01:32-05:00")
				}
				return a
			},
		}),
	}

	return handler
}

// Handle partially implements slog.Handler and decorates log lines with HTTP
// request info if any exists on the provided context.
func (h JSONHandler) Handle(ctx context.Context, r slog.Record) error {
	if ri := getRequestInfo(ctx); ri != nil {
		r.Add(slog.Group(
			"httpRequest",
			slog.String("id", ri.ID),
			slog.String("ip", ri.IP),
			slog.String("method", ri.Method),
			slog.String("path", ri.Path),
			slog.String("ua", ri.UserAgent),
		))
	}
	return h.Handler.Handle(ctx, r)
}
