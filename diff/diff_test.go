package diff

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/lifenod/assert/diff/internal"
	"github.com/lifenod/assert/diff/test/draw_tree"
)

type S1 struct {
	bool          bool
	int           int
	int8          int8
	int16         int16
	int32         int32
	int64         int64
	uint          uint
	uint8         uint8
	uint16        uint16
	uint32        uint32
	uint64        uint64
	uintptr       uintptr
	float32       float32
	float64       float64
	complex64     complex64
	complex128    complex128
	array         [2]int
	ch            chan string
	fn            func()
	inter         interface{}
	mapping       map[string]string
	ptr           *uint
	slice         []float32
	string        string
	stru          S2
	unsafePointer unsafe.Pointer
}

type S2 struct {
	bool bool
}

type S3 struct {
	pppBool       ***bool
	pInt          *int
	pUint64       *uint64
	pMapping      *map[S4]*S4
	pArray        *[2]float32
	pSlice        *[]string
	pStru         *S4
	ch            chan int
	fn            func()
	mapping       map[string]bool
	unsafePointer unsafe.Pointer
	inter         interface{}
	slice         []S4
}

type S4 struct {
	int int
}

type S5 struct {
	str  string
	self *S5
}

func TestDiff__AllKindsAreNotEqual(t *testing.T) {
	var uintX uint = 23
	var uintY uint = 123
	var chanX = make(chan string)
	var chanY = make(chan string)
	var fnX = func() {}
	var fnY = func() {}
	x := S1{
		bool:       true,
		int:        1,
		int8:       2,
		int16:      3,
		int32:      4,
		int64:      5,
		uint:       6,
		uint8:      7,
		uint16:     8,
		uint32:     9,
		uint64:     10,
		uintptr:    11,
		float32:    12,
		float64:    13,
		complex64:  14,
		complex128: 15,
		array:      [2]int{16, 17},
		ch:         chanX,
		fn:         fnX,
		inter:      18,
		mapping: map[string]string{
			"19": "20",
			"21": "22",
		},
		ptr:    &uintX,
		slice:  []float32{23, 24},
		string: "25",
		stru: S2{
			bool: true,
		},
		unsafePointer: unsafe.Pointer(&uintX),
	}

	y := S1{
		bool:       false,
		int:        101,
		int8:       102,
		int16:      103,
		int32:      104,
		int64:      105,
		uint:       106,
		uint8:      107,
		uint16:     108,
		uint32:     109,
		uint64:     110,
		uintptr:    111,
		float32:    112,
		float64:    113,
		complex64:  114,
		complex128: 115,
		array:      [2]int{116, 117},
		ch:         chanY,
		fn:         fnY,
		inter:      118,
		mapping: map[string]string{
			"119": "120",
			"21":  "22",
		},
		ptr:    &uintY,
		slice:  []float32{123, 124},
		string: "125",
		stru: S2{
			bool: false,
		},
		unsafePointer: unsafe.Pointer(&uintY),
	}

	expectedDiff := fmt.Sprintf(`  diff.S1{
-     bool: bool(true)
+     bool: bool(false)
-     int: int(1)
+     int: int(101)
-     int8: int8(2)
+     int8: int8(102)
-     int16: int16(3)
+     int16: int16(103)
-     int32: int32(4)
+     int32: int32(104)
-     int64: int64(5)
+     int64: int64(105)
-     uint: uint(6)
+     uint: uint(106)
-     uint8: uint8(7)
+     uint8: uint8(107)
-     uint16: uint16(8)
+     uint16: uint16(108)
-     uint32: uint32(9)
+     uint32: uint32(109)
-     uint64: uint64(10)
+     uint64: uint64(110)
-     uintptr: uintptr(11)
+     uintptr: uintptr(111)
-     float32: float32(12)
+     float32: float32(112)
-     float64: float64(13)
+     float64: float64(113)
-     complex64: complex64((14+0i))
+     complex64: complex64((114+0i))
-     complex128: complex128((15+0i))
+     complex128: complex128((115+0i))
      array: [2]int{
-         0: int(16)
+         0: int(116)
-         1: int(17)
+         1: int(117)
      }
-     ch: chan string(%v)
+     ch: chan string(%v)
-     fn: func()(%v)
+     fn: func()(%v)
-     inter: int(18)
+     inter: int(118)
      mapping: map[string]string{
-         "19": string("20")
+         "119": string("120")
      }
-     ptr: *(uint)(23)
+     ptr: *(uint)(123)
      slice: []float32{
-         0: float32(23)
+         0: float32(123)
-         1: float32(24)
+         1: float32(124)
      }
-     string: string("25")
+     string: string("125")
      stru: diff.S2{
-         bool: bool(true)
+         bool: bool(false)
      }
-     unsafePointer: unsafe.Pointer(%v)
+     unsafePointer: unsafe.Pointer(%v)
  }
`, chanX, chanY, reflect.ValueOf(fnX), reflect.ValueOf(fnY), &uintX, &uintY)
	diff := Diff(x, y)
	if diff != expectedDiff {
		t.Fatal()
	}
}

