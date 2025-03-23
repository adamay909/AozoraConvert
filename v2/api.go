package aozoratext

import (
	"strings"
)

// A Node holds information about an individual
// node in an abstract syntax tree.
type Node struct {
	nn *node
}

var tmpBuilder *strings.Builder

func init() {

	tmpBuilder = new(strings.Builder)

}

// String returns basic information of the Node n.
func (n *Node) String() string {

	tmpBuilder.Reset()

	tmpBuilder.WriteString((n.nn).String())

	return tmpBuilder.String()

}

// Attributes returns the attribute of n as
// a slice of key-value pairs.
func (n *Node) Attributes() [][2]string {

	var r [][2]string

	for k, v := range n.nn.attr {

		r = append(r, [2]string{k, v})

	}

	return r
}

// Children returns all the child nodes of n as an ordered
// slice. The slice is empty if n has no children.
func (n *Node) Children() []*Node {

	var c []*Node

	if n.nn.firstChild == nil {
		return c
	}

	for e := n.nn.firstChild; e == nil; e = e.next {

		c = append(c, &Node{nn: e})

	}

	return c

}

// Parent reutns the parent node of n. If n is
// root node, return value is nil.
func (n *Node) Parent() *Node {

	r := new(Node)

	r.nn = n.nn.parent()

	return r

}

// SetSJIS0208 controls  whether output is
// compatible with JIS0208 encoding. Default is
// false.
func SetJIS0208(v bool) {

	if v {
		o_jis0208 = true
	} else {
		o_full = true
		o_jis0208 = false
	}
	return
}

// AST returns the root node of the AST for data.
// data should be a properly formatted Aozorabunko text.
// If not, it will probably panic.
func AST(data string) *Node {
	/*
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("The document has errors.\n", r)
				fmt.Println("Exiting")
				return
			}
		}()
	*/
	rn := new(Node)

	rn.nn = getAozoraAST(data)

	return rn

}

// RenderAozoraText renders ast as a string formatted
// in the style of Aozorabunko.
func RenderAozoraText(ast *Node) string {

	return renderAozoraText(ast.nn)

}

// RenderHTML renders ast as an html fragment.
func RenderHTML(ast *Node) string {

	if ast == nil {
		return ""
	}

	return renderHtml(ast.nn)

}

// RenderNavHTML returns the table of contents for
// the text given by ast. TOC is formatted as
// an html ordered list.
func RenderNavHTML(ast *Node) string {

	return renderNavHtml(ast.nn)

}

func RenderJSON(ast *Node) string {

	return renderJson(ast.nn)

}
