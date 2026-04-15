package metrics

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rajveer43/goagentflow/runtime"
)

// MetricsObserver tracks runtime metrics using atomic operations and histograms.
// Pattern: Observer - collects metrics for all runtime events
// DSA: Atomic counters for O(1) increments, histogram buckets as sorted slice
type MetricsObserver struct {
	// Atomic counters (O(1) increments)
	toolCalls      int64
	toolSuccesses  int64
	toolFailures   int64
	totalErrors    int64
	totalSteps     int64
	runsCompleted  int64

	mu sync.RWMutex
	// Tool-specific metrics
	toolLatencies map[string][]time.Duration // histogram of latencies per tool

	// Event tracking for latency computation
	eventTimestamps map[string]time.Time // event ID -> start time
}

// New creates a new metrics observer.
func New() *MetricsObserver {
	return &MetricsObserver{
		toolLatencies:   make(map[string][]time.Duration),
		eventTimestamps: make(map[string]time.Time),
	}
}

// Observe records a runtime event and updates metrics.
func (m *MetricsObserver) Observe(ctx context.Context, event runtime.RuntimeEvent) {
	eventID := buildEventID(event)

	switch event.Type {
	case runtime.RuntimeEventToolStarted:
		m.mu.Lock()
		m.eventTimestamps[eventID] = time.Now()
		m.mu.Unlock()
		atomic.AddInt64(&m.toolCalls, 1)

	case runtime.RuntimeEventToolFinished:
		atomic.AddInt64(&m.toolSuccesses, 1)
		m.recordToolLatency(eventID, event)

	case runtime.RuntimeEventToolFailed:
		atomic.AddInt64(&m.toolFailures, 1)
		m.recordToolLatency(eventID, event)

	case runtime.RuntimeEventError:
		atomic.AddInt64(&m.totalErrors, 1)

	case runtime.RuntimeEventStateUpdated:
		atomic.AddInt64(&m.totalSteps, 1)

	case runtime.RuntimeEventCompleted:
		atomic.AddInt64(&m.runsCompleted, 1)
	}
}

// recordToolLatency computes and stores tool execution latency.
// DSA: Histogram buckets stored as sorted slice for percentile queries.
func (m *MetricsObserver) recordToolLatency(eventID string, event runtime.RuntimeEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	startTime, found := m.eventTimestamps[eventID]
	if !found {
		return
	}

	latency := time.Since(startTime)
	delete(m.eventTimestamps, eventID)

	// Extract tool name from event
	toolName := "unknown"
	if toolCall, ok := event.Payload.(runtime.ToolCall); ok {
		toolName = toolCall.Name
	}

	m.toolLatencies[toolName] = append(m.toolLatencies[toolName], latency)
}

// buildEventID creates a unique identifier for an event.
func buildEventID(event runtime.RuntimeEvent) string {
	return event.TraceID + "-" + string(event.Type) + "-" + string(rune(event.Step))
}

// Snapshot represents a point-in-time view of all metrics.
type Snapshot struct {
	ToolCalls       int64
	ToolSuccesses   int64
	ToolFailures    int64
	TotalErrors     int64
	TotalSteps      int64
	RunsCompleted   int64
	ToolLatencies   map[string]LatencyStats
	SuccessRate     float64
}

// LatencyStats represents percentile statistics for tool latencies.
type LatencyStats struct {
	Count   int
	Min     time.Duration
	Max     time.Duration
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
	Average time.Duration
}

// GetSnapshot returns a read-only snapshot of current metrics.
func (m *MetricsObserver) GetSnapshot() Snapshot {
	snapshot := Snapshot{
		ToolCalls:     atomic.LoadInt64(&m.toolCalls),
		ToolSuccesses: atomic.LoadInt64(&m.toolSuccesses),
		ToolFailures:  atomic.LoadInt64(&m.toolFailures),
		TotalErrors:   atomic.LoadInt64(&m.totalErrors),
		TotalSteps:    atomic.LoadInt64(&m.totalSteps),
		RunsCompleted: atomic.LoadInt64(&m.runsCompleted),
		ToolLatencies: make(map[string]LatencyStats),
	}

	// Compute success rate
	if snapshot.ToolCalls > 0 {
		snapshot.SuccessRate = float64(snapshot.ToolSuccesses) / float64(snapshot.ToolCalls)
	}

	// Compute latency statistics for each tool (DSA: percentile calculation)
	m.mu.RLock()
	defer m.mu.RUnlock()

	for toolName, latencies := range m.toolLatencies {
		if len(latencies) == 0 {
			continue
		}

		// Sort latencies for percentile calculation
		sorted := make([]time.Duration, len(latencies))
		copy(sorted, latencies)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

		stats := LatencyStats{
			Count: len(sorted),
			Min:   sorted[0],
			Max:   sorted[len(sorted)-1],
		}

		// Calculate percentiles
		stats.P50 = percentile(sorted, 0.50)
		stats.P95 = percentile(sorted, 0.95)
		stats.P99 = percentile(sorted, 0.99)

		// Calculate average
		var sum time.Duration
		for _, lat := range sorted {
			sum += lat
		}
		stats.Average = sum / time.Duration(len(sorted))

		snapshot.ToolLatencies[toolName] = stats
	}

	return snapshot
}

// percentile calculates the percentile value from a sorted slice.
// DSA: Binary search would apply here, but for percentiles, linear interpolation is standard.
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	index := float64(len(sorted)-1) * p
	lower := int(index)
	upper := lower + 1

	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}

	fraction := index - float64(lower)
	return time.Duration(float64(sorted[lower]) + fraction*float64(sorted[upper]-sorted[lower]))
}

// Reset clears all metrics (useful for testing or explicit reset).
func (m *MetricsObserver) Reset() {
	atomic.StoreInt64(&m.toolCalls, 0)
	atomic.StoreInt64(&m.toolSuccesses, 0)
	atomic.StoreInt64(&m.toolFailures, 0)
	atomic.StoreInt64(&m.totalErrors, 0)
	atomic.StoreInt64(&m.totalSteps, 0)
	atomic.StoreInt64(&m.runsCompleted, 0)

	m.mu.Lock()
	m.toolLatencies = make(map[string][]time.Duration)
	m.eventTimestamps = make(map[string]time.Time)
	m.mu.Unlock()
}
