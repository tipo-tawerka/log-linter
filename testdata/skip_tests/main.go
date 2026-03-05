package main

import "log/slog"

func logRequest(id string) {
	slog.Info("handling request", "id", id)
	slog.Debug("request processed")
}

func main() {
	logRequest("req-1")
}
