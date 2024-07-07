package src

import (
	"log/slog"
	"os"
)

var (
	mylogger *slog.Logger
)

func Init(level *string, format *string) {
	var l_level slog.Level
	switch *level {
	case "debug":
		l_level = slog.LevelDebug
	case "warn":
		l_level = slog.LevelWarn
	case "error":
		l_level = slog.LevelError
	default:
		l_level = slog.LevelInfo
	}

	// Select format
	var handler slog.Handler
	switch *format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l_level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l_level})
	}
	mylogger = slog.New(handler)
}

func GetLogger() *slog.Logger {
	return mylogger
}
