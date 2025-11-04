package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/go-app/internal/metrics"
)

// MetricsResponse represents the metrics endpoint response
type MetricsResponse struct {
	TotalRequests   uint64            `json:"total_requests"`
	ActiveRequests  int64             `json:"active_requests"`
	ErrorCount      uint64            `json:"error_count"`
	ErrorRate       float64           `json:"error_rate_percent"`
	AverageDuration string            `json:"average_duration"`
	Uptime          string            `json:"uptime"`
	StatusCodes     map[int]uint64    `json:"status_codes"`
}

// MetricsHandler returns metrics data
func MetricsHandler(m *metrics.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := MetricsResponse{
			TotalRequests:   m.RequestCount(),
			ActiveRequests:  m.ActiveRequests(),
			ErrorCount:      m.ErrorCount(),
			ErrorRate:       m.ErrorRate(),
			AverageDuration: m.AverageDuration().String(),
			Uptime:          m.Uptime().String(),
			StatusCodes:     m.StatusCodes(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// HelloHandler is a simple example endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Hello, World!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"error": "Resource not found",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}
