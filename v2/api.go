package aozoraconvert

import (
	_ "embed" //for embedding
	"errors"
	"strings"
)

// LaTeXdefinitions contains the default definitions of commands and environments
// for LaTeX output.
//
//go:embed assets/azcommands.tex
var LaTeXdefinitions string

// AozoraCSS is the default CSS to be used with HTML files rendered through aozoraconvert.
//
//go:embed assets/aozora.css
var AozoraCSS string

var tmpBuilder *strings.Builder

func init() {

	tmpBuilder = new(strings.Builder)

}

// SetJIS0208 sets output to JIS0208.
func SetJIS0208() {

	oJis0208 = true

	oJis0213 = false

	oFull = false

}

// SetJIS0213 sets output to JIS0213.
func SetJIS0213() {

	oJis0208 = false

	oJis0213 = true

	oFull = false

}

// SetFullUnicode sets output to allow the full range of unicode codepoints.
func SetFullUnicode() {

	oJis0208 = false

	oJis0213 = false

	oFull = true

}

// SetParsable sets output to a parsable form if v is true. Only useful for text
// output.
func SetParsable(v bool) {

	oParsable = v

}

// SetRubyEmph sets whether ruby are handled as text-emphasis or as ruby
// (only relevant for (X)HTML output).
func SetRubyEmph(v bool) {

	oCompatible = v

}

// SetFragment controls whether input should be treated
// as a full aozorabunko text or just a fragment.
func SetFragment(v bool) {

	if v {

		oFragment = true

	} else {

		oFragment = false
	}

}

// SetStrict tells the parser whether fixable errors in the input
// file should abort the parsing.
func SetStrict(v bool) {

	if v {

		oTolerant = false

	} else {

		oTolerant = true

	}
}

// AST returns the root node of the AST for data.
// data should be a properly formatted Aozorabunko text.
// If not, it will probably panic.
func AST(data string) (n *Node, err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))
			return

		}
	}()

	return parse(data)

}

// RenderAozoraText renders ast as a string formatted
// in the style of Aozorabunko.
func RenderAozoraText(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))

			return

		}
	}()

	renderAozoraText(ast, w)

	return
}

// RenderHTML renders ast as an html fragment.
func RenderHTML(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))
			return

		}
	}()
	if ast == nil {
		return
	}

	renderHTML(ast, w)

	return
}

//go:embed assets/skel.html
var htmltemplate string

// RenderHTMLFull renders ast as a full HTML file including doctype declararation and head element.
func RenderHTMLFull(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))

			return

		}
	}()

	title := ast.getTitle()

	if title == "" {
		title = "青空文庫テキスト"
	}

	header := strings.Replace(htmltemplate, "XXXYYY", title, 1)

	w.WriteString(header)

	RenderHTML(ast, w)

	w.WriteString(`</body>` + "\n")

	w.WriteString(`</html>` + "\n")

	return
}

// RenderNavHTML returns the table of contents for
// the text given by ast. TOC is formatted as
// an html ordered list.
func RenderNavHTML(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))
			return

		}
	}()

	renderNavHTML(ast, w)

	return
}

// RenderJSON renders ast in JSON format.
func RenderJSON(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))
			return

		}
	}()

	renderJSON(ast, w)

	return
}

// RenderLaTeX renders ast as a LaTeX fragment.
func RenderLaTeX(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))
			return

		}
	}()

	renderLaTeX(ast, w)

	return
}

// RenderLaTeXFull renders ast as a whole compileable
// LaTeX document. You will need to use the uplatex engine.
func RenderLaTeXFull(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))

			return

		}
	}()

	w.WriteString(`\documentclass[tate,paper=b6j,jafontsize=9pt]{jlreq}` + "\n\n")

	w.WriteString(`\input{azcommands}` + "\n\n")

	w.WriteString(`\setlength{\parindent}{0em}` + "\n\n")

	documentStartLaTeX(ast, w)

	RenderLaTeX(ast, w)

	documentEndLaTeX(ast, w)

	return
}

func SetVerbose() {

	oVerbose = true

}
