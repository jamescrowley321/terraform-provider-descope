package infra

import (
	"context"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/stretchr/testify/require"
)

func TestRetryRejectsAmbiguousServerFailures(t *testing.T) {
	calls := 0
	expected := (&descope.Error{Code: "E999999"}).WithInfo(descope.ErrorInfoKeys.HTTPResponseStatusCode, http.StatusInternalServerError)
	_, err := RetryOnRateLimit(t.Context(), func() (string, error) { calls++; return "", expected })
	require.ErrorIs(t, err, expected)
	require.Equal(t, 1, calls)
}

func TestRateLimitRetryHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := RetryOnRateLimit(ctx, func() (string, error) { return "", descope.ErrRateLimitExceeded })
	require.ErrorIs(t, err, context.Canceled)
}

func TestRateLimitDelay(t *testing.T) {
	for _, tc := range []struct {
		seconds int
		want    time.Duration
	}{{0, 10 * time.Second}, {-1, 10 * time.Second}, {2, 2 * time.Second}, {60, time.Minute}, {math.MaxInt, time.Minute}} {
		err := descope.ErrRateLimitExceeded.WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, tc.seconds)
		require.Equal(t, tc.want, rateLimitDelay(err))
	}
}
