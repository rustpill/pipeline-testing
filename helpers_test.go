package pipeline

import "testing"

func assertEqual[T comparable](t *testing.T, got, want T, label string) {
	t.Helper()

	if got != want {
		t.Errorf("%s got %v, want %v", label, got, want)
	}
}
