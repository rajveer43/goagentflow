package tests

import (
	"context"
	"testing"

	"goagentflow/runtime"
)

type streamingAgent struct {
	called bool
}

func (a *streamingAgent) Plan(_ context.Context, _ *runtime.State) (*runtime.Plan, error) {
	if a.called {
		return &runtime.Plan{Done: true, Output: "ok"}, nil
	}
	a.called = true
	return &runtime.Plan{Actions: []runtime.ToolCall{{Name: "stream", Args: map[string]any{"value": "x"}}}}, nil
}

type streamingTool struct{}

func (streamingTool) Name() string { return "stream" }
func (streamingTool) Description() string { return "stream tool" }
func (streamingTool) ParamsSchema() map[string]any { return map[string]any{} }
func (streamingTool) Call(_ context.Context, _ map[string]any, stream runtime.StreamWriter) (any, error) {
	_ = stream.Write("chunk-1")
	return "done", nil
}

func TestStreamingOrdering(t *testing.T) {
	runner := runtime.NewRunner()
	runner.RegisterTool(streamingTool{})
	events, err := runner.Run(context.Background(), &streamingAgent{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var order []runtime.RuntimeEventType
	for event := range events {
		order = append(order, event.Type)
	}
	if len(order) == 0 || order[0] != runtime.RuntimeEventPlanCreated {
		t.Fatalf("unexpected event ordering: %v", order)
	}
}
