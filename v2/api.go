package aozoraConvert

import (
	_ "embed" //for embedding
	"errors"
	"strings"
)

//go:embed assets/azcommands.tex
var LaTeXdefinitions string

//go:embed assets/aozora.css
var AozoraCSS string

var tmpBuilder *strings.Builder

func init() {

	tmpBuilder = new(strings.Builder)

}

// SetSJIS0208 sets output to JIS0208.
func SetJIS0208() {

	o_jis0208 = true

	o_jis0213 = false

	o_full = false

}

func SetJIS0213() {

	o_jis0208 = false

	o_jis0213 = true

	o_full = false

}

func SetFullUnicode() {

	o_jis0208 = false

	o_jis0213 = false

	o_full = true

}

func SetRubyEmph(v bool) {

	o_compatible = v

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

func SetStrict(v bool) {

	if v {

		o_strict = true

	} else {

		o_strict = false

	}
}

// AST returns the root node of the AST for data.
// data should be a properly formatted Aozorabunko text.
// If not, it will probably panic.
func AST(data string) (n *Node, err error) {

	defer func() {

		if r := recover(); r != nil {

			//			log.Println(r)

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

			//			log.Println(r)

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

			//			log.Println(r)

			err = errors.New(r.(string))
			return

		}
	}()
	if ast == nil {
		return
	}

	renderHtml(ast, w)

	return
}

//go:embed assets/skel.html
var htmltemplate string

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

	renderNavHtml(ast, w)

	return
}

func RenderJSON(ast *Node, w *strings.Builder) (err error) {

	defer func() {

		if r := recover(); r != nil {

			err = errors.New(r.(string))
			return

		}
	}()

	renderJson(ast, w)

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
