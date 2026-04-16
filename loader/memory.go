package loader

import (
	"context"
	"github.com/rajveer43/goagentflow/runtime"
)

// InjectIntoMemory stores loaded documents in runtime.Memory.
// If asMessages is true, each document is also added as a system message.
// Pattern: Dependency Injection
func InjectIntoMemory(
	ctx context.Context,
	mem runtime.Memory,
	key string,
	docs []Document,
	asMessages bool,
) error {
	// Store documents under key
	if err := mem.Set(ctx, key, docs); err != nil {
		return err
	}

	// Optionally add each doc as a system message
	if asMessages {
		for _, doc := range docs {
			msg := runtime.Message{
				Role:    "system",
				Content: doc.PageContent,
			}
			if err := mem.AddMessage(ctx, msg); err != nil {
				return err
			}
		}
	}

	return nil
}

// InjectIntoState stores loaded documents in runtime.State under the given key.
// Pattern: Dependency Injection
func InjectIntoState(state *runtime.State, key string, docs []Document) {
	state.Set(key, docs)
}
