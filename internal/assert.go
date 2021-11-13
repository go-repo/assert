package internal

import (
	"reflect"
	"testing"

	"github.com/go-repo/assert/diff"
)

func Equal(t *testing.T, actual, expected interface{}) bool {
	t.Helper()

	if reflect.DeepEqual(actual, expected) {
		return true
	}

	t.Log("Actual (-) and expected (+) are not equal:\n" + diff.Diff(actual, expected))
	return false
}

func NotEqual(t *testing.T, actual, expected interface{}) bool {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		return true
	}

	t.Logf("Actual and expected are equal: %#v\n", actual)
	return false
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)
	kind := val.Kind()
	if kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Map ||
		kind == reflect.Ptr ||
		kind == reflect.UnsafePointer ||
		kind == reflect.Interface ||
		kind == reflect.Slice {
		return val.IsNil()
	}

	return false
}

func NoError(t *testing.T, err error) bool {
	t.Helper()

	if isNil(err) {
		return true
	}

	t.Logf("Got unexpected error: %v\n", err)
	return false
}

func Nil(t *testing.T, actual interface{}) bool {
	t.Helper()

	if isNil(actual) {
		return true
	}

	t.Logf("Expected nil but got: %#v\n", actual)
	return false
}

func NotNil(t *testing.T, actual interface{}) bool {
	t.Helper()

	if !isNil(actual) {
		return true
	}

	t.Logf("Expected not nil but got nil: %#v\n", actual)
	return false
}
