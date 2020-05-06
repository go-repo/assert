package internal

import (
	"reflect"
	"testing"

	"github.com/lifenod/assert/diff"
)

func Equal(t *testing.T, actual, expected interface{}) bool {
	t.Helper()

	if reflect.DeepEqual(actual, expected) {
		return true
	}

	t.Log("actual (-) and expected (+) are not equal\n" + diff.Diff(actual, expected))
	return false
}
