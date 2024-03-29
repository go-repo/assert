package assert_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/go-repo/assert"
	"github.com/go-repo/assert/errorassert"
)

const TestRunNameEnvKey = "TEST_RUN_NAME"

var (
	testRunNameMapping = map[string]func(*testing.T){}
)

type testStruct struct {
	Field1 string
}

var tests = []struct {
	fn func(*testing.T)

	expectedOutput      string
	expectedIsExitError bool
}{
	{
		fn:                  testEqual_Expected,
		expectedOutput:      "",
		expectedIsExitError: false,
	},

	{
		fn: testEqual_Unexpected,
		expectedOutput: `        assert_test.go:%v: Actual (-) and expected (+) are not equal:
            - int(1)
            + int(2)`,
		expectedIsExitError: true,
	},

	{
		fn: testEqual_Unexpected_SlicesWithDifferentLength,
		expectedOutput: `        assert_test.go:%v: Actual (-) and expected (+) are not equal:
              []int{
            -     1: int(2)
            -     2: int(3)
              }`,
		expectedIsExitError: true,
	},

	{
		fn: testEqual_Unexpected_NilAndEmptySlice,
		expectedOutput: `        assert_test.go:%v: Actual (-) and expected (+) are not equal:
            - []uint8(nil)
            + []uint8([])`,
		expectedIsExitError: true,
	},

	{
		fn: testErrorAssert_Equal_Unexpected,
		expectedOutput: `        assert_test.go:%v: Actual (-) and expected (+) are not equal:
            - string("123")
            + string("456")
            
        assert_test.go:%v: Actual (-) and expected (+) are not equal:
            - string("78")
            + string("90")`,
		expectedIsExitError: true,
	},

	{
		fn:                  testNotEqual_Expected,
		expectedOutput:      "",
		expectedIsExitError: false,
	},

	{
		fn:                  testNotEqual_Unexpected,
		expectedOutput:      `        assert_test.go:%v: Actual and expected are equal: []string{"1"}`,
		expectedIsExitError: true,
	},

	{
		fn: testErrorAssert_NotEqual_Unexpected,
		expectedOutput: `        assert_test.go:%v: Actual and expected are equal: []struct {}{}
        assert_test.go:%v: Actual and expected are equal: &errors.errorString{s:"err"}`,
		expectedIsExitError: true,
	},

	{
		fn:                  testNoError_Unexpected,
		expectedOutput:      "assert_test.go:%v: Got unexpected error: error message",
		expectedIsExitError: true,
	},

	{
		fn:                  testNoError_Expected,
		expectedOutput:      "",
		expectedIsExitError: false,
	},

	{
		fn: testErrorAssert_NoError_Unexpected,
		expectedOutput: `        assert_test.go:%v: Got unexpected error: error message 1
        assert_test.go:%v: Got unexpected error: error message 2`,
		expectedIsExitError: true,
	},

	{
		fn:                  testNil_Expected,
		expectedOutput:      "",
		expectedIsExitError: false,
	},

	{
		fn:                  testNil_Unexpected,
		expectedOutput:      `        assert_test.go:%v: Expected nil but got: &assert_test.testStruct{Field1:""}`,
		expectedIsExitError: true,
	},

	{
		fn: testErrorAssert_Nil_Unexpected,
		expectedOutput: `        assert_test.go:%v: Expected nil but got: 1
        assert_test.go:%v: Expected nil but got: []int{}`,
		expectedIsExitError: true,
	},

	{
		fn:                  testNotNil_Expected,
		expectedOutput:      "",
		expectedIsExitError: false,
	},

	{
		fn:                  testNotNil_Unexpected,
		expectedOutput:      `        assert_test.go:%v: Expected not nil but got nil: <nil>`,
		expectedIsExitError: true,
	},

	{
		fn: testErrorAssert_NotNil_Unexpected,
		expectedOutput: `        assert_test.go:%v: Expected not nil but got nil: <nil>
        assert_test.go:%v: Expected not nil but got nil: (*assert_test.testStruct)(nil)`,
		expectedIsExitError: true,
	},
}

