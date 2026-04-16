package loader

import "context"

// LoaderChain wraps a Loader as a runtime.Chain.
// Pattern: Adapter — bridges Loader to runtime.Chain interface
// This lets loaders be the first step in a runtime.ChainPipeline.
type LoaderChain struct {
	L Loader
}

// Run implements runtime.Chain interface.
// Input is ignored; output is []Document.
func (lc LoaderChain) Run(ctx context.Context, _ any) (any, error) {
	return lc.L.Load(ctx)
}

// NewLoaderChain creates a new LoaderChain.
func NewLoaderChain(loader Loader) LoaderChain {
	return LoaderChain{L: loader}
}
