package runtime

import "context"

type Memory interface {
	AddMessage(ctx context.Context, msg Message) error
	GetMessages(ctx context.Context) ([]Message, error)
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
}

type Message struct {
	Role    string
	Content string
}
