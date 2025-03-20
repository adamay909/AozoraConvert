package aozoratext

import (
	"strings"
)

// linearizes AST into a slice of *node.
// isolate will isolate the tree rooted at n. I.e., siblings
// of n are not included.

func linearizeNode(n *node, isolate bool) []*node {

	var l []*node

	var prev *node

	var flatten func(*node)

	flatten = func(m *node) {

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

func linearize(n *node) []*node {

	return linearizeNode(n, false)

}

func linearizeIsolate(n *node) []*node {

	return linearizeNode(n, true)
}

// linearize AST as a string. ingressFunc controls output when entering a node, egressFunc controls the output when exiting a node.
func stringify(n *node, ingressFunc, egressFunc func(*node) string) string {

	var output = new(strings.Builder)

	const (
		oEnterNode = false
		oLeaveNode = true
	)

	prev := new(node)

	var flatten func(*node)

	flatten = func(m *node) {

		output.WriteString(ingressFunc(m))

		if prev != m.lastChild() {

			if m.firstChild != nil {

				prev = m.firstChild

				flatten(m.firstChild)

			}

		}

		output.WriteString(egressFunc(m))

		//		if m == n {
		//			return
		//		}

		if m.next == nil {
			return
		}

		prev = m.next

		flatten(m.next)

		return

	}

	prev = n

	flatten(n)

	return output.String()
}

func listTypes(n *node) string {

	return renderJson(n)

	ingressFunc := func(n *node) string {

		output := new(strings.Builder)

		o_prettyStrings = true

		for i := 0; i < n.level; i++ {

			output.WriteString("\t")

		}
		output.WriteString("{\n" + n.String())

		return output.String() + "\n"

	}

	egressFunc := func(n *node) string {

		output := new(strings.Builder)

		o_prettyStrings = true

		for i := 0; i < n.level; i++ {

			output.WriteString("\t")

		}

		return output.String() + "}\n"

	}

	return stringify(n, ingressFunc, egressFunc)
}

// list tokens the parser did not handle but recognized as unimplemented.
func listUnimplemented(n *node) string {

	ingressFunc := func(n *node) string {

		output := new(strings.Builder)

		if n.attr["type"] == "unimplemented" {

			addToStringsBuilder(output, n.attr["content"], "\n")

		}

		return output.String()

	}

	egressFunc := func(n *node) string {

		return ""

	}

	return "Unimplemented:\n\n" + stringify(n, ingressFunc, egressFunc)
}

func listSections(n *node) string {

	egressFunc := func(n *node) string {

		return ""

	}

	ingressFunc := func(n *node) string {
		if n.attr["type"] != "section" {
			return ""
		}

		output := new(strings.Builder)

		for i := 0; i < n.level; i++ {

			output.WriteString("\t")

		}

		output.WriteString(n.String())

		return output.String()

	}

	return stringify(n, ingressFunc, egressFunc)
}
