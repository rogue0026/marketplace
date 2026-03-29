package logger

import (
	"context"
	"log/slog"
	"os"
)

func New() *slog.Logger {

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	return slog.New(h)
}

func Extract(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value("logger").(*slog.Logger)
	if ok {
		return l
	}

	return slog.Default()
}
