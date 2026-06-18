package aozoraconvert

import (
	"sort"
	"strings"
)

var oPrettystrings bool

func renderJSON(n *Node, w *strings.Builder) {

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

		for k, line := range strings.Split(n.jsonString(), "\n") {

			addToStringsBuilder(w, lead, line)

			if k < len(strings.Split(n.jsonString(), "\n"))-1 {
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

func (n *Node) jsonString() string {
	output.Reset()
	addToStringsBuilder(output, `"type": `, `"`, n.Attr["type"], `"`, ",\n")
	if n.Attr["raw"] != "" {
		addToStringsBuilder(output, `"raw": `, `"`, escapeJSONstring(n.Attr["raw"]), `"`, ",\n")
	}
	if n.Attr["raw closer"] != "" {
		addToStringsBuilder(output, `"raw closer": `, `"`, escapeJSONstring(n.Attr["raw closer"]), `"`, ",\n")
	}
	var keys []string
	for k := range n.Attr {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch k {
		case "type", "raw", "raw closer":
			continue
		default:
			addToStringsBuilder(output, `"`, k, `": `, `"`, escapeJSONstring(n.Attr[k]), `"`, ",\n")
		}
	}
	return strings.TrimSuffix(output.String(), ",\n")
}

var needsJSONescape [256]bool

func init() {
	for i := 0; i < 0x20; i++ {
		needsJSONescape[i] = true
	}
	needsJSONescape['"'] = true
	needsJSONescape['\\'] = true
}

const hexDigits = "0123456789abcdef"

func escapeJSONstring(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)

	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !needsJSONescape[c] {
			continue
		}
		if start < i {
			b.WriteString(s[start:i]) //minimize writes by grouping safe chars
		}
		switch c {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		default:
			b.WriteString(`\u00`)
			b.WriteByte(hexDigits[c>>4])
			b.WriteByte(hexDigits[c&0xF])
		}
		start = i + 1
	}
	if start < len(s) {
		b.WriteString(s[start:])
	}
	return b.String()
}
