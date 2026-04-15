package tests

import (
	"context"
	"testing"

	"goagentflow/memory/inmemory"
)

func TestMemory(t *testing.T) {
	mem := inmemory.New()
	if err := mem.Set(context.Background(), "k", "v"); err != nil {
		t.Fatal(err)
	}
	value, err := mem.Get(context.Background(), "k")
	if err != nil || value != "v" {
		t.Fatalf("got %v %v", value, err)
	}
}
