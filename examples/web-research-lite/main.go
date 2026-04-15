package main

import (
	"context"
	"fmt"

	"github.com/rajveer43/goagentflow/runtime"
	"github.com/rajveer43/goagentflow/memory/inmemory"
)

type researchAgent struct {
	step int
}

func (a *researchAgent) Plan(_ context.Context, state *runtime.State) (*runtime.Plan, error) {
	a.step++
	if a.step == 1 {
		return &runtime.Plan{Actions: []runtime.ToolCall{{Name: "search", Args: map[string]any{"query": state.Input}}}}, nil
	}
	return &runtime.Plan{Done: true, Output: "research complete"}, nil
}

type searchTool struct{}

func (searchTool) Name() string { return "search" }
func (searchTool) Description() string { return "search stub" }
func (searchTool) ParamsSchema() map[string]any { return map[string]any{"type": "object"} }
func (searchTool) Call(_ context.Context, args map[string]any, stream runtime.StreamWriter) (any, error) {
	_ = stream.Write("streamed search result")
	return fmt.Sprintf("results for %v", args["query"]), nil
}

func main() {
	runner := runtime.NewRunner(runtime.WithMemory(inmemory.New()))
	runner.RegisterTool(searchTool{})
	events, _ := runner.Run(context.Background(), &researchAgent{}, "golang agents")
	for event := range events {
		fmt.Println(event.Type, event.Payload)
	}
}
