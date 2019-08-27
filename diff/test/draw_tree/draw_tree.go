package draw_tree

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/lifenod/go-assert/diff/internal"
)

var nextNodeName = func() func() string {
	nextNodeNum := 0
	return func() string {
		num := nextNodeNum
		nextNodeNum = nextNodeNum + 1
		return fmt.Sprint(num)
	}
}()

func DrawTreeToFile(tree *internal.Node, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString("graph NodeGraph {\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("node [shape = \"box\", fontname = \"menlo\"]\n")
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	drawTree(nextNodeName(), tree, buffer)

	_, err = f.WriteString(buffer.String())
	if err != nil {
		return err
	}

	_, err = f.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}

func printNodeContent(nodeName string, node *internal.Node, buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`%v [label="\
key:      %v\l\
kind:     %v\l\
type:     %v\l\
diff_num: %v\l\
"]`+"\n",
		nodeName,
		strings.Replace(node.Key, `"`, `\"`, -1),
		node.Kind.String(),
		node.Type,
		node.DiffNum,
	))
}

func printNodeDiffXYContent(nodeName string, node *internal.Node, buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`%v [label="\
key:    %v\l\
kind:   %v\l\
type:   %v\l\
diff_num: %v\l\
diffxy:\l\
  x:\l\
    kind: %v\l\
    type: %v\l\
    val:  %v\l\
  y:\l\
    kind: %v\l\
    type: %v\l\
    val:  %v\l\
"]`+"\n",
		nodeName,
		strings.Replace(node.Key, `"`, `\"`, -1),
		node.Kind.String(),
		node.Type,
		node.DiffNum,
		node.DiffXY.X.Kind.String(),
		node.DiffXY.X.Type,
		strings.Replace(node.DiffXY.X.Val, `"`, `\"`, -1),
		node.DiffXY.Y.Kind.String(),
		node.DiffXY.Y.Type,
		strings.Replace(node.DiffXY.Y.Val, `"`, `\"`, -1),
	))
}

func printNodeDiffXContent(nodeName string, node *internal.Node, buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`%v [label="\
key:    %v\l\
kind:   %v\l\
type:   %v\l\
diff_num: %v\l\
diffxy:\l\
  x:\l\
    kind: %v\l\
    type: %v\l\
    val:  %v\l\
"]`+"\n",
		nodeName,
		strings.Replace(node.Key, `"`, `\"`, -1),
		node.Kind.String(),
		node.Type,
		node.DiffNum,
		node.DiffXY.X.Kind.String(),
		node.DiffXY.X.Type,
		strings.Replace(node.DiffXY.X.Val, `"`, `\"`, -1),
	))
}

func printNodeDiffYContent(nodeName string, node *internal.Node, buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`%v [label="\
key:    %v\l\
kind:   %v\l\
type:   %v\l\
diff_num: %v\l\
diffxy:\l\
  y:\l\
    kind: %v\l\
    type: %v\l\
    val:  %v\l\
"]`+"\n",
		nodeName,
		strings.Replace(node.Key, `"`, `\"`, -1),
		node.Kind.String(),
		node.Type,
		node.DiffNum,
		node.DiffXY.Y.Kind.String(),
		node.DiffXY.Y.Type,
		strings.Replace(node.DiffXY.Y.Val, `"`, `\"`, -1),
	))
}

func drawNodeContent(nodeName string, node *internal.Node, buffer *bytes.Buffer) {
	if node.DiffXY == nil {
		printNodeContent(nodeName, node, buffer)
		return
	}

	if node.DiffXY.X != nil && node.DiffXY.Y != nil {
		printNodeDiffXYContent(nodeName, node, buffer)
		return
	}

	if node.DiffXY.X != nil {
		printNodeDiffXContent(nodeName, node, buffer)
		return
	}

	if node.DiffXY.Y != nil {
		printNodeDiffYContent(nodeName, node, buffer)
		return
	}
}

func drawTree(parentNodeName string, node *internal.Node, buffer *bytes.Buffer) {
	drawNodeContent(parentNodeName, node, buffer)

	for _, child := range node.Children {
		childNodeName := nextNodeName()
		buffer.WriteString(fmt.Sprintf("%v -- %v\n", parentNodeName, childNodeName))

		drawTree(childNodeName, child, buffer)
	}
}
