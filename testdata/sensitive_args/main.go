package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	password := "hunter2"
	secret := "top-secret-value"
	apikey := "sk-abc123"
	userToken := "eyJhbGciOiJIUzI1NiJ9"
	username := "alice"

	slog.Info("auth attempt", "user", username, "password", password)
	slog.Debug("loaded config", "secret", secret, "key", apikey)
	slog.Info("session created", "token", userToken)

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("api call", zap.String("apikey", apikey))
	logger.Info("user auth", zap.String("pass", password))

	requestID := "req-999"
	slog.Info("handling request", "id", requestID)
}
