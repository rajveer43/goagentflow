package goagentflow

import (
	"context"
	"time"

	"goagentflow/internal/idempotency"
	"goagentflow/internal/stream"
)

type Runner struct {
	cfg      Config
	registry *ToolRegistry
}

type RunOption func(*runConfig)

type runConfig struct {
	maxSteps int
}

func WithRunMaxSteps(maxSteps int) RunOption {
	return func(cfg *runConfig) { cfg.maxSteps = maxSteps }
}

func NewRunner(opts ...Option) *Runner {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Runner{cfg: cfg, registry: NewToolRegistry()}
}

func (r *Runner) RegisterTool(tool Tool) {
	r.registry.Register(tool)
}

func (r *Runner) Run(ctx context.Context, agent Agent, input any, opts ...RunOption) (<-chan RuntimeEvent, error) {
	runCfg := runConfig{maxSteps: r.cfg.MaxSteps}
	for _, opt := range opts {
		opt(&runCfg)
	}
	events := stream.NewEventStream(64)
	out := make(chan RuntimeEvent, 64)
	state := NewState(input)
	go func() {
		defer close(out)
		for raw := range events.C {
			if event, ok := raw.(RuntimeEvent); ok {
				out <- event
			}
		}
	}()
	go r.execute(ctx, agent, state, runCfg, events)
	return out, nil
}

func (r *Runner) execute(ctx context.Context, agent Agent, state *State, runCfg runConfig, events *stream.EventStream) {
	defer events.Close()
	if runCfg.maxSteps <= 0 {
		runCfg.maxSteps = r.cfg.MaxSteps
	}
	if runCfg.maxSteps <= 0 {
		runCfg.maxSteps = 1
	}
	_ = idempotency.NewKey()
	for step := 0; step < runCfg.maxSteps; step++ {
		if err := ctx.Err(); err != nil {
			events.TryEmit(RuntimeEvent{Type: RuntimeEventError, Timestamp: time.Now(), Payload: ErrContextCanceled, Step: step})
			return
		}
		state.Step = step
		plan, err := agent.Plan(ctx, state)
		if err != nil {
			events.TryEmit(RuntimeEvent{Type: RuntimeEventError, Timestamp: time.Now(), Payload: err, Step: step})
			return
		}
		if r.cfg.Memory != nil {
			_ = r.cfg.Memory.AddMessage(ctx, Message{Role: "agent", Content: "plan created"})
		}
		r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventPlanCreated, Timestamp: time.Now(), Payload: plan, Step: step})
		if plan == nil {
			continue
		}
		if plan.Done {
			state.Output = plan.Output
			r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventCompleted, Timestamp: time.Now(), Payload: plan.Output, Step: step})
			return
		}
		for _, call := range plan.Actions {
			tool, ok := r.registry.Get(call.Name)
			if !ok {
				r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventToolFailed, Timestamp: time.Now(), Payload: ErrToolNotFound, Step: step})
				continue
			}
			r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventToolStarted, Timestamp: time.Now(), Payload: call, Step: step})
			var result any
			err := r.cfg.RetryPolicy.Do(ctx, func() error {
				out, callErr := tool.Call(ctx, call.Args, events)
				if callErr != nil {
					return callErr
				}
				result = out
				return nil
			})
			if err != nil {
				r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventToolFailed, Timestamp: time.Now(), Payload: err, Step: step})
				continue
			}
			state.Set(call.Name, result)
			if r.cfg.Memory != nil {
				_ = r.cfg.Memory.Set(ctx, call.Name, result)
			}
			r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventToolFinished, Timestamp: time.Now(), Payload: result, Step: step})
		}
		r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventStateUpdated, Timestamp: time.Now(), Payload: state, Step: step})
	}
	r.emit(ctx, events, RuntimeEvent{Type: RuntimeEventError, Timestamp: time.Now(), Payload: ErrMaxStepsExceeded, Step: runCfg.maxSteps})
}

func (r *Runner) emit(ctx context.Context, events *stream.EventStream, event RuntimeEvent) {
	events.TryEmit(event)
	for _, observer := range r.cfg.Observers {
		observer.Observe(ctx, event)
	}
}
