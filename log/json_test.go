package log_test

import (
	"context"
	"log/slog"

	"github.com/haleyrc/lib/log"
)

func Example() {
	h := log.NewJSONHandler(
		log.FreezeTime(),
		log.Debug(),
	)
	logger := slog.New(h)

	logger.Debug("hello", slog.String("target", "world"))
	logger.Info("hello", slog.String("target", "world"))
	logger.Warn("hello", slog.String("target", "world"))
	logger.Error("hello", slog.String("target", "world"))

	ctx := log.SetRequestInfo(context.Background(), &log.RequestInfo{
		ID:        "46f7e63a-5c89-42e5-a488-9809dbf47760",
		IP:        "192.168.1.100",
		Method:    "GET",
		Path:      "/test",
		UserAgent: "GoTest 1.1",
	})
	logger.DebugContext(ctx, "hello", slog.String("target", "world"))
	logger.InfoContext(ctx, "hello", slog.String("target", "world"))
	logger.WarnContext(ctx, "hello", slog.String("target", "world"))
	logger.ErrorContext(ctx, "hello", slog.String("target", "world"))

	// Output:
	//
	// {"time":"2024-10-03T12:01:32-05:00","level":"DEBUG","msg":"hello","target":"world"}
	// {"time":"2024-10-03T12:01:32-05:00","level":"INFO","msg":"hello","target":"world"}
	// {"time":"2024-10-03T12:01:32-05:00","level":"WARN","msg":"hello","target":"world"}
	// {"time":"2024-10-03T12:01:32-05:00","level":"ERROR","msg":"hello","target":"world"}
	// {"time":"2024-10-03T12:01:32-05:00","level":"DEBUG","msg":"hello","target":"world","httpRequest":{"id":"46f7e63a-5c89-42e5-a488-9809dbf47760","ip":"192.168.1.100","method":"GET","path":"/test","ua":"GoTest 1.1"}}
	// {"time":"2024-10-03T12:01:32-05:00","level":"INFO","msg":"hello","target":"world","httpRequest":{"id":"46f7e63a-5c89-42e5-a488-9809dbf47760","ip":"192.168.1.100","method":"GET","path":"/test","ua":"GoTest 1.1"}}
	// {"time":"2024-10-03T12:01:32-05:00","level":"WARN","msg":"hello","target":"world","httpRequest":{"id":"46f7e63a-5c89-42e5-a488-9809dbf47760","ip":"192.168.1.100","method":"GET","path":"/test","ua":"GoTest 1.1"}}
	// {"time":"2024-10-03T12:01:32-05:00","level":"ERROR","msg":"hello","target":"world","httpRequest":{"id":"46f7e63a-5c89-42e5-a488-9809dbf47760","ip":"192.168.1.100","method":"GET","path":"/test","ua":"GoTest 1.1"}}
}
