package goagentflow

import "context"

type StreamWriter interface {
	Write(event any) error
}

type Tool interface {
	Name() string
	Description() string
	ParamsSchema() map[string]any
	Call(ctx context.Context, args map[string]any, stream StreamWriter) (any, error)
}

type ToolCall struct {
	Name string
	Args map[string]any
}
