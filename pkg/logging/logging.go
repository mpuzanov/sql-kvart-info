package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// NewLogger ...
func NewLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
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

type ctxLogger struct{}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// LoggerFromContext returns logger from context
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return l
	}
	return NewLogger(envProd)
}
