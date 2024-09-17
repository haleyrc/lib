// Package log provides structured logging in which log records are JSON objects
// separated by newlines and containing, at minimum, a timestamp, severity, and
// message. Log records can be further decorated with additional attributes
// supplied as key-value pairs.
package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type config struct {
	freezeTime bool
	level      slog.Level
	output     io.Writer
}

// A Logger records structured information about each call to its Debug, Info,
// and Error methods.
//
// To create a new logger, call New with any desired Options.
type Logger struct {
	l *slog.Logger
}

// New creates a new logger that outputs log lines as one-line-per-object
// JSON.
func New(opts ...Option) *Logger {
	cfg := config{
		freezeTime: false,
		level:      slog.LevelInfo,
		output:     os.Stderr,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	logger := &Logger{
		l: slog.New(slog.NewJSONHandler(
			cfg.output,
			&slog.HandlerOptions{
				Level: cfg.level,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey && cfg.freezeTime {
						a.Value = slog.StringValue("2024-02-01T12:01:32-05:00")
					}
					return a
				},
			},
		)),
	}

	return logger
}

// Debug emits a log line at the debug level.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.l.DebugContext(ctx, msg, args...)
}

// Error emits a log line at the error level.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.l.ErrorContext(ctx, msg, args...)
}

// Info emits a log line at the info level.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.l.InfoContext(ctx, msg, args...)
}

// An Option modifies the configuration of the Logger created by calling New.
type Option func(*config)

// Debug configures a logger to output messages at the debug log level. This is
// not usually safe to use in a production environment.
func Debug() Option {
	return func(cfg *config) {
		cfg.level = slog.LevelDebug
	}
}

// FreezeTime configures a logger to output a static timestamp. This option is
// available for testing to make example output deterministic.
func FreezeTime() Option {
	return func(cfg *config) {
		cfg.freezeTime = true
	}
}

// WithOutput configures a logger to write to w.
func WithOutput(w io.Writer) Option {
	return func(cfg *config) {
		cfg.output = w
	}
}
