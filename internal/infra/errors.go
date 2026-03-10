package infra

import "github.com/descope/go-sdk/descope"

const (
	errCodeValidationError = "E113007"
)

func AsValidationError(err error) (failure string, ok bool) {
	if de := descope.AsError(err, errCodeValidationError); de != nil && de.Message != "" {
		return de.Message, true
	}
	return
}

// IsNotFoundError returns true if the error is a Descope not-found error.
func IsNotFoundError(err error) bool {
	if de := descope.AsError(err); de != nil {
		return de.IsNotFound()
	}
	return false
}
