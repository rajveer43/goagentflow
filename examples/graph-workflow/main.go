package main

import (
	"context"
	"fmt"

	"goagentflow/runtime"
)

func main() {
	graph := runtime.NewGraph(runtime.WithGraphStart("plan"))
	graph.AddNode("plan", func(_ context.Context, state *runtime.State) (string, error) {
		state.Set("plan", "research -> execute -> summarize")
		return "execute", nil
	})
	graph.AddNode("execute", func(_ context.Context, state *runtime.State) (string, error) {
		state.Output = "workflow complete"
		return "summarize", nil
	})
	graph.AddNode("summarize", func(_ context.Context, state *runtime.State) (string, error) {
		fmt.Println(state.Output)
		return "", nil
	})
	events, _ := runtime.NewGraphRunner(graph).Run(context.Background(), "build a graph")
	for event := range events {
		fmt.Println(event.Type, event.Payload)
	}
}
