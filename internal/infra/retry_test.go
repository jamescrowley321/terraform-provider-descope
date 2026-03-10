package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/descope/go-sdk/descope"
)

func TestRetryOnRateLimit(t *testing.T) {
	t.Run("returns result without retrying on success", func(t *testing.T) {
		calls := 0
		result, err := RetryOnRateLimit(context.Background(), func() (string, error) {
			calls++
			return "ok", nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "ok" {
			t.Fatalf("expected 'ok', got %q", result)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	t.Run("does not retry non-rate-limit errors", func(t *testing.T) {
		calls := 0
		_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
			calls++
			return "", errors.New("some other error")
		})
		if err == nil || err.Error() != "some other error" {
			t.Fatalf("expected 'some other error', got %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	t.Run("aborts retry wait and returns context error when cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		calls := 0
		rateLimitErr := descope.ErrRateLimitExceeded.WithMessage("rate limited")
		_, err := RetryOnRateLimit(ctx, func() (string, error) {
			calls++
			return "", rateLimitErr
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("does not retry descope errors with non-rate-limit codes", func(t *testing.T) {
		calls := 0
		descopeErr := &descope.Error{Code: "E999999", Message: "bad request"}
		_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
			calls++
			return "", descopeErr
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	t.Run("retries on rate limit then succeeds", func(t *testing.T) {
		calls := 0
		start := time.Now()
		result, err := RetryOnRateLimit(context.Background(), func() (string, error) {
			calls++
			if calls <= 2 {
				return "", descope.ErrRateLimitExceeded.
					WithMessage("rate limited").
					WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, 1)
			}
			return "recovered", nil
		})
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "recovered" {
			t.Fatalf("expected 'recovered', got %q", result)
		}
		if calls != 3 {
			t.Fatalf("expected 3 calls (2 retries + success), got %d", calls)
		}
		// Should have waited ~2 seconds (1s per retry × 2 retries)
		if elapsed < 1*time.Second {
			t.Fatalf("expected at least 1s of retry wait, got %v", elapsed)
		}
	})

	t.Run("exhausts retries and makes final call", func(t *testing.T) {
		calls := 0
		rateLimitErr := descope.ErrRateLimitExceeded.
			WithMessage("rate limited").
			WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, 1)
		_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
			calls++
			return "", rateLimitErr
		})
		if err == nil {
			t.Fatal("expected error after exhausting retries")
		}
		// maxRetries (3) loop iterations + 1 final call = 4 total
		if calls != maxRetries+1 {
			t.Fatalf("expected %d calls (maxRetries + final), got %d", maxRetries+1, calls)
		}
	})

	t.Run("uses Retry-After from error info", func(t *testing.T) {
		calls := 0
		start := time.Now()
		_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
			calls++
			if calls == 1 {
				return "", descope.ErrRateLimitExceeded.
					WithMessage("rate limited").
					WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, 2)
			}
			return "ok", nil
		})
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should have waited ~2 seconds from Retry-After
		if elapsed < 2*time.Second {
			t.Fatalf("expected at least 2s retry wait from Retry-After header, got %v", elapsed)
		}
		if elapsed > 5*time.Second {
			t.Fatalf("waited too long (%v), Retry-After may not be respected", elapsed)
		}
	})

	t.Run("falls back to default wait when Retry-After missing", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		calls := 0
		_, err := RetryOnRateLimit(ctx, func() (string, error) {
			calls++
			// No Retry-After info — should use defaultRetryWait (10s)
			return "", descope.ErrRateLimitExceeded.WithMessage("rate limited")
		})
		// Context should cancel before the 10s default wait completes
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context.DeadlineExceeded (default wait > context timeout), got %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call before context timeout, got %d", calls)
		}
	})
}

func TestRetryOnRateLimitNoResult(t *testing.T) {
	t.Run("returns nil on success", func(t *testing.T) {
		calls := 0
		err := RetryOnRateLimitNoResult(context.Background(), func() error {
			calls++
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	t.Run("returns error on failure", func(t *testing.T) {
		err := RetryOnRateLimitNoResult(context.Background(), func() error {
			return errors.New("failed")
		})
		if err == nil || err.Error() != "failed" {
			t.Fatalf("expected 'failed', got %v", err)
		}
	})

	t.Run("retries on rate limit then succeeds", func(t *testing.T) {
		calls := 0
		err := RetryOnRateLimitNoResult(context.Background(), func() error {
			calls++
			if calls == 1 {
				return descope.ErrRateLimitExceeded.
					WithMessage("rate limited").
					WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, 1)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calls != 2 {
			t.Fatalf("expected 2 calls (1 retry + success), got %d", calls)
		}
	})
}
