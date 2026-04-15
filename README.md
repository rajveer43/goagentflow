# goagentflow

`goagentflow` is an idiomatic Go agent runtime for struct-based agents, tool execution, streaming events, retries, pluggable memory, and graph workflows.

## Why this exists

Most agent frameworks lean on dynamic typing and heavy abstractions. This package keeps the core small and explicit:

- agents are plain Go structs
- control flow is driven by `context.Context`
- tools are simple interfaces
- streaming uses channels
- dependencies arrive through options

## Architecture

- `runtime/`: core interfaces, runner, events, retry policy, and state
- `memory/inmemory`: concurrency-safe memory implementation
- `graph.go` / `graph_runner.go`: graph workflow engine with conditional transitions
- `chain.go`: composable sequential pipelines
- `provider/openai`: minimal adapter stub for OpenAI-compatible usage
- `tracer/noop`: no-op tracer
- `internal/`: backoff, idempotency, and stream helpers

## Quickstart

```go
runner := runtime.NewRunner()
runner.RegisterTool(myTool{})
events, err := runner.Run(ctx, myAgent, input)
```

Consume the returned event channel until it closes. The runner emits plan, tool, state, completion, and error events in order.

## Notes

- Go 1.22+
- Core has no network dependency
- Examples are runnable starting points for extension
