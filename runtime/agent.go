package runtime

import "context"

type Agent interface {
	Plan(ctx context.Context, state *State) (*Plan, error)
}

type Observer interface {
	Observe(ctx context.Context, event RuntimeEvent)
}
