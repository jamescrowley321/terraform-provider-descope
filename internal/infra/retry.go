package infra

import (
	"context"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RetryOnRateLimit wraps an SDK call with rate-limit retry logic.
// It retries up to maxRetries times when the Descope SDK returns a rate
// limit error, respecting the Retry-After duration from the error info.
func RetryOnRateLimit[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	for attempt := range maxRetries {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		de := descope.AsError(err, descope.ErrRateLimitExceeded.Code)
		if de == nil {
			return result, err
		}

		wait := defaultRetryWait
		if retryAfter, ok := de.Info[descope.ErrorInfoKeys.RateLimitExceededRetryAfter].(int); ok && retryAfter > 0 {
			wait = time.Duration(retryAfter) * time.Second
		}
		if wait > maxRetryWait {
			wait = maxRetryWait
		}

		tflog.Warn(ctx, "Rate limited by Descope API, retrying", map[string]any{
			"attempt": attempt + 1,
			"max":     maxRetries,
			"wait":    wait.String(),
		})

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}

	return fn()
}

// RetryOnRateLimitNoResult wraps an SDK call that returns only an error.
func RetryOnRateLimitNoResult(ctx context.Context, fn func() error) error {
	_, err := RetryOnRateLimit(ctx, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	return err
}