func TestDiff__AllKindsAreEqual(t *testing.T) {
	var uintX uint = 23
	var uintY uint = 23
	var ch = make(chan string)

	x := S1{
		bool:       true,
		int:        1,
		int8:       2,
		int16:      3,
		int32:      4,
		int64:      5,
		uint:       6,
		uint8:      7,
		uint16:     8,
		uint32:     9,
		uint64:     10,
		uintptr:    11,
		float32:    12,
		float64:    13,
		complex64:  14,
		complex128: 15,
		array:      [2]int{16, 17},
		ch:         ch,
		fn:         nil,
		inter:      18,
		mapping: map[string]string{
			"19": "20",
			"21": "22",
		},
		ptr:    &uintX,
		slice:  []float32{23, 24},
		string: "25",
		stru: S2{
			bool: true,
		},
		unsafePointer: unsafe.Pointer(&uintX),
	}

	y := S1{
		bool:       true,
		int:        1,
		int8:       2,
		int16:      3,
		int32:      4,
		int64:      5,
		uint:       6,
		uint8:      7,
		uint16:     8,
		uint32:     9,
		uint64:     10,
		uintptr:    11,
		float32:    12,
		float64:    13,
		complex64:  14,
		complex128: 15,
		array:      [2]int{16, 17},
		ch:         ch,
		fn:         nil,
		inter:      18,
		mapping: map[string]string{
			"19": "20",
			"21": "22",
		},
		ptr:    &uintY,
		slice:  []float32{23, 24},
		string: "25",
		stru: S2{
			bool: true,
		},
		unsafePointer: unsafe.Pointer(&uintX),
	}

	diff := Diff(x, y)
	if diff != "" {
		t.Fatal()
	}
}

func TestDiff__NilParam(t *testing.T) {
	diff := Diff(nil, 2)
	if diff != "- <nil>\n+ int(2)\n" {
		t.Fatal()
	}

	diff = Diff("s", nil)
	if diff != "- string(\"s\")\n+ <nil>\n" {
		t.Fatal()
	}

	diff = Diff(nil, nil)
	if diff != "" {
		t.Fatal()
	}
}

