package loggermw

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client := r.RemoteAddr
			start := time.Now()

			next.ServeHTTP(w, r)

			since := time.Since(start)

			log.Info(
				fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				slog.String("client", client),
				slog.String("latency", since.String()),
			)
		})
	}
}
