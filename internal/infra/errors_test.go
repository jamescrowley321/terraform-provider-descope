package infra

import (
	"errors"
	"testing"

	"github.com/descope/go-sdk/descope"
)

func TestAsValidationError(t *testing.T) {
	t.Run("extracts message from descope validation error", func(t *testing.T) {
		err := &descope.Error{
			Code:    errCodeValidationError,
			Message: "name is required",
		}
		msg, ok := AsValidationError(err)
		if !ok {
			t.Fatal("expected ok=true for validation error")
		}
		if msg != "name is required" {
			t.Fatalf("expected 'name is required', got %q", msg)
		}
	})

	t.Run("rejects descope errors with non-validation error codes", func(t *testing.T) {
		err := &descope.Error{
			Code:    "E999999",
			Message: "something else",
		}
		_, ok := AsValidationError(err)
		if ok {
			t.Fatal("expected ok=false for non-validation error")
		}
	})

	t.Run("returns false for nil error", func(t *testing.T) {
		_, ok := AsValidationError(nil)
		if ok {
			t.Fatal("expected ok=false for nil error")
		}
	})

	t.Run("returns false for non-descope error", func(t *testing.T) {
		err := errors.New("generic error")
		_, ok := AsValidationError(err)
		if ok {
			t.Fatal("expected ok=false for non-descope error")
		}
	})

	t.Run("rejects validation error when message is empty", func(t *testing.T) {
		err := &descope.Error{
			Code:    errCodeValidationError,
			Message: "",
		}
		_, ok := AsValidationError(err)
		if ok {
			t.Fatal("expected ok=false for empty message")
		}
	})
}

func TestIsNotFoundError(t *testing.T) {
	t.Run("detects 404 status code via SDK IsNotFound", func(t *testing.T) {
		// The SDK's IsNotFound() checks the HTTP status code in Info
		err := &descope.Error{
			Code: "E000000",
			Info: map[string]any{
				descope.ErrorInfoKeys.HTTPResponseStatusCode: 404,
			},
		}
		if !IsNotFoundError(err) {
			t.Fatal("expected true for not-found error")
		}
	})

	t.Run("rejects non-404 descope errors", func(t *testing.T) {
		err := &descope.Error{
			Code: "E999999",
			Info: map[string]any{
				descope.ErrorInfoKeys.HTTPResponseStatusCode: 400,
			},
		}
		if IsNotFoundError(err) {
			t.Fatal("expected false for non-404 error")
		}
	})

	t.Run("returns false for nil error", func(t *testing.T) {
		if IsNotFoundError(nil) {
			t.Fatal("expected false for nil error")
		}
	})

	t.Run("returns false for non-descope error", func(t *testing.T) {
		err := errors.New("generic error")
		if IsNotFoundError(err) {
			t.Fatal("expected false for non-descope error")
		}
	})
}
