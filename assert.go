package assert

import (
	"testing"

	"github.com/go-repo/assert/internal"
)

func Equal(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if !internal.Equal(t, actual, expected) {
		t.FailNow()
	}
}

func NotEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if !internal.NotEqual(t, actual, expected) {
		t.FailNow()
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if !internal.NoError(t, err) {
		t.FailNow()
	}
}

func Nil(t *testing.T, actual interface{}) {
	t.Helper()

	if !internal.Nil(t, actual) {
		t.FailNow()
	}
}

func NotNil(t *testing.T, actual interface{}) {
	t.Helper()

	if !internal.NotNil(t, actual) {
		t.FailNow()
	}
}
