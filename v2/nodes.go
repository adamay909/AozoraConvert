package aozoratext

import (
	"sort"
	"strings"
)

type Node struct {
	next           *Node
	prev           *Node
	firstChild     *Node
	up             *Node
	level          int
	Attr           map[string]string
	hasGaijiWithin bool
	tok            *token
	closed         bool
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

// String returns the attributes of n.
func (n *Node) String() string {

	lead := new(strings.Builder)

	if o_prettyStrings {

		for i := 0; i < n.nestingLevel(); i++ {

			lead.WriteString("\t")

		}
	}

	output.Reset()

	addToStringsBuilder(output, lead.String(), `"type": `, `"`, n.Attr["type"], `"`, ",\n")

	if n.Attr["raw"] != "" {

		addToStringsBuilder(output, lead.String(), `"raw": `, `"`, n.Attr["raw"], `"`, ",\n")

	}

	if n.Attr["raw closer"] != "" {

		addToStringsBuilder(output, lead.String(), `"raw closer": `, `"`, n.Attr["raw closer"], `"`, ",\n")

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
			addToStringsBuilder(output, lead.String(), `"`, k, `": `, `"`, n.Attr[k], `"`, ",\n")

		}
	}
	return strings.TrimSuffix(output.String(), ",\n")
}

// Parent returns the parent node of n. nil if n is top node.
func (n *Node) Parent() *Node {

	return n.firstSibling().up

}

// Siblings returns a slice of all siblings including and after n.
func (n *Node) Siblings() []*Node {

	var siblings []*Node

	if n == nil {
		return siblings
	}

	for e := n; e != nil; e = e.next {

		siblings = append(siblings, e)

	}

	return siblings
}

// Children returns the children of n in order as a slice, starting with
// the the first child of n.
func (n *Node) Children() []*Node {

	return n.firstChild.Siblings()

}

// IsLastSiblings returns whether or not n has any further siblings.
func (n *Node) IsLastSibling() bool {

	return n.next == nil

}

func (n *Node) topNode() *Node {

	e := new(Node)

	for e = n; e.Parent() != nil; e = e.Parent() {
	}

	return e

}

func (n *Node) setLevel(l int) {

	n.level = l

}

func (n *Node) setType(t string) {

	n.SetAttr("type", t)

}

func (n *Node) setRaw(s string) {

	n.SetAttr("raw", s)

}

func (n *Node) setContent(s string) {

	n.SetAttr("content", s)

}

func (n *Node) addChild(n2 *Node) {

	if n.firstChild == nil {

		n.firstChild = n2

		n2.up = n

		n2.level = n.level + 1

	} else {

		n.firstChild.addSibling(n2)

	}

	return

}

func (n *Node) addFirstChild(n2 *Node) {

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

func (n *Node) addSibling(n2 *Node) {

	if n == nil {

		return

	}

	e := n.lastSibling()

	e.next = n2

	n2.prev = e

	n2.level = n.level

	return

}

func (n *Node) lastChild() *Node {

	if n.firstChild == nil {

		return nil

	}

	return n.firstChild.lastSibling()

}

func (n *Node) firstSibling() (e *Node) {

	for e = n; e.prev != nil; e = e.prev {
	}

	return e

}

func (n *Node) lastSibling() (e *Node) {

	if n == nil {
		return n
	}

	for e = n; e.next != nil; e = e.next {
	}

	return e

}

// make node with "type" attribute set to t
func newNode(t string) *Node {

	n := new(Node)

	n.Attr = make(map[string]string)

	n.Attr["type"] = t

	return n

}

// SetAttr sets an Attr of n with the key-val pair.
func (n *Node) SetAttr(key string, val string) {

	n.Attr[key] = val

	return

}

func (n *Node) nestingLevel() int {

	l := -1

	for e := n; e != nil; e = e.Parent() {

		l++
	}

	return l
}
