package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New().String()

			logger := base.With(
				slog.String("request_id", requestId),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			ctx := context.WithValue(r.Context(), "logger", logger)

			start := time.Now()
			next.ServeHTTP(w, r.WithContext(ctx))
			logger.Info("request finished", slog.String("duration", time.Since(start).String()))

		})
	}
}
