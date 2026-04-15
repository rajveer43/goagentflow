package backoff

import "time"

func Exponential(base time.Duration, factor float64, max time.Duration) func(int) time.Duration {
	if base <= 0 {
		base = time.Millisecond
	}
	if factor < 1 {
		factor = 1
	}
	return func(attempt int) time.Duration {
		if attempt < 1 {
			attempt = 1
		}
		delay := float64(base)
		for i := 1; i < attempt; i++ {
			delay *= factor
		}
		if d := time.Duration(delay); d > max && max > 0 {
			return max
		}
		return time.Duration(delay)
	}
}
