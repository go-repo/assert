package errorassert

import (
	"testing"

	"github.com/go-repo/assert/internal"
)

func Equal(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if !internal.Equal(t, actual, expected) {
		t.Fail()
	}
}

func NotEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if !internal.NotEqual(t, actual, expected) {
		t.Fail()
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if !internal.NoError(t, err) {
		t.Fail()
	}
}
