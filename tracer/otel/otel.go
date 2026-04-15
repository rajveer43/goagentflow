package otel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer wraps OpenTelemetry trace.Tracer and implements runtime.Tracer interface.
// Pattern: Adapter - bridges runtime.Tracer to OTel trace.Tracer
type Tracer struct {
	tracer trace.Tracer
	once   sync.Once
}

// New creates a new OTel tracer wrapper.
// tracerName: semantic name for the tracer (e.g., "goagentflow")
func New(tracerName string) *Tracer {
	t := &Tracer{}
	t.once.Do(func() {
		t.tracer = otel.Tracer(tracerName)
	})
	return t
}

// StartSpan starts an OTel span and returns a function to end it.
// Implements runtime.Tracer interface (StartSpan(name string) func()).
// Note: runtime.Tracer interface doesn't pass context, so we use a background context.
// In production, consider extending the interface to accept context.
func (t *Tracer) StartSpan(name string) func() {
	// Use a background context since the interface doesn't provide one
	// This is a limitation of the current runtime.Tracer interface
	ctx, span := t.tracer.Start(context.Background(), name)
	_ = ctx // context available if needed

	// Return a closure to end the span
	return func() {
		span.End()
	}
}

// Direct returns the underlying OTel tracer for advanced usage.
// Allows direct access to context-aware tracing when needed.
func (t *Tracer) Direct() trace.Tracer {
	return t.tracer
}
