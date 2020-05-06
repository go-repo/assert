package assert_test

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/lifenod/assert"
	"github.com/lifenod/assert/errorassert"
)

const TestRunNameEnvKey = "TEST_RUN_NAME"

var (
	testRunNameMapping = map[string]func(*testing.T){}
)

var tests = []struct {
	fn func(*testing.T)

	expectedOutput      string
	expectedIsExitError bool
}{
	{
		fn:                  testEqual__IsEqual,
		expectedOutput:      "",
		expectedIsExitError: false,
	},

	{
		fn: testEqual__IsNotEqual,
		expectedOutput: `        assert_test.go:%v: actual (-) and expected (+) are not equal
            - int(1)
            + int(2)`,
		expectedIsExitError: true,
	},

	{
		fn: testEqual__IsNotEqual__SlicesWithDifferentLength,
		expectedOutput: `        assert_test.go:%v: actual (-) and expected (+) are not equal
              []int{
            -     1: int(2)
            -     2: int(3)
              }`,
		expectedIsExitError: true,
	},

	{
		fn: testEqual__IsNotEqual__NilAndEmptySlice,
		expectedOutput: `        assert_test.go:%v: actual (-) and expected (+) are not equal
            - []uint8(nil)
            + []uint8([])`,
		expectedIsExitError: true,
	},

	{
		fn: testErrorEqual__IsNotEqual,
		expectedOutput: `        assert_test.go:%v: actual (-) and expected (+) are not equal
            - string("123")
            + string("456")
            
        assert_test.go:%v: actual (-) and expected (+) are not equal
            - string("78")
            + string("90")`,
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

func testEqual__IsEqual(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, 5, 5)
	assert.Equal(t, 6, 6)
}

func testEqual__IsNotEqual(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, 1, 2)
	assert.Equal(t, 2, 5)
}

func testEqual__IsNotEqual__SlicesWithDifferentLength(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, []int{1, 2, 3}, []int{1})
}

func testEqual__IsNotEqual__NilAndEmptySlice(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Println(line + 2)
	assert.Equal(t, []byte(nil), []byte{})
}

func testErrorEqual__IsNotEqual(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("%v:%v\n", line+2, line+3)
	errorassert.Equal(t, "123", "456")
	errorassert.Equal(t, "78", "90")
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
