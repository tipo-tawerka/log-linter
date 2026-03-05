package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	password := "s3cr3t"
	apikey := "key-xyz"

	slog.Info("User connected!")

	slog.Debug("соединение установлено")

	slog.Warn("Ошибка соединения")

	slog.Error("auth failed", "password", password, "key", apikey)

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Auth attempt #1", zap.String("apikey", apikey))

	slog.Info("done", "status", "ok")
}
