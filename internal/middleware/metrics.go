package middleware

import (
	"net/http"
	"time"

	"github.com/yourusername/go-app/internal/metrics"
)

// Metrics middleware tracks request metrics
func Metrics(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			m.RecordRequest()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			m.RecordResponse(wrapped.statusCode, duration)
		})
	}
}
