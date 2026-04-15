package runtime

import "errors"

var (
	ErrToolNotFound     = errors.New("tool not found")
	ErrMaxStepsExceeded = errors.New("max steps exceeded")
	ErrContextCanceled  = errors.New("context canceled")
)
