package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"kvart-info/internal/config"
	"kvart-info/internal/services"
	"kvart-info/internal/storage"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustConfig("config.json")

	//установка логгера по умолчанию на основании конфиг файла
	logger := NewLogger(cfg.Env)
	slog.SetDefault(logger)

	logger.Debug("debug", "cfg", cfg)
	db, err := storage.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	if err = services.New(cfg, db).Run(); err != nil {
		panic(err)
	}

}

// NewLogger ...
func NewLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
