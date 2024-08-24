package middleware

import (
	"context"
	"net/http"

	"route256/cart/pkg/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infow(context.Background(), "incoming request",
			"method", r.Method,
			"host", r.Host,
			"path", r.URL.Path,
			"user_agent", r.UserAgent(),
			"proto", r.Proto,
			"service", "cart",
		)

		next.ServeHTTP(w, r)
	})
}
