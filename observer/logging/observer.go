package logging

import (
	"context"
	"log/slog"

	"github.com/rajveer43/goagentflow/runtime"
)

// LoggingObserver implements runtime.Observer and emits structured logs for each event.
// Pattern: Observer - logs all runtime events
type LoggingObserver struct {
	logger *slog.Logger
}

// New creates a new logging observer.
// logger: slog logger instance (uses default if nil)
func New(logger *slog.Logger) *LoggingObserver {
	if logger == nil {
		logger = slog.Default()
	}
	return &LoggingObserver{logger: logger}
}

// Observe logs a runtime event with structured fields.
func (l *LoggingObserver) Observe(ctx context.Context, event runtime.RuntimeEvent) {
	// Extract payload details based on event type
	attrs := []slog.Attr{
		slog.String("event_type", string(event.Type)),
		slog.Time("timestamp", event.Timestamp),
		slog.Int("step", event.Step),
		slog.String("trace_id", event.TraceID),
	}

	// Add type-specific attributes based on payload
	switch event.Type {
	case runtime.RuntimeEventPlanCreated:
		if plan, ok := event.Payload.(*runtime.Plan); ok {
			attrs = append(attrs,
				slog.Int("action_count", len(plan.Actions)),
				slog.Bool("done", plan.Done),
			)
		}

	case runtime.RuntimeEventToolStarted:
		if toolCall, ok := event.Payload.(runtime.ToolCall); ok {
			attrs = append(attrs, slog.String("tool_name", toolCall.Name))
		}

	case runtime.RuntimeEventToolFinished:
		if toolCall, ok := event.Payload.(runtime.ToolCall); ok {
			attrs = append(attrs, slog.String("tool_name", toolCall.Name))
		}

	case runtime.RuntimeEventToolFailed:
		if err, ok := event.Payload.(error); ok {
			attrs = append(attrs, slog.String("error", err.Error()))
		}

	case runtime.RuntimeEventStateUpdated:
		if state, ok := event.Payload.(*runtime.State); ok {
			attrs = append(attrs, slog.Any("state_values", state.Values))
		}

	case runtime.RuntimeEventError:
		if err, ok := event.Payload.(error); ok {
			attrs = append(attrs, slog.String("error", err.Error()))
		}
	}

	// Emit structured log
	l.logger.LogAttrs(ctx, slog.LevelInfo, "runtime_event", attrs...)
}
