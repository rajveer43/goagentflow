package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rajveer43/goagentflow/runtime"
)

func TestRetry(t *testing.T) {
	policy := runtime.RetryPolicy{MaxAttempts: 3, Backoff: func(int) time.Duration { return 0 }}
	count := 0
	err := policy.Do(context.Background(), func() error {
		count++
		if count < 3 {
			return errors.New("fail")
		}
		return nil
	})
	if err != nil || count != 3 {
		t.Fatalf("retry failed: %v count=%d", err, count)
	}
}
