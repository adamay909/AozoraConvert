package aozoratext

import (
	"strings"
)

var o_prettyStrings bool

func renderJson(n *Node, w *strings.Builder) {

	egressFunc := func(n *Node, w *strings.Builder) {

		output := new(strings.Builder)

		for i := 0; i < n.NestingLevel(); i++ {

			output.WriteString("\t")

		}

		lead := output.String()

		w.WriteString(lead)
		w.WriteString("}")

		switch {
		case n.next != nil:

			w.WriteString(",\n")

		case n.Parent() != nil:

			addToStringsBuilder(w, "\n", lead, `]`, "\n")

		default:

			w.WriteString("\n")

		}

		return

	}

	ingressFunc := func(n *Node, w *strings.Builder) {

		output := new(strings.Builder)

		for i := 0; i < n.NestingLevel(); i++ {

			output.WriteString("\t")

		}

		lead := output.String()

		w.WriteString(lead)

		w.WriteString("{\n")

		for k, line := range strings.Split(n.String(), "\n") {

			addToStringsBuilder(w, lead, line)

			if k < len(strings.Split(n.String(), "\n"))-1 {
				w.WriteString("\n")
			}
		}

		if n.HasChild() {

			addToStringsBuilder(w, ",\n", lead, `"children": [`)

		}

		w.WriteString("\n")
	}

	Serialize(n, w, ingressFunc, egressFunc)
}
