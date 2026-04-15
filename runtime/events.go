package goagentflow

import "time"

type RuntimeEventType string

const (
	RuntimeEventPlanCreated    RuntimeEventType = "plan_created"
	RuntimeEventToolStarted    RuntimeEventType = "tool_started"
	RuntimeEventToolFinished   RuntimeEventType = "tool_finished"
	RuntimeEventToolFailed     RuntimeEventType = "tool_failed"
	RuntimeEventStateUpdated   RuntimeEventType = "state_updated"
	RuntimeEventObservationMade RuntimeEventType = "observation_made"
	RuntimeEventCompleted      RuntimeEventType = "completed"
	RuntimeEventError          RuntimeEventType = "error"
)

type RuntimeEvent struct {
	Type      RuntimeEventType
	Timestamp time.Time
	Payload   any
	Step      int
	TraceID   string
}

type Plan struct {
	Actions []ToolCall
	Done    bool
	Output  any
}
