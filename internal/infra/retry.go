package infra

import (
	"context"
	"slices"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Transient races with the backend's async post-creation work: E111009 = management key ReBAC tuples not yet
// written, E111604 = theme version conflict with the creation handler republishing it. Both clear on the next apply.
var retryableErrorCodes = []string{"E111009", "E111604"}

var retryDelays = []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second}

// Safe to replay: the backend returns these codes before the request has any effect.
func retrying(ctx context.Context, call func() (*api.HTTPResponse, error)) (*api.HTTPResponse, error) {
	return RetryOnRateLimit(ctx, call)
}

func RetryOnRateLimit[T any](ctx context.Context, call func() (T, error)) (T, error) {
	res, err := call()
	for _, delay := range retryDelays {
		de := descope.AsError(err)
		if de == nil {
			break
		}
		if de.Code == descope.ErrRateLimitExceeded.Code {
			delay = rateLimitDelay(de)
		} else if !slices.Contains(retryableErrorCodes, de.Code) {
			break
		}
		tflog.Info(ctx, "Retrying after a transient backend error", map[string]any{"code": de.Code, "delay": delay.String()})
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case <-time.After(delay):
		}
		res, err = call()
	}
	return res, err
}

func RetryOnRateLimitNoResult(ctx context.Context, call func() error) error {
	_, err := RetryOnRateLimit(ctx, func() (struct{}, error) { return struct{}{}, call() })
	return err
}

func rateLimitDelay(de *descope.Error) time.Duration {
	if seconds, ok := de.Info[descope.ErrorInfoKeys.RateLimitExceededRetryAfter].(int); ok && seconds > 0 {
		if seconds > 60 {
			return time.Minute
		}
		return time.Duration(seconds) * time.Second
	}
	return 10 * time.Second
}
