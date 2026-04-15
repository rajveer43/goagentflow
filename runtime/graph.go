package goagentflow

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type NodeFunc func(ctx context.Context, state *State) (string, error)

type Graph struct {
	nodes map[string]NodeFunc
	edges map[string][]Edge
	start string
	end   string
}

type Edge struct {
	To       string
	When     func(*State) bool
	Priority int
}

type GraphOption func(*Graph)

func WithGraphStart(node string) GraphOption {
	return func(g *Graph) { g.start = node }
}

func WithGraphEnd(node string) GraphOption {
	return func(g *Graph) { g.end = node }
}

func NewGraph(opts ...GraphOption) *Graph {
	g := &Graph{
		nodes: make(map[string]NodeFunc),
		edges: make(map[string][]Edge),
		end:   "__end__",
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *Graph) AddNode(name string, fn NodeFunc) {
	g.nodes[name] = fn
}

func (g *Graph) AddEdge(from, to string, when func(*State) bool) {
	g.edges[from] = append(g.edges[from], Edge{To: to, When: when})
}

func (g *Graph) Run(ctx context.Context, state *State, sink EventSink) error {
	current := g.start
	if current == "" {
		return errors.New("graph start node not set")
	}
	for steps := 0; steps < 1024; steps++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		fn, ok := g.nodes[current]
		if !ok {
			return fmt.Errorf("graph node %q not found", current)
		}
		if sink != nil {
			sink.Emit(RuntimeEvent{Type: RuntimeEventPlanCreated, Timestamp: time.Now(), Step: state.Step, Payload: current})
		}
		next, err := fn(ctx, state)
		if err != nil {
			if sink != nil {
				sink.Emit(RuntimeEvent{Type: RuntimeEventError, Timestamp: time.Now(), Step: state.Step, Payload: err})
			}
			return err
		}
		if next == "" {
			next = g.nextNode(current, state)
		}
		if next == "" || next == g.end {
			if sink != nil {
				sink.Emit(RuntimeEvent{Type: RuntimeEventCompleted, Timestamp: time.Now(), Step: state.Step, Payload: state.Output})
			}
			return nil
		}
		if sink != nil {
			sink.Emit(RuntimeEvent{Type: RuntimeEventStateUpdated, Timestamp: time.Now(), Step: state.Step, Payload: next})
		}
		current = next
		state.Step++
	}
	return ErrMaxStepsExceeded
}

func (g *Graph) nextNode(from string, state *State) string {
	candidates := append([]Edge(nil), g.edges[from]...)
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].Priority > candidates[i].Priority {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	for _, edge := range candidates {
		if edge.When == nil || edge.When(state) {
			return edge.To
		}
	}
	return ""
}

type EventSink interface {
	Emit(RuntimeEvent)
}

