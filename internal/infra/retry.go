package infra

import (
	"context"
	"net/http"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	maxRetries             = 3
	defaultRetryWait       = 10 * time.Second
	maxRetryWait           = 60 * time.Second
	transientRetryBaseWait = 2 * time.Second
)

// RetryOnRateLimit wraps an SDK call with retry logic for transient errors.
// It retries up to maxRetries times when the Descope SDK returns a rate limit
// error or a server error (5xx), using appropriate backoff for each case.
func RetryOnRateLimit[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	for attempt := range maxRetries {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		de, retryable := isRetryableError(err)
		if !retryable {
			return result, err
		}

		wait := retryWaitDuration(de, attempt)

		tflog.Warn(ctx, "Transient Descope API error, retrying", map[string]any{
			"attempt": attempt + 1,
			"max":     maxRetries,
			"wait":    wait.String(),
			"code":    de.Code,
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

// isRetryableError checks if an error is a transient Descope API error that
// should be retried. Returns the Descope error and true for rate limit errors
// and server errors (5xx status codes).
func isRetryableError(err error) (*descope.Error, bool) {
	de := descope.AsError(err)
	if de == nil {
		return nil, false
	}
	if de.Code == descope.ErrRateLimitExceeded.Code {
		return de, true
	}
	if statusCode, ok := de.Info[descope.ErrorInfoKeys.HTTPResponseStatusCode].(int); ok {
		if statusCode >= http.StatusInternalServerError && statusCode < 600 {
			return de, true
		}
	}
	return nil, false
}

// retryWaitDuration returns the appropriate wait duration for a retryable error.
// Rate limit errors use the Retry-After header if available, falling back to
// defaultRetryWait. Server errors use exponential backoff starting at 2 seconds.
func retryWaitDuration(de *descope.Error, attempt int) time.Duration {
	if de.Code == descope.ErrRateLimitExceeded.Code {
		if retryAfter, ok := de.Info[descope.ErrorInfoKeys.RateLimitExceededRetryAfter].(int); ok && retryAfter > 0 {
			wait := time.Duration(retryAfter) * time.Second
			if wait > maxRetryWait {
				return maxRetryWait
			}
			return wait
		}
		return defaultRetryWait
	}
	// Exponential backoff for server errors: 2s, 4s, 8s, ...
	wait := transientRetryBaseWait << uint(attempt)
	if wait > maxRetryWait {
		return maxRetryWait
	}
	return wait
}
