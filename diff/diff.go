package diff

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/lifenod/assert/diff/internal"
)

const indent = "    "

func sprintDiffXYWithKey(
	prefix string,
	deep, ptrDeep int,
	key, typ, val string,
) string {
	if typ != "" {
		var format string
		if ptrDeep > 0 {
			format = "%s %s%s: %s(%s)(%s)\n"
		} else {
			format = "%s %s%s: %s%s(%s)\n"
		}

		return fmt.Sprintf(format,
			prefix,
			strings.Repeat(indent, deep),
			key,
			strings.Repeat("*", ptrDeep),
			typ,
			val,
		)
	} else {
		return fmt.Sprintf("%s %s%s: %s%s\n",
			prefix,
			strings.Repeat(indent, deep),
			key,
			strings.Repeat("*", ptrDeep),
			val,
		)
	}
}

func sprintDiffXY(
	prefix string,
	deep, ptrDeep int,
	typ, val string,
) string {
	if typ != "" {
		var format string
		if ptrDeep > 0 {
			format = "%s %s%s(%s)(%s)\n"
		} else {
			format = "%s %s%s%s(%s)\n"
		}

		return fmt.Sprintf(format,
			prefix,
			strings.Repeat(indent, deep),
			strings.Repeat("*", ptrDeep),
			typ,
			val,
		)
	} else {
		return fmt.Sprintf("%s %s%s%s\n",
			prefix,
			strings.Repeat(indent, deep),
			strings.Repeat("*", ptrDeep),
			val,
		)
	}
}

func diffXYStr(node *internal.Node, deep, ptrDeep int) string {
	var str string
	if node.Key != "" {
		if node.DiffXY.X != nil {
			str = str + sprintDiffXYWithKey("-", deep, ptrDeep, node.Key,
				node.DiffXY.X.Type, node.DiffXY.X.Val)
		}
		if node.DiffXY.Y != nil {
			str = str + sprintDiffXYWithKey("+", deep, ptrDeep, node.Key,
				node.DiffXY.Y.Type, node.DiffXY.Y.Val)
		}
	} else {
		if node.DiffXY.X != nil {
			str = str + sprintDiffXY("-", deep, ptrDeep,
				node.DiffXY.X.Type, node.DiffXY.X.Val)
		}
		if node.DiffXY.Y != nil {
			str = str + sprintDiffXY("+", deep, ptrDeep,
				node.DiffXY.Y.Type, node.DiffXY.Y.Val)
		}
	}
	return str
}

func levelStrWithKey(deep, ptrDeep int, key, typ string) string {
	return fmt.Sprintf("  %s%s: %s%s{\n",
		strings.Repeat(indent, deep),
		key,
		strings.Repeat("&", ptrDeep),
		typ,
	)
}

func levelStr(deep, ptrDeep int, typ string) string {
	return fmt.Sprintf("  %s%s%s{\n",
		strings.Repeat(indent, deep),
		strings.Repeat("&", ptrDeep),
		typ,
	)
}

func sprintTree(node *internal.Node, deep int, ptrDeep *int, buffer *bytes.Buffer) {
	for _, child := range node.Children {
		if child.DiffXY == nil && child.DiffNum == 0 {
			continue
		}

		if child.DiffXY != nil {
			buffer.WriteString(diffXYStr(child, deep, *ptrDeep))
			*ptrDeep = 0
			continue
		}

		switch child.Kind {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			if child.Key != "" {
				buffer.WriteString(
					levelStrWithKey(deep, *ptrDeep, child.Key, child.Type),
				)
			} else {
				buffer.WriteString(
					levelStr(deep, *ptrDeep, child.Type),
				)
			}

			*ptrDeep = 0

			sprintTree(child, deep+1, ptrDeep, buffer)

			buffer.WriteString(strings.Repeat(indent, deep) + "  }\n")
		case reflect.Ptr:
			*ptrDeep = *ptrDeep + 1
			sprintTree(child, deep, ptrDeep, buffer)
		case reflect.Interface:
			sprintTree(child, deep, ptrDeep, buffer)
		default:
			panic(fmt.Sprintf("%v kind should be handled as level", child.Kind.String()))
		}
	}
}

func calcNodeDiffNum(node *internal.Node) int {
	for _, child := range node.Children {
		if child.DiffXY != nil {
			node.DiffNum++
		} else {
			node.DiffNum = node.DiffNum + calcNodeDiffNum(child)
		}
	}
	return node.DiffNum
}

func Diff(x, y interface{}) string {
	tree := internal.Diff(x, y)

	calcNodeDiffNum(tree)

	buffer := bytes.NewBuffer(nil)
	ptrDeep := 0
	sprintTree(tree, 0, &ptrDeep, buffer)

	return buffer.String()
}
