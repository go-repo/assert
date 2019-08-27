package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func createNewCurrentNode(current *Node, v reflect.Value, key string) *Node {
	child := &Node{
		Key:  key,
		Kind: v.Kind(),
		Type: v.Type().String(),
	}
	current.Children = append(current.Children, child)
	return child
}

func diffXYVal(v reflect.Value) string {
	if v.Kind() == reflect.String ||
		v.Kind() == reflect.Interface {
		return fmt.Sprintf("%#v", v)
	}

	return fmt.Sprintf("%v", v)
}

func createChildNodes(current *Node, x, y reflect.Value, key string) {
	current.Children = append(current.Children, &Node{
		Key: key,
		DiffXY: &DiffXY{
			X: &XY{
				Kind: x.Kind(),
				Type: x.Type().String(),
				Val:  diffXYVal(x),
			},
			Y: &XY{
				Kind: y.Kind(),
				Type: y.Type().String(),
				Val:  diffXYVal(y),
			},
		},
	})
}

func createChildNodeForX(current *Node, x reflect.Value, key string) {
	current.Children = append(current.Children, &Node{
		Key: key,
		DiffXY: &DiffXY{
			X: &XY{
				Kind: x.Kind(),
				Type: x.Type().String(),
				Val:  diffXYVal(x),
			},
		},
	})
}

func createChildNodeForY(current *Node, y reflect.Value, key string) {
	current.Children = append(current.Children, &Node{
		Key: key,
		DiffXY: &DiffXY{
			Y: &XY{
				Kind: y.Kind(),
				Type: y.Type().String(),
				Val:  diffXYVal(y),
			},
		},
	})
}

func diffNil(current *Node, x, y reflect.Value, key string) bool {
	if x.IsNil() || y.IsNil() {
		if x.IsNil() == y.IsNil() {
			return true
		}

		createChildNodes(current, x, y, key)
		return true
	}

	return false
}

func ifFalseThenCreateChildNodes(b bool, current *Node, x, y reflect.Value, key string) {
	if !b {
		createChildNodes(current, x, y, key)
	}
}

func cmpMap(curr *Node, x, y reflect.Value, key string, visited map[visit]bool) {
	if diffNil(curr, x, y, key) {
		return
	}

	yMap := map[interface{}]bool{}
	newNode := createNewCurrentNode(curr, x, key)

	for _, k := range x.MapKeys() {
		vx := x.MapIndex(k)
		vy := y.MapIndex(k)
		keyStr := fmt.Sprintf("%#v", k)

		if !vy.IsValid() {
			createChildNodeForX(newNode, vx, keyStr)
			continue
		}

		yMap[newValue(k).Interface()] = true
		deepDiff(newNode, vx, vy, keyStr, visited)
	}

	for _, k := range y.MapKeys() {
		if yMap[newValue(k).Interface()] {
			continue
		}

		createChildNodeForY(newNode, y.MapIndex(k), fmt.Sprintf("%#v", k))
	}
}

func diffIsValid(curr *Node, x, y reflect.Value, key string) bool {
	if !x.IsValid() {
		if !y.IsValid() {
			return true
		}

		curr.Children = append(curr.Children, &Node{
			Key: key,
			DiffXY: &DiffXY{
				X: &XY{
					Val: "<nil>",
				},
				Y: &XY{
					Kind: y.Kind(),
					Type: y.Type().String(),
					Val:  diffXYVal(y),
				},
			},
		})
		return true
	}

	if !y.IsValid() {
		curr.Children = append(curr.Children, &Node{
			Key: key,
			DiffXY: &DiffXY{
				X: &XY{
					Kind: x.Kind(),
					Type: x.Type().String(),
					Val:  diffXYVal(x),
				},
				Y: &XY{
					Val: "<nil>",
				},
			},
		})
		return true
	}

	return false
}

// for access unexported value
func newValue(v reflect.Value) reflect.Value {
	return reflect.NewAt(v.Type(), valueDataPtr(v)).Elem()
}

