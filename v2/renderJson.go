package aozoratext

import (
	"strings"
)

var o_prettyStrings bool

func renderJson(n *Node) string {

	egressFunc := func(n *Node) string {

		output := new(strings.Builder)

		o_prettyStrings = true

		for i := 0; i < n.level; i++ {

			output.WriteString("\t")

		}

		lead := output.String()

		output.Reset()

		output.WriteString(lead + "}")

		if n.next != nil {

			output.WriteString(",\n")

		} else {

			if n.Parent() != nil {

				output.WriteString("\n" + lead + `]` + "\n")

			} else {

				output.WriteString("\n")

			}
		}

		return output.String()

	}

	ingressFunc := func(n *Node) string {

		output := new(strings.Builder)

		o_prettyStrings = true

		for i := 0; i < n.level; i++ {

			output.WriteString("\t")

		}

		lead := output.String()

		output.Reset()

		output.WriteString(lead + "{\n" + n.String())

		if n.firstChild != nil {

			output.WriteString(",\n")

			output.WriteString(lead + `"children": [`)

		}

		return output.String() + "\n"

	}

	return Serialize(n, ingressFunc, egressFunc)
}
