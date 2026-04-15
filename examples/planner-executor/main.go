package main

import (
	"context"
	"fmt"

	"goagentflow/runtime"
	"goagentflow/memory/inmemory"
)

type planner struct {
	step int
}

func (p *planner) Plan(_ context.Context, state *runtime.State) (*runtime.Plan, error) {
	p.step++
	if p.step == 1 {
		return &runtime.Plan{Actions: []runtime.ToolCall{{Name: "remember", Args: map[string]any{"key": "topic", "value": state.Input}}}}, nil
	}
	return &runtime.Plan{Done: true, Output: "planned and executed"}, nil
}

type rememberTool struct{}

func (rememberTool) Name() string { return "remember" }
func (rememberTool) Description() string { return "store a value" }
func (rememberTool) ParamsSchema() map[string]any { return map[string]any{"type": "object"} }
func (rememberTool) Call(ctx context.Context, args map[string]any, _ runtime.StreamWriter) (any, error) {
	return args, ctx.Err()
}

func main() {
	runner := runtime.NewRunner(runtime.WithMemory(inmemory.New()))
	runner.RegisterTool(rememberTool{})
	events, _ := runner.Run(context.Background(), &planner{}, "build a plan")
	for event := range events {
		fmt.Printf("%s %v\n", event.Type, event.Payload)
	}
}
