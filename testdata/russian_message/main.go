package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	slog.Info("пользователь вошёл в систему")
	slog.Debug("загрузка данных из кэша")
	slog.Warn("повторная попытка подключения")
	slog.Error("не удалось открыть файл")

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("сервер запущен")
	logger.Error("ошибка запроса к базе данных")

	slog.Info("latin text is fine")
}
