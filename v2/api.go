package aozoratext

import (
	"strings"
)

var tmpBuilder *strings.Builder

func init() {

	tmpBuilder = new(strings.Builder)

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

// SetFragment controls whether input should be treated
// as a full aozorabunko text or just a fragment.
func SetFragment(v bool) {

	if v {

		o_fragment = true

	} else {

		o_fragment = false
	}

}

// AST returns the root node of the AST for data.
// data should be a properly formatted Aozorabunko text.
// If not, it will probably panic.
func AST(data string) (*Node, error) {

	return parse(data)

}

// RenderAozoraText renders ast as a string formatted
// in the style of Aozorabunko.
func RenderAozoraText(ast *Node, w *strings.Builder) {

	renderAozoraText(ast, w)

}

// RenderHTML renders ast as an html fragment.
func RenderHTML(ast *Node, w *strings.Builder) {

	if ast == nil {
		return
	}

	renderHtml(ast, w)

}

// RenderNavHTML returns the table of contents for
// the text given by ast. TOC is formatted as
// an html ordered list.
func RenderNavHTML(ast *Node, w *strings.Builder) {

	renderNavHtml(ast, w)

}

func RenderJSON(ast *Node, w *strings.Builder) {

	renderJson(ast, w)

}
