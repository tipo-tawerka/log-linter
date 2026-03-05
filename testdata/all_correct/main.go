package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	slog.Info("user logged in")
	slog.Debug("fetching data from cache")
	slog.Warn("connection retry attempt")
	slog.Error("failed to open file")

	username := "alice"
	slog.Info("processing request", "user", username)

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("server started")
	logger.Debug("cache hit", zap.String("key", "config"))
	logger.Warn("high memory usage")
	logger.Error("db query failed")
}
