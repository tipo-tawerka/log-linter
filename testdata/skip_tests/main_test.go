package main

import (
	"log/slog"
	"testing"
)

func TestSomething(_ *testing.T) {
	password := "test-pass"

	slog.Info("Test started!")
	slog.Debug("Loading данные")
	slog.Warn("check password", "p", password)
}
