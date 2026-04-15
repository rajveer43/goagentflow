package tests

import (
	"context"
	"testing"

	"github.com/rajveer43/goagentflow/runtime"
)

type testAgent struct {
	calls int
}

func (a *testAgent) Plan(_ context.Context, _ *runtime.State) (*runtime.Plan, error) {
	a.calls++
	if a.calls == 1 {
		return &runtime.Plan{Actions: []runtime.ToolCall{{Name: "echo", Args: map[string]any{"value": "hello"}}}}, nil
	}
	return &runtime.Plan{Done: true, Output: "done"}, nil
}

type echoTool struct{}

func (echoTool) Name() string { return "echo" }
func (echoTool) Description() string { return "echo tool" }
func (echoTool) ParamsSchema() map[string]any { return map[string]any{"type": "object"} }
func (echoTool) Call(_ context.Context, args map[string]any, _ runtime.StreamWriter) (any, error) {
	return args["value"], nil
}

func TestRunner(t *testing.T) {
	runner := runtime.NewRunner()
	runner.RegisterTool(echoTool{})
	events, err := runner.Run(context.Background(), &testAgent{}, "input")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	var sawCompleted bool
	for ev := range events {
		if ev.Type == runtime.RuntimeEventCompleted {
			sawCompleted = true
		}
	}
	if !sawCompleted {
		t.Fatal("expected completion event")
	}
}
