package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/descope/go-sdk/descope"
)

const (
	rateLimitMsg = "rate limited"
)

func newRateLimitError(retryAfter ...int) *descope.Error {
	err := descope.ErrRateLimitExceeded.WithMessage(rateLimitMsg)
	if len(retryAfter) > 0 {
		err = err.WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, retryAfter[0])
	}
	return err
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireCallCount(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected %d call(s), got %d", want, got)
	}
}

func TestRetryOnRateLimit_Success(t *testing.T) {
	calls := 0
	result, err := RetryOnRateLimit(context.Background(), func() (string, error) {
		calls++
		return "ok", nil
	})
	requireNoError(t, err)
	if result != "ok" {
		t.Fatalf("expected 'ok', got %q", result)
	}
	requireCallCount(t, calls, 1)
}

func TestRetryOnRateLimit_NonRateLimitError(t *testing.T) {
	calls := 0
	_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
		calls++
		return "", errors.New("some other error")
	})
	if err == nil || err.Error() != "some other error" {
		t.Fatalf("expected 'some other error', got %v", err)
	}
	requireCallCount(t, calls, 1)
}

func TestRetryOnRateLimit_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	_, err := RetryOnRateLimit(ctx, func() (string, error) {
		calls++
		return "", newRateLimitError()
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRetryOnRateLimit_NonRateLimitDescopeError(t *testing.T) {
	calls := 0
	descopeErr := &descope.Error{Code: "E999999", Message: "bad request"}
	_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
		calls++
		return "", descopeErr
	})
	if err == nil {
		t.Fatal("expected error")
	}
	requireCallCount(t, calls, 1)
}

func TestRetryOnRateLimit_RetriesThenSucceeds(t *testing.T) {
	calls := 0
	start := time.Now()
	result, err := RetryOnRateLimit(context.Background(), func() (string, error) {
		calls++
		if calls <= 2 {
			return "", newRateLimitError(1)
		}
		return "recovered", nil
	})
	elapsed := time.Since(start)
	requireNoError(t, err)
	if result != "recovered" {
		t.Fatalf("expected 'recovered', got %q", result)
	}
	requireCallCount(t, calls, 3)
	if elapsed < 1*time.Second {
		t.Fatalf("expected at least 1s of retry wait, got %v", elapsed)
	}
}

func TestRetryOnRateLimit_ExhaustsRetries(t *testing.T) {
	calls := 0
	_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
		calls++
		return "", newRateLimitError(1)
	})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	// maxRetries (3) loop iterations + 1 final call = 4 total
	requireCallCount(t, calls, maxRetries+1)
}

func TestRetryOnRateLimit_UsesRetryAfter(t *testing.T) {
	calls := 0
	start := time.Now()
	_, err := RetryOnRateLimit(context.Background(), func() (string, error) {
		calls++
		if calls == 1 {
			return "", newRateLimitError(2)
		}
		return "ok", nil
	})
	elapsed := time.Since(start)
	requireNoError(t, err)
	if elapsed < 2*time.Second {
		t.Fatalf("expected at least 2s retry wait from Retry-After, got %v", elapsed)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("waited too long (%v), Retry-After may not be respected", elapsed)
	}
}

func TestRetryOnRateLimit_DefaultWaitFallback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	calls := 0
	_, err := RetryOnRateLimit(ctx, func() (string, error) {
		calls++
		return "", newRateLimitError()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded (default wait > context timeout), got %v", err)
	}
	requireCallCount(t, calls, 1)
}

func TestRetryOnRateLimitNoResult_Success(t *testing.T) {
	calls := 0
	err := RetryOnRateLimitNoResult(context.Background(), func() error {
		calls++
		return nil
	})
	requireNoError(t, err)
	requireCallCount(t, calls, 1)
}

func TestRetryOnRateLimitNoResult_Error(t *testing.T) {
	err := RetryOnRateLimitNoResult(context.Background(), func() error {
		return errors.New("failed")
	})
	if err == nil || err.Error() != "failed" {
		t.Fatalf("expected 'failed', got %v", err)
	}
}

func TestRetryOnRateLimitNoResult_RetriesThenSucceeds(t *testing.T) {
	calls := 0
	err := RetryOnRateLimitNoResult(context.Background(), func() error {
		calls++
		if calls == 1 {
			return newRateLimitError(1)
		}
		return nil
	})
	requireNoError(t, err)
	requireCallCount(t, calls, 2)
}
