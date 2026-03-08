package infra

import "github.com/descope/go-sdk/descope"

const (
	errCodeValidationError = "E113007"
)

func AsValidationError(err error) (failure string, ok bool) {
	if err, ok := err.(*descope.Error); ok && err.Code == errCodeValidationError && err.Message != "" {
		return err.Message, true
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