func init() {
	for _, test := range tests {
		testRunNameMapping[funcName(test.fn)] = test.fn
	}
}

func funcName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func testEqual_Expected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, 5, 5)
	assert.Equal(t, 6, 6)
}

func testEqual_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, 1, 2)
	assert.Equal(t, 2, 5)
}

func testEqual_Unexpected_SlicesWithDifferentLength(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, []int{1, 2, 3}, []int{1})
}

func testEqual_Unexpected_NilAndEmptySlice(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, []byte(nil), []byte{})
}

func testErrorAssert_Equal_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("%v:%v\n", line+2, line+3)
	errorassert.Equal(t, "123", "456")
	errorassert.Equal(t, "78", "90")
}

func testNotEqual_Expected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.NotEqual(t, 0, 1)
}

func testNotEqual_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.NotEqual(t, []string{"1"}, []string{"1"})
}

func testErrorAssert_NotEqual_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("%v:%v\n", line+2, line+3)
	errorassert.NotEqual(t, []struct{}{}, []struct{}{})
	errorassert.NotEqual(t, errors.New("err"), errors.New("err"))
}

func testNoError_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.NoError(t, errors.New("error message"))
}

func testNoError_Expected(t *testing.T) {
	var err error
	assert.NoError(t, err)
}

func testErrorAssert_NoError_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("%v:%v\n", line+2, line+3)
	errorassert.NoError(t, errors.New("error message 1"))
	errorassert.NoError(t, errors.New("error message 2"))
}

func testNil_Expected(t *testing.T) {
	assert.Nil(t, nil)
	assert.Nil(t, (*struct{})(nil))

	var ts *testStruct
	assert.Nil(t, ts)
}

func testNil_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Nil(t, &testStruct{})
}

func testErrorAssert_Nil_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("%v:%v\n", line+2, line+3)
	errorassert.Nil(t, 1)
	errorassert.Nil(t, []int{})
}

func testNotNil_Expected(t *testing.T) {
	assert.NotNil(t, 100)
	assert.NotNil(t, &testStruct{})
}

func testNotNil_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 3)
	var err error
	assert.NotNil(t, err)
}

func testErrorAssert_NotNil_Unexpected(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("%v:%v\n", line+2, line+4)
	errorassert.NotNil(t, nil)
	var ts *testStruct
	errorassert.NotNil(t, ts)
}

// Used to run a test via shell command.
func TestRun(t *testing.T) {
	testRunName := os.Getenv(TestRunNameEnvKey)
	if testRunName == "" {
		return
	}

	testFn, ok := testRunNameMapping[testRunName]
	if !ok {
		t.Fatalf("Can't find \"%v\" test function, please update testRunNameMapping variable.", testRunName)
	}

	t.Run(testRunName, testFn)
}

func execTest(fn func(t *testing.T)) (_ string, isExitError bool, _ error) {
	cmd := exec.Command("go", "test", "-count=1", "-run", "^TestRun$")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("%s=%s", TestRunNameEnvKey, funcName(fn)),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return string(output), true, nil
		} else {
			return "", false, err
		}
	}

	return string(output), false, nil
}

func TestAll(t *testing.T) {
	for _, test := range tests {
		fnName := funcName(test.fn)
		t.Logf("Run test %v\n", fnName)

		output, isExitError, err := execTest(test.fn)
		if err != nil {
			t.Fatal(fnName + " > " + err.Error())
		}

		if isExitError != test.expectedIsExitError {
			t.Fatal(fnName + " > isExitError is unexpected")
		}

		if test.expectedOutput == "" {
			continue
		}

		idx := strings.IndexByte(output, '\n')
		args := strings.Split(output[0:idx], ":")
		var argsI []interface{}
		for _, a := range args {
			argsI = append(argsI, a)
		}
		expected := fmt.Sprintf(test.expectedOutput, argsI...)
		if !strings.Contains(output, expected) {
			t.Fatalf("%v > expectedOutput is unexpected, actual output is \n```\n%v\n```\n", funcName(test.fn), output)
		}
	}
}