func valueDataPtr(v reflect.Value) unsafe.Pointer {
	return *(*unsafe.Pointer)(
		unsafe.Pointer(
			reflect.ValueOf(&v).Elem().FieldByName("ptr").UnsafeAddr(),
		),
	)
}

// Copy from https://github.com/golang/go/blob/919594830f17f25c9e971934d825615463ad8a10/src/reflect/deepequal.go#L11-L19
// During deepValueEqual, must keep track of checks that are
// in progress. The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Copy from https://github.com/golang/go/blob/919594830f17f25c9e971934d825615463ad8a10/src/reflect/deepequal.go#L34-L63
func isReferenceCycle(v1, v2 reflect.Value, visited map[visit]bool) bool {
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	return false
}

func deepDiff(curr *Node, x, y reflect.Value, key string, visited map[visit]bool) {
	if diffIsValid(curr, x, y, key) {
		return
	}

	if x.Type() != y.Type() {
		createChildNodes(curr, x, y, key)
		return
	}

	if isReferenceCycle(x, y, visited) {
		return
	}

	switch x.Kind() {
	case reflect.Bool:
		ifFalseThenCreateChildNodes(x.Bool() == y.Bool(), curr, x, y, key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ifFalseThenCreateChildNodes(x.Int() == y.Int(), curr, x, y, key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ifFalseThenCreateChildNodes(x.Uint() == y.Uint(), curr, x, y, key)
	case reflect.Float32, reflect.Float64:
		ifFalseThenCreateChildNodes(x.Float() == y.Float(), curr, x, y, key)
	case reflect.Complex64, reflect.Complex128:
		ifFalseThenCreateChildNodes(x.Complex() == y.Complex(), curr, x, y, key)
	case reflect.Array:
		newNode := createNewCurrentNode(curr, x, key)
		for i := 0; i < x.Len(); i++ {
			deepDiff(newNode, x.Index(i), y.Index(i), strconv.Itoa(i), visited)
		}
	case reflect.Chan:
		ifFalseThenCreateChildNodes(
			newValue(x).Interface() == newValue(y).Interface(),
			curr, x, y, key,
		)
	case reflect.Func:
		ifFalseThenCreateChildNodes(x.IsNil() && y.IsNil(), curr, x, y, key)
	case reflect.Interface:
		if diffNil(curr, x, y, key) {
			return
		}
		newNode := createNewCurrentNode(curr, x, key)
		deepDiff(newNode, x.Elem(), y.Elem(), key, visited)
	case reflect.Map:
		cmpMap(curr, x, y, key, visited)
	case reflect.Ptr:
		if diffNil(curr, x, y, key) {
			return
		}
		newNode := createNewCurrentNode(curr, x, key)
		deepDiff(newNode, x.Elem(), y.Elem(), key, visited)
	case reflect.Slice:
		if diffNil(curr, x, y, key) {
			return
		}
		newNode := createNewCurrentNode(curr, x, key)
		for i := 0; i < x.Len(); i++ {
			deepDiff(newNode, x.Index(i), y.Index(i), strconv.Itoa(i), visited)
		}
	case reflect.String:
		ifFalseThenCreateChildNodes(x.String() == y.String(), curr, x, y, key)
	case reflect.Struct:
		newNode := createNewCurrentNode(curr, x, key)
		for i, n := 0, x.NumField(); i < n; i++ {
			deepDiff(newNode, x.Field(i), y.Field(i), x.Type().Field(i).Name, visited)
		}
	case reflect.UnsafePointer:
		ifFalseThenCreateChildNodes(x.Pointer() == y.Pointer(), curr, x, y, key)
	default:
		panic(fmt.Sprintf("%v kind is not supported", x.Kind().String()))
	}
}

func Diff(x, y interface{}) *Node {
	root := &Node{}
	deepDiff(root, reflect.ValueOf(x), reflect.ValueOf(y), "", make(map[visit]bool))
	return root
}
