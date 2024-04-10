package logging

import (
	"log/slog"
	"os"
)

const (
	LocalMode Mode = "local"
	DevMode   Mode = "dev"
	ProdMode  Mode = "prod"
)

type (
	Mode   string
	Logger = slog.Logger
)

func getLocalLogger() *Logger {
	return slog.New(newPrettyHandler(&slog.HandlerOptions{Level: slog.LevelDebug}))
}

func getDevLogger() *Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func getProdLogger() *Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func GetLogger(mode Mode) *Logger {
	var log *Logger

	switch mode {
	case LocalMode:
		log = getLocalLogger()
	case DevMode:
		log = getDevLogger()
	case ProdMode:
		log = getProdLogger()
	}

	return log
}
