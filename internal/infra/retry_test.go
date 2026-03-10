package infra

import (
	"context"
	"errors"
	"testing"

	"github.com/descope/go-sdk/descope"
)

func TestRetryOnRateLimit(t *testing.T) {
	t.Run("returns result on first success", func(t *testing.T) {
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

	t.Run("returns immediately for non-rate-limit error", func(t *testing.T) {
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

	t.Run("returns context error when cancelled", func(t *testing.T) {
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

	t.Run("returns non-rate-limit descope error immediately", func(t *testing.T) {
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
}
