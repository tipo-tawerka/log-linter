package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	slog.Info("user logged in!")
	slog.Debug("cache hit: key found")
	slog.Warn("retrying... please wait")
	slog.Error("fatal error #42 occurred")

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("server started @ port 8080")
	logger.Warn("queue overflow (>1000 items)")

	slog.Info("all good here 42")
}
