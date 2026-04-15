package goagentflow

import "context"

type Chain interface {
	Run(ctx context.Context, input any) (any, error)
}

type ChainFunc func(ctx context.Context, input any) (any, error)

func (f ChainFunc) Run(ctx context.Context, input any) (any, error) {
	return f(ctx, input)
}

type ChainPipeline struct {
	steps []Chain
}

func NewChainPipeline(steps ...Chain) *ChainPipeline {
	return &ChainPipeline{steps: append([]Chain(nil), steps...)}
}

func (p *ChainPipeline) Run(ctx context.Context, input any) (any, error) {
	current := input
	var err error
	for _, step := range p.steps {
		current, err = step.Run(ctx, current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

