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

	return getAozoraAST(data)

}

// RenderAozoraText renders ast as a string formatted
// in the style of Aozorabunko.
func RenderAozoraText(ast *Node) string {

	return renderAozoraText(ast)

}

// RenderHTML renders ast as an html fragment.
func RenderHTML(ast *Node) string {

	if ast == nil {
		return ""
	}

	return renderHtml(ast)

}

// RenderNavHTML returns the table of contents for
// the text given by ast. TOC is formatted as
// an html ordered list.
func RenderNavHTML(ast *Node) string {

	return renderNavHtml(ast)

}

func RenderJSON(ast *Node) string {

	return renderJson(ast)

}
