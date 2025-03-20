package aozoratext

import (
	"log"
	"sort"
	"strconv"
	"strings"
)

type node struct {
	next           *node
	prev           *node
	firstChild     *node
	up             *node
	level          int
	attr           map[string]string
	hasGaijiWithin bool
	lineNo         int
}

type nodeType int

const (
	dummyNode nodeType = iota
	documentNode
	sectionNode
	blockNode
	commandNode
	imageNode
	noteNode
	rubyNode
	decorationNode
	textNode
	aozoraEndnoteNode
)

func (n *node) String() string {

	lead := new(strings.Builder)

	if o_prettyStrings {

		for i := 0; i < n.level; i++ {

			lead.WriteString("\t")

		}
	}

	output.Reset()

	addToStringsBuilder(output, lead.String(), `"type": `, `"`, n.attr["type"], `"`, ",\n")

	if n.attr["raw"] != "" {

		addToStringsBuilder(output, lead.String(), `"raw": `, `"`, n.attr["raw"], `"`, ",\n")

	}

	if n.attr["raw closer"] != "" {

		addToStringsBuilder(output, lead.String(), `"raw closer": `, `"`, n.attr["raw closer"], `"`, ",\n")

	}

	var keys []string

	for k := range n.attr {

		keys = append(keys, k)

	}

	sort.Strings(keys)

	for _, k := range keys {

		switch k {

		case "type", "raw", "raw closer":
			continue

		default:
			addToStringsBuilder(output, lead.String(), `"`, k, `": `, `"`, n.attr[k], `"`, ",\n")

		}
	}
	return strings.TrimSuffix(output.String(), ",\n")
}

func (n *node) parent() *node {

	return n.firstSibling().up

}

func (n *node) topNode() *node {

	e := new(node)

	for e = n; e.parent() != nil; e = e.parent() {
	}

	return e

}

func (n *node) setLevel(l int) {

	n.level = l

}

func (n *node) setType(t string) {

	n.setAttr("type", t)

}

func (n *node) setRaw(s string) {

	n.setAttr("raw", s)

}

func (n *node) setContent(s string) {

	n.setAttr("content", s)

}

func (n *node) addChild(n2 *node) {

	if n.firstChild == nil {

		n.firstChild = n2

		n2.up = n

		n2.level = n.level + 1

	} else {

		n.firstChild.addSibling(n2)

	}

	return

}

func (n *node) addFirstChild(n2 *node) {

	if n.firstChild == nil {

		n.addChild(n2)

		return

	}

	oldfc := n.firstChild

	n2.level = oldfc.level

	n2.next = oldfc

	n2.up = n

	n.firstChild = n2

}

func (n *node) addSibling(n2 *node) {

	if n == nil {

		log.Fatal("FATAL: attempting to add sibling to NIL " + n2.attr["raw"] + " line: " + strconv.Itoa(n2.lineNo))

	}

	e := n.lastSibling()

	e.next = n2

	n2.prev = e

	n2.level = n.level

	return

}

func (n *node) lastChild() *node {

	if n.firstChild == nil {

		return nil

	}

	return n.firstChild.lastSibling()

}

func (n *node) firstSibling() (e *node) {

	for e = n; e.prev != nil; e = e.prev {
	}

	return e

}

func (n *node) lastSibling() (e *node) {

	if n == nil {
		return n
	}

	for e = n; e.next != nil; e = e.next {
	}

	return e

}

// make node with "type" attribute set to t
func newNode(t string) *node {

	n := new(node)

	n.attr = make(map[string]string)

	n.attr["type"] = t

	return n

}

func (n *node) setAttr(key string, val string) {

	n.attr[key] = val

	return

}

func (n *node) nestingLevel() (l int) {

	l = 0

	e := n

	for {

		if e.parent() == nil {

			return
		}

		e = e.parent()

		l++
	}
}

/*
func sectionLevelN(n *node) secLevel {

		switch n.attr["sectionLevel"] {

		case sectionLevelMarker[0]:
			return topsection

		case sectionLevelMarker[1]:
			return subsection

		default:
			return subsubsection

		}
	}
*/
func (n *node) firstTextOffspring() *node {

	e := new(node)

	for e = n.firstChild; e.attr["type"] != "text"; e = e.firstChild {

		if e.firstChild == nil {
			e.setAttr("type", "text")
			e.setAttr("content", "NO TEXT CHILD!")
			break
		}
	}

	return e
}
