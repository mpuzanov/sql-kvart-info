package main

import (
	"context"
	"kvart-info/internal/config"
	"kvart-info/internal/services"
	"kvart-info/internal/storage"
	"kvart-info/pkg/logging"
	"log/slog"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	cfg := config.MustConfig("config.yml")

	//установка логгера по умолчанию на основании конфиг файла
	logger := logging.NewLogger(cfg.Env)
	slog.SetDefault(logger)

	logger.Debug("debug", "cfg", cfg)

	// добавим логгер в контекст
	ctx := logging.ContextWithLogger(context.Background(), logger)

	db, err := storage.NewDB(ctx, cfg)
	if err != nil {
		panic(err)
	}

	if err = services.New(ctx, cfg, db).Run(); err != nil {
		panic(err)
	}

}
