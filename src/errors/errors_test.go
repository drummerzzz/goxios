package goxios_errors

import (
	"testing"
)

func TestErrors(t *testing.T) {
	errs := []error{
		ErrInvalidBaseURL,
		ErrEmptyURL,
		ErrRelativeURL,
		ErrNilClient,
		ErrInvalidTimeout,
		ErrEmptyHeaderKey,
		ErrUnsupportedAuth,
	}

	for _, err := range errs {
		if err == nil {
			t.Error("error should not be nil")
		}
		if err.Error() == "" {
			t.Error("error message should not be empty")
		}
	}
}

