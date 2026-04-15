package tests

import (
	"context"
	"testing"

	"goagentflow/runtime"
)

type registryTool struct{}

func (registryTool) Name() string { return "registry" }
func (registryTool) Description() string { return "registry tool" }
func (registryTool) ParamsSchema() map[string]any { return map[string]any{} }
func (registryTool) Call(_ context.Context, _ map[string]any, _ runtime.StreamWriter) (any, error) {
	return "ok", nil
}

func TestToolRegistry(t *testing.T) {
	registry := runtime.NewToolRegistry()
	registry.Register(registryTool{})
	if _, ok := registry.Get("registry"); !ok {
		t.Fatal("expected tool in registry")
	}
	if len(registry.List()) != 1 {
		t.Fatal("expected one tool")
	}
}
