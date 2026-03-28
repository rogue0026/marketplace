package interceptors

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func Logging(l *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = context.WithValue(ctx, "logger", l)
		start := time.Now()
		resp, err := handler(ctx, req)

		l.Info("request finished", slog.String("duration", time.Since(start).String()))

		return resp, err
	}
}