func ptrInt(i int) *int {
	return &i
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

func TestDiff__NilValues(t *testing.T) {
	yBool := true
	yPtrBool := &yBool
	yPtrPtrBool := &yPtrBool

	yInt := 1
	yUint64 := uint64(2)

	yFn := func() {}
	yCh := make(chan int)

	x := &S3{}
	y := &S3{
		pppBool: &yPtrPtrBool,
		pInt:    &yInt,
		pUint64: &yUint64,
		pMapping: &map[S4]*S4{
			S4{int: 3}: nil,
		},
		pArray: &[2]float32{4},
		pSlice: &[]string{"5"},
		pStru: &S4{
			int: 6,
		},
		ch:            yCh,
		fn:            yFn,
		mapping:       map[string]bool{"7": true},
		unsafePointer: unsafe.Pointer(&yBool),
		inter:         S4{int: 8},
		slice:         []S4{{int: 9}},
	}

	expectedDiff := fmt.Sprintf(`  &diff.S3{
-     pppBool: ***bool(<nil>)
+     pppBool: ***bool(%v)
-     pInt: *int(<nil>)
+     pInt: *int(%v)
-     pUint64: *uint64(<nil>)
+     pUint64: *uint64(%v)
-     pMapping: *map[diff.S4]*diff.S4(<nil>)
+     pMapping: *map[diff.S4]*diff.S4(&map[{3}:<nil>])
-     pArray: *[2]float32(<nil>)
+     pArray: *[2]float32(&[4 0])
-     pSlice: *[]string(<nil>)
+     pSlice: *[]string(&[5])
-     pStru: *diff.S4(<nil>)
+     pStru: *diff.S4(&{6})
-     ch: chan int(<nil>)
+     ch: chan int(%v)
-     fn: func()(<nil>)
+     fn: func()(%v)
-     mapping: map[string]bool(map[])
+     mapping: map[string]bool(map[7:true])
-     unsafePointer: unsafe.Pointer(<nil>)
+     unsafePointer: unsafe.Pointer(%v)
-     inter: interface {}(interface {}(nil))
+     inter: interface {}(diff.S4{int:8})
-     slice: []diff.S4([])
+     slice: []diff.S4([{9}])
  }
`, &yPtrPtrBool, &yInt, &yUint64,
		reflect.ValueOf(yCh), reflect.ValueOf(yFn), &yBool,
	)
	diff := Diff(x, y)
	if diff != expectedDiff {
		t.Fatal()
	}
}

func TestDiff__MultilevelPointers(t *testing.T) {
	x := 1
	px := &x
	ppx := &px
	pppx := &ppx

	y := 2
	py := &y
	ppy := &py
	pppy := &ppy

	diff := Diff(pppx, pppy)
	if diff != "- ***(int)(1)\n+ ***(int)(2)\n" {
		t.Fatal()
	}

	var emptyPP *int
	**pppy = emptyPP

	diff = Diff(pppx, pppy)
	if diff != fmt.Sprintf("- **(*int)(%v)\n+ **(*int)(<nil>)\n",
		unsafe.Pointer(&x)) {
		t.Fatal()
	}
}

func TestDiff__ReferenceCycle(t *testing.T) {
	x := &S5{
		str: "1",
	}
	x.self = x

	y := &S5{
		str: "2",
	}
	y.self = y

	diff := Diff(x, y)
	expectedDiff := `  &diff.S5{
-     str: string("1")
+     str: string("2")
      self: &diff.S5{
-         str: string("1")
+         str: string("2")
      }
  }
`
	if diff != expectedDiff {
		t.Fatal()
	}
}

// The test can generate a tree graph, used to analyze tree structure.
func testDrawTree(t *testing.T) {
	yBool := true
	yPtrBool := &yBool
	yPtrPtrBool := &yPtrBool

	x := &S3{}
	y := &S3{
		pppBool: &yPtrPtrBool,
		pInt:    ptrInt(1),
		pUint64: ptrUint64(2),
		pMapping: &map[S4]*S4{
			S4{int: 3}: nil,
		},
		pArray: &[2]float32{4},
		pSlice: &[]string{"5"},
		pStru: &S4{
			int: 6,
		},
		ch:            make(chan int),
		fn:            func() {},
		mapping:       map[string]bool{"7": true},
		unsafePointer: unsafe.Pointer(&yBool),
		inter:         S4{int: 8},
		slice:         []S4{{int: 9}},
	}

	tree := internal.Diff(x, y)
	calcNodeDiffNum(tree)

	err := draw_tree.DrawTreeToFile(tree, "tree.dot")
	if err != nil {
		t.Fatal(err)
	}
}
