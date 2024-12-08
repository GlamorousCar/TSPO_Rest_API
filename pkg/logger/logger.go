package logger

import (
	"log/slog"
	"os"
	"strings"
)

func CreateLogger(level string) *slog.Logger {
	level = strings.ToLower(level)
	var currentLogLevel slog.Level
	switch level {
	case "debug":
		currentLogLevel = slog.LevelDebug
	case "info":
		currentLogLevel = slog.LevelInfo
	case "error":
		currentLogLevel = slog.LevelError
	}

	slogOpts := slog.HandlerOptions{
		Level: currentLogLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, &slogOpts)
	logger := slog.New(handler)
	return logger

}
