package aozoratext

import (
	"strings"
)

var o_prettyStrings bool

func renderJson(n *node) string {

	egressFunc := func(n *node) string {

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

			if n.parent() != nil {

				output.WriteString("\n" + lead + `]` + "\n")

			} else {

				output.WriteString("\n")

			}
		}

		return output.String()

	}

	ingressFunc := func(n *node) string {

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

	return stringify(n, ingressFunc, egressFunc)
}
