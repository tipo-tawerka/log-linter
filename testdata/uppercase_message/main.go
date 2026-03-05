package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	slog.Info("User logged in")
	slog.Debug("Fetching data from cache")
	slog.Warn("Connection retry attempt")
	slog.Error("Failed to open file")

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Server started")
	logger.Error("Database connection lost")

	slog.Info("everything is fine")
}
