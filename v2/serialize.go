package aozoraConvert

import (
	"strings"
)

// linearizes AST branch rooted at n into a slice of *node.

func linearizeNode(n *Node) []*Node {

	var l []*Node

	var flatten func(*Node)

	flatten = func(m *Node) {

		l = append(l, m)

		for _, childNode := range m.Children() {

			flatten(childNode)

		}

		return

	}

	flatten(n)

	return l

}

func linearizeDescendants(n *Node) []*Node {

	return linearizeNode(n)[1:]

}

// Serialize AST given by n as a string. ingressFunc controls output when entering a node, egressFunc controls the output when exiting a node.
func Serialize(n *Node, w *strings.Builder, ingressFunc, egressFunc func(*Node, *strings.Builder)) {

	var linearize func(*Node)

	linearize = func(m *Node) {

		ingressFunc(m, w)

		for _, childNode := range m.Children() {

			linearize(childNode)

			egressFunc(childNode, w)
		}

		return

	}

	linearize(n)

	egressFunc(n, w)

	return
}

// SerializeDescendants leaves out the top node n in serializing.
func SerializeDescendants(n *Node, w *strings.Builder, ingressFunc, egressFunc func(*Node, *strings.Builder)) {

	var linearize func(*Node)

	linearize = func(m *Node) {

		if m != n {
			ingressFunc(m, w)
		}

		for _, childNode := range m.Children() {

			linearize(childNode)

			egressFunc(childNode, w)
		}

		return

	}

	linearize(n)

	//	egressFunc(n, w)

	return
}

/*

var linearize func(*Node)

	linearize = func(m *Node) {

		for _, childNode := range m.Children() {

			ingressFunc(childNode, w)

			linearize(childNode)

			egressFunc(childNode, w)

		}

		return

	}

	linearize(n)

	return
}
*/
