package middleware

import (
	"gateway/pkg/logger"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type logKey struct{}

func LogginMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New()

			l := base.With(
				"request_id", requestId,
				"method", r.Method,
				"path", r.URL.Path,
			)

			ctx := logger.WithContext(r.Context(), l)

			start := time.Now()

			next.ServeHTTP(w, r.WithContext(ctx))

			l.Info("request completed", "duration", time.Since(start))
		})
	}
}
