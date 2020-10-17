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
