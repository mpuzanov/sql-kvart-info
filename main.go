package main

import (
	"io"
	"kvart-info/internal/config"
	"kvart-info/internal/services"
	"kvart-info/internal/storage"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	cfg := config.MustConfig("config.yml")

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
		// добавим запись логов в файл
		ex, _ := os.Executable()
		fileName := filepath.Base(ex)
		fileLog := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".log"
		workDir := filepath.Dir(ex) // путь к программе
		fileLog = filepath.Join(workDir, fileLog)

		file, _ := os.OpenFile(fileLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		multi := io.MultiWriter(file, os.Stdout) //, os.Stderr

		log = slog.New(slog.NewJSONHandler(multi, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	log = log.With(slog.String("env", env))
	return log
}
