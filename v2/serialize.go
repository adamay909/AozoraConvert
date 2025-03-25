package aozoratext

import (
	"strings"
)

// linearizes AST into a slice of *node.
// isolate will isolate the tree rooted at n. I.e., siblings
// of n are not included.

func linearizeNode(n *Node, isolate bool) []*Node {

	var l []*Node

	var prev *Node

	var flatten func(*Node)

	flatten = func(m *Node) {

		l = append(l, m)

		if prev != m.lastChild() {

			if m.firstChild != nil {

				prev = m.firstChild

				flatten(m.firstChild)

			}
		}

		if isolate {

			if m == n {
				return
			}
		}

		if m.next == nil {
			return
		}

		prev = m.next

		flatten(m.next)

		return

	}

	flatten(n)

	return l

}

func linearize(n *Node) []*Node {

	return linearizeNode(n, false)

}

func linearizeIsolate(n *Node) []*Node {

	return linearizeNode(n, true)
}

// Serialize AST given by n as a string. ingressFunc controls output when entering a node, egressFunc controls the output when exiting a node.
func Serialize(n *Node, ingressFunc, egressFunc func(*Node) string) string {

	var output = new(strings.Builder)

	var flatten func(*Node)

	flatten = func(m *Node) {

		output.WriteString(ingressFunc(m))

		for _, childNode := range m.Children() {

			flatten(childNode)

			output.WriteString(egressFunc(childNode))

		}

		return

	}

	flatten(n)

	return output.String()
}

func listSections(n *Node) string {

	egressFunc := func(n *Node) string {

		return ""

	}

	ingressFunc := func(n *Node) string {
		if n.Attr["type"] != "section" {
			return ""
		}

		output := new(strings.Builder)

		for i := 0; i < n.level; i++ {

			output.WriteString("\t")

		}

		output.WriteString(n.String())

		return output.String()

	}

	return Serialize(n, ingressFunc, egressFunc)
}
