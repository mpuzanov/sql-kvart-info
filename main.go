package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"kvart-info/internal/config"
	"kvart-info/internal/services"
	"kvart-info/internal/storage"
	"log/slog"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}
	// TODO: установка логгера по умолчанию на основании конфиг файла
	//slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	slog.Debug("debug", "cfg", cfg)
	db, err := storage.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	if err = services.New(cfg, db).GetLicTotal(); err != nil {
		panic(err)
	}

}
