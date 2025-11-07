package metrics

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("Expected non-nil metrics")
	}

	if m.RequestCount() != 0 {
		t.Errorf("Expected initial request count 0, got %d", m.RequestCount())
	}

	if m.ErrorCount() != 0 {
		t.Errorf("Expected initial error count 0, got %d", m.ErrorCount())
	}
}

func TestRecordRequest(t *testing.T) {
	m := New()
	m.RecordRequest()

	if count := m.RequestCount(); count != 1 {
		t.Errorf("Expected request count 1, got %d", count)
	}

	if active := m.ActiveRequests(); active != 1 {
		t.Errorf("Expected active requests 1, got %d", active)
	}
}

func TestRecordResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		duration   time.Duration
		shouldErr  bool
	}{
		{"success response", 200, 100 * time.Millisecond, false},
		{"client error", 404, 50 * time.Millisecond, false},
		{"server error", 500, 200 * time.Millisecond, true},
		{"another server error", 503, 150 * time.Millisecond, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.RecordRequest()
			initialErrors := m.ErrorCount()

			m.RecordResponse(tt.statusCode, tt.duration)

			if active := m.ActiveRequests(); active != 0 {
				t.Errorf("Expected active requests 0 after response, got %d", active)
			}

			codes := m.StatusCodes()
			if codes[tt.statusCode] != 1 {
				t.Errorf("Expected status code %d count 1, got %d", tt.statusCode, codes[tt.statusCode])
			}

			if tt.shouldErr {
				if m.ErrorCount() != initialErrors+1 {
					t.Errorf("Expected error count to increment for status %d", tt.statusCode)
				}
			}
		})
	}
}

func TestAverageDuration(t *testing.T) {
	m := New()

	// No requests yet
	if avg := m.AverageDuration(); avg != 0 {
		t.Errorf("Expected average duration 0 for no requests, got %v", avg)
	}

	// Add some requests
	m.RecordRequest()
	m.RecordResponse(200, 100*time.Millisecond)

	m.RecordRequest()
	m.RecordResponse(200, 200*time.Millisecond)

	// Average should be 150ms
	avg := m.AverageDuration()
	expected := 150 * time.Millisecond
	tolerance := 1 * time.Millisecond

	if avg < expected-tolerance || avg > expected+tolerance {
		t.Errorf("Expected average duration around %v, got %v", expected, avg)
	}
}

func TestAverageDurationOverflow(t *testing.T) {
	m := New()

	// Simulate a scenario with very large total duration
	// Set totalDuration to a value that would overflow int64 if not handled
	m.RecordRequest()
	// Directly set a very large totalDuration using atomic operation
	// This simulates extreme edge case where average would overflow
	atomic.StoreUint64(&m.totalDuration, uint64(1<<63))

	// Should not panic and should return max duration
	avg := m.AverageDuration()

	// Should return max int64 duration
	maxDuration := time.Duration(1<<63 - 1)
	if avg != maxDuration {
		t.Errorf("Expected max duration %v on overflow, got %v", maxDuration, avg)
	}
}

func TestErrorRate(t *testing.T) {
	m := New()

	// No requests
	if rate := m.ErrorRate(); rate != 0 {
		t.Errorf("Expected error rate 0 for no requests, got %f", rate)
	}

	// Add requests: 2 success, 1 error
	m.RecordRequest()
	m.RecordResponse(200, 10*time.Millisecond)

	m.RecordRequest()
	m.RecordResponse(200, 10*time.Millisecond)

	m.RecordRequest()
	m.RecordResponse(500, 10*time.Millisecond)

	// Error rate should be 33.33%
	rate := m.ErrorRate()
	expected := 33.33
	tolerance := 0.1

	if rate < expected-tolerance || rate > expected+tolerance {
		t.Errorf("Expected error rate around %f%%, got %f%%", expected, rate)
	}
}

func TestUptime(t *testing.T) {
	m := New()
	time.Sleep(10 * time.Millisecond)

	uptime := m.Uptime()
	if uptime < 10*time.Millisecond {
		t.Errorf("Expected uptime >= 10ms, got %v", uptime)
	}
}

func TestStatusCodes(t *testing.T) {
	m := New()

	m.RecordRequest()
	m.RecordResponse(200, 10*time.Millisecond)

	m.RecordRequest()
	m.RecordResponse(200, 10*time.Millisecond)

	m.RecordRequest()
	m.RecordResponse(404, 10*time.Millisecond)

	codes := m.StatusCodes()

	if codes[200] != 2 {
		t.Errorf("Expected 2 requests with status 200, got %d", codes[200])
	}

	if codes[404] != 1 {
		t.Errorf("Expected 1 request with status 404, got %d", codes[404])
	}
}
