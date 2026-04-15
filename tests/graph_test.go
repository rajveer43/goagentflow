package tests

import (
	"context"
	"testing"

	"github.com/rajveer43/goagentflow/runtime"
)

func TestGraphExecution(t *testing.T) {
	graph := runtime.NewGraph(runtime.WithGraphStart("start"))
	graph.AddNode("start", func(_ context.Context, state *runtime.State) (string, error) {
		state.Output = "done"
		return "end", nil
	})
	graph.AddNode("end", func(_ context.Context, _ *runtime.State) (string, error) {
		return "", nil
	})
	events := make([]runtime.RuntimeEvent, 0, 4)
	if err := graph.Run(context.Background(), runtime.NewState("input"), eventCollector{events: &events}); err != nil {
		t.Fatalf("graph run: %v", err)
	}
	if len(events) == 0 {
		t.Fatal("expected graph events")
	}
}

type eventCollector struct {
	events *[]runtime.RuntimeEvent
}

func (c eventCollector) Emit(event runtime.RuntimeEvent) {
	*c.events = append(*c.events, event)
}
