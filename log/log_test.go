package log_test

import (
	"context"
	"os"

	"github.com/haleyrc/lib/log"
)

func Example() {
	ctx := context.Background()
	logger := log.New(
		log.FreezeTime(),
		log.Debug(),
		log.WithOutput(os.Stdout),
	)

	logger.Debug(ctx, "debug msg", "string", "Hello, World!")
	logger.Error(ctx, "error msg", "string", "Hello, World!")
	logger.Info(ctx, "info msg", "string", "Hello, World!")

	// Output:
	//
	// {"time":"2024-02-01T12:01:32-05:00","level":"DEBUG","msg":"debug msg","string":"Hello, World!"}
	// {"time":"2024-02-01T12:01:32-05:00","level":"ERROR","msg":"error msg","string":"Hello, World!"}
	// {"time":"2024-02-01T12:01:32-05:00","level":"INFO","msg":"info msg","string":"Hello, World!"}
}
