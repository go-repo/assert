package internal

import "reflect"

type Node struct {
	Key  string
	Kind reflect.Kind
	Type string
	// Diff number of all children.
	DiffNum int
	DiffXY  *DiffXY

	Children []*Node
}

type XY struct {
	Kind reflect.Kind
	Type string
	Val  string
}

type DiffXY struct {
	X *XY
	Y *XY
}
