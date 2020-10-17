## Usage

```go
package assert_test

import (
	"testing"

	"github.com/go-repo/assert"
)

type StructA struct {
	A int64
	B string
}

func TestEqual(t *testing.T) {
	assert.Equal(t,
		&StructA{
			A: 1,
			B: "str",
		},
		&StructA{
			A: 2,
			B: "str",
		},
	)
}
```

Run the test and output:

```
=== RUN   TestEqual
    TestEqual: main_test.go:15: actual (-) and expected (+) are not equal
          &assert_test.StructA{
        -     A: int64(1)
        +     A: int64(2)
          }

--- FAIL: TestEqual (0.00s)
FAIL
```

## errorassert

Useful for table test, you can test all cases even if one of them is failed, for example:

```go
package errorassert_test

import (
	"testing"

	"github.com/go-repo/assert/errorassert"
)

type StructA struct {
	A int64
	B string
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{
			name:     "equal",
			actual:   "123str",
			expected: "123str",
		},

		{
			name:     "not equal for int type",
			actual:   1,
			expected: 2,
		},

		{
			name:     "not equal for struct type",
			actual:   StructA{B: "abc"},
			expected: StructA{B: "def"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errorassert.Equal(t, test.actual, test.expected)
		})
	}
}
```

Run the test and output:

```
=== RUN   TestEqual
=== RUN   TestEqual/equal
=== RUN   TestEqual/not_equal_for_int_type
    TestEqual/not_equal_for_int_type: main_test.go:40: actual (-) and expected (+) are not equal
        - int(1)
        + int(2)

=== RUN   TestEqual/not_equal_for_struct_type
    TestEqual/not_equal_for_struct_type: main_test.go:40: actual (-) and expected (+) are not equal
          errorassert_test.StructA{
        -     B: string("abc")
        +     B: string("def")
          }

--- FAIL: TestEqual (0.00s)
    --- PASS: TestEqual/equal (0.00s)
    --- FAIL: TestEqual/not_equal_for_int_type (0.00s)
    --- FAIL: TestEqual/not_equal_for_struct_type (0.00s)
FAIL
```
