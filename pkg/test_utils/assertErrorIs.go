package test_utils

import (
	"errors"
	"testing"
)

func AssertErrorIs(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("NewInput returned wrong error, got %q want %q", got, want)
	}
}
