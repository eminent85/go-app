package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eminent85/go-app/internal/metrics"
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hello", nil)
	w := httptest.NewRecorder()

	HelloHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "Hello, World!" {
		t.Errorf("Expected message 'Hello, World!', got '%s'", response["message"])
	}
}

func TestNotFoundHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	NotFoundHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "Resource not found" {
		t.Errorf("Expected error 'Resource not found', got '%s'", response["error"])
	}
}

func TestMetricsHandler(t *testing.T) {
	m := metrics.New()

	// Simulate some requests
	m.RecordRequest()
	m.RecordResponse(200, 100)
	m.RecordRequest()
	m.RecordResponse(500, 200)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	handler := MetricsHandler(m)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response MetricsResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.TotalRequests != 2 {
		t.Errorf("Expected total requests 2, got %d", response.TotalRequests)
	}

	if response.ErrorCount != 1 {
		t.Errorf("Expected error count 1, got %d", response.ErrorCount)
	}

	if response.StatusCodes[200] != 1 {
		t.Errorf("Expected 1 request with status 200, got %d", response.StatusCodes[200])
	}

	if response.StatusCodes[500] != 1 {
		t.Errorf("Expected 1 request with status 500, got %d", response.StatusCodes[500])
	}
}
