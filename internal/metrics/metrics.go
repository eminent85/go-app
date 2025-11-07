package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics holds application metrics.
type Metrics struct {
	requestCount   uint64
	errorCount     uint64
	totalDuration  uint64 // in nanoseconds
	activeRequests int64
	startTime      time.Time
	mu             sync.RWMutex
	statusCodes    map[int]uint64
}

// New creates a new Metrics instance.
func New() *Metrics {
	return &Metrics{
		startTime:   time.Now(),
		statusCodes: make(map[int]uint64),
	}
}

// RecordRequest increments the request counter.
func (m *Metrics) RecordRequest() {
	atomic.AddUint64(&m.requestCount, 1)
	atomic.AddInt64(&m.activeRequests, 1)
}

// RecordResponse records response metrics.
func (m *Metrics) RecordResponse(statusCode int, duration time.Duration) {
	atomic.AddInt64(&m.activeRequests, -1)
	atomic.AddUint64(&m.totalDuration, uint64(duration.Nanoseconds()))

	m.mu.Lock()
	m.statusCodes[statusCode]++
	m.mu.Unlock()

	if statusCode >= 500 {
		atomic.AddUint64(&m.errorCount, 1)
	}
}

// RequestCount returns the total number of requests.
func (m *Metrics) RequestCount() uint64 {
	return atomic.LoadUint64(&m.requestCount)
}

// ErrorCount returns the total number of 5xx errors.
func (m *Metrics) ErrorCount() uint64 {
	return atomic.LoadUint64(&m.errorCount)
}

// ActiveRequests returns the current number of active requests.
func (m *Metrics) ActiveRequests() int64 {
	return atomic.LoadInt64(&m.activeRequests)
}

// AverageDuration returns the average request duration.
func (m *Metrics) AverageDuration() time.Duration {
	count := atomic.LoadUint64(&m.requestCount)
	if count == 0 {
		return 0
	}
	totalNanos := atomic.LoadUint64(&m.totalDuration)
	avgNanos := totalNanos / count

	// Check for overflow before converting to int64
	if avgNanos > uint64(1<<63-1) {
		// Return max duration if overflow would occur
		return time.Duration(1<<63 - 1)
	}

	return time.Duration(avgNanos)
}

// Uptime returns the server uptime.
func (m *Metrics) Uptime() time.Duration {
	return time.Since(m.startTime)
}

// StatusCodes returns a copy of status code counts.
func (m *Metrics) StatusCodes() map[int]uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	codes := make(map[int]uint64, len(m.statusCodes))
	for k, v := range m.statusCodes {
		codes[k] = v
	}
	return codes
}

// ErrorRate returns the error rate as a percentage.
func (m *Metrics) ErrorRate() float64 {
	requests := atomic.LoadUint64(&m.requestCount)
	if requests == 0 {
		return 0
	}
	errors := atomic.LoadUint64(&m.errorCount)
	return (float64(errors) / float64(requests)) * 100
}
