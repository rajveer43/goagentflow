package runtime

import "context"

type GraphRunner struct {
	graph *Graph
}

func NewGraphRunner(graph *Graph) *GraphRunner {
	return &GraphRunner{graph: graph}
}

func (r *GraphRunner) Run(ctx context.Context, input any) (<-chan RuntimeEvent, error) {
	events := make(chan RuntimeEvent, 64)
	state := NewState(input)
	sink := &channelSink{events: events}
	go func() {
		defer close(events)
		_ = r.graph.Run(ctx, state, sink)
	}()
	return events, nil
}

type channelSink struct {
	events chan RuntimeEvent
}

func (s *channelSink) Emit(event RuntimeEvent) {
	select {
	case s.events <- event:
	default:
		s.events <- event
	}
}

