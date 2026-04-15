package goagentflow

import (
	"context"
	"time"

	"goagentflow/internal/backoff"
)

type RetryPolicy struct {
	MaxAttempts int
	Backoff     func(int) time.Duration
}

func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 3,
		Backoff:     backoff.Exponential(50*time.Millisecond, 2, 500*time.Millisecond),
	}
}

func (p RetryPolicy) Do(ctx context.Context, fn func() error) error {
	attempts := p.MaxAttempts
	if attempts <= 0 {
		attempts = 1
	}
	for attempt := 1; attempt <= attempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		if attempt == attempts {
			return err
		}
		delay := time.Duration(0)
		if p.Backoff != nil {
			delay = p.Backoff(attempt)
		}
		if delay > 0 {
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return nil
}
