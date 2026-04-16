package runtime

import "context"

// Memory is the interface all memory backends must satisfy.
// Pattern: Repository - interchangeable storage backends
type Memory interface {
	AddMessage(ctx context.Context, msg Message) error
	GetMessages(ctx context.Context) ([]Message, error)
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
}

// Message represents a single turn in a conversation.
type Message struct {
	Role    string // e.g., "user", "assistant", "system"
	Content string
}

// Summarizer is an optional interface for memory backends that can compress conversations.
// Pattern: Strategy - different summarization strategies can be swapped
type Summarizer interface {
	Summarize(ctx context.Context, messages []Message) (string, error)
}

// Compressor is an optional interface for memory backends that can intelligently compress messages.
// Different from Summarizer: Compressor returns compressed Message slice, not a string summary.
// Pattern: Strategy - different compression strategies can be swapped
type Compressor interface {
	Compress(ctx context.Context, messages []Message) ([]Message, error)
}
