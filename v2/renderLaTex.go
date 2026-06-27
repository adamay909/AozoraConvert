package aozoraconvert

import (
	"log"
	"strconv"
	"strings"
)

func renderLaTeX(n *Node, w *strings.Builder) {

	Serialize(n, w, azrLaTeXFormatterOpen, azrLaTeXFormatterClose)

}

func renderInnerLaTeX(n *Node, w *strings.Builder) {

	SerializeDescendants(n, w, azrLaTeXFormatterOpen, azrLaTeXFormatterClose)

}

func azrLaTeXFormatterOpen(n *Node, w *strings.Builder) {

	if _, ok := n.Attr["ignore"]; ok {
		return
	}

	switch n.Attr["type"] {

	case "text":
		w.WriteString(n.rawStringLaTeX())

	case "paragraph":
		paragraphStartLaTeX(n, w)

	case "kanbun":
		kanbunStartLaTeX(n, w)

	case "ruby group":
		rubyGroupOpenLaTeX(n, w)

	case "ruby parent":
		rubyParentStartLaTeX(n, w)

	case "ruby":
		rubyOpenLaTeX(n, w)

	case "empty line":
		emptyLineLaTeX(n, w)

	case "section":
		return

	case "indentation":
		indentationOpenLaTeX(n, w)

	case "bottom align":
		bottomAlignOpenLaTeX(n, w)

	case "emphasis":
		emphOpenLaTeX(n, w)

	case "line decoration":
		lineDecorationOpenLaTeX(n, w)

	case "narrow paragraph":
		narrowParOpenLaTeX(n, w)

	case "inline section":
		inlineHeaderOpenLaTeX(n, w)

	case "window section":
		inlineHeaderOpenLaTeX(n, w)

	case "inline note":
		inlineNoteOpenLaTeX(n, w)

	case "rubylike note":
		rubyLikeOpenLaTeX(n, w)

	case "caption":
		captionOpenLaTeX(n, w)

	case "font shape":
		fontShapeOpenLaTeX(n, w)

	case "font size":
		fontSizeOpenLaTeX(n, w)

	case "offset":
		subsupOpenLaTeX(n, w)

	case "text direction":
		textDirOpenLaTeX(n, w)

	case "section title":
		sectionTitleLaTeX(n, w)

	case "figure":
		figureOpenLaTeX(n, w)

	case "image":
		imageLaTeX(n, w)

	case "pagination":
		paginationLaTeX(n, w)

	case "accent start":
		return

	case "accent end":
		return

	case "kunten":
		kuntenLaTeX(n, w)

	case "okurigana":
		okuriganaLaTeX(n, w)

	case "centering":
		centeringOpenLaTeX(n, w)

	case "bibliographical info":
		bibInfoOpenLaTeX(n, w)

	case "kunoji":
		kunojiOpenLaTeX(n, w)

	case "gaiji note":
		gaijiNoteOpenLaTeX(n, w)

	case "warichu line break":

		w.WriteString("\n\n")

	case "gaiji char":
		gaijiCharOpenLaTeX(n, w)

	case "special char":
		specialCharOpenLaTeX(n, w)

	case "metadata":
		metadataOpenLaTeX(n, w)

	case "meta title":
		metaTitleOpenLaTeX(n, w)

	case "meta subtitle":
		metaSubtitleOpenLaTeX(n, w)

	case "meta contributor":
		metaContributorOpenLaTeX(n, w)

	case "main text":
		return

	case "end marker":
		return

	case "document":
		return

	case "unknown":
		unknownOpenLaTeX(n, w)

	case "unknown block type":
		unknownBlockOpenLaTeX(n, w)

	default:
		msglog.Println("Renderer: unknown node type: " + n.String())
		unknownOpenLaTeX(n, w)
	}

}

func azrLaTeXFormatterClose(n *Node, w *strings.Builder) {

	if _, ok := n.Attr["ignore"]; ok {
		return
	}

	switch n.Attr["type"] {

	case "text":
		return

	case "paragraph":
		if n.Attr["latex lbr"] == "true" {
			w.WriteString(`\\` + "\n")
			return
		}
		w.WriteString("\n")

		if !n.insideSingleLineCommand() {
			w.WriteString("\n")
		}

	case "kanbun":
		kanbunEndLaTeX(n, w)

	case "ruby group":
		rubyGroupCloseLaTeX(n, w)

	case "ruby parent":
		rubyParentCloseLaTeX(n, w)

	case "ruby":
		rubyCloseLaTeX(n, w)

	case "empty line":
		return

	case "section":
		return

	case "indentation":
		indentationCloseLaTeX(n, w)

	case "bottom align":
		bottomAlignCloseLaTeX(n, w)

	case "emphasis":
		standardCloserLaTeX(n, w)

	case "line decoration":
		lineDecorationCloseLaTeX(n, w)

	case "narrow paragraph":
		narrowParCloseLaTeX(n, w)

	case "inline section":
		standardCloserLaTeX(n, w)

	case "window section":
		standardCloserLaTeX(n, w)

	case "inline note":
		inlineNoteCloseLaTeX(n, w)

	case "rubylike note":
		rubyLikeCloseLaTeX(n, w)

	case "caption":
		captionCloseLaTeX(n, w)

	case "font shape":
		standardCloserLaTeX(n, w)

	case "font size":
		standardCloserLaTeX(n, w)

	case "offset":
		subsupCloseLaTeX(n, w)

	case "text direction":
		textDirCloseLaTeX(n, w)

	case "section title":
		sectionTitleEndLaTeX(n, w)

	case "figure":
		figureCloseLaTeX(n, w)

	case "image":
		return

	case "pagination":
		return

	case "accent start":
		return

	case "accent end":
		return

	case "kunten":
		return

	case "okurigana":
		return

	case "centering":
		centeringCloseLaTeX(n, w)

	case "bibliographical info":
		return

	case "gaiji char":
		return

	case "gaiji note":
		return

	case "warichu line break":
		return

	case "special char":
		return

	case "metadata":
		metadataCloseLaTeX(n, w)

	case "meta title":
		metaTitleCloseLaTeX(n, w)

	case "meta subtitle":
		metaSubtitleCloseLaTeX(n, w)

	case "meta contributor":
		metaContributorCloseLaTeX(n, w)

	case "main text":
		return

	case "end marker":
		return

	case "document":
		return

	case "unknown":
		return

	case "unknown block type":
		unknownBlockCloseLaTeX(n, w)

	default:
		return
	}

	return
}

func standardCloserLaTeX(n *Node, w *strings.Builder) {

	w.WriteString("}")

	if !n.hasGaijiWithin {
		return
	}

	w.WriteString(n.firstDescendantOfType("gaiji token").gaijiNoteString())

}

func paragraphStartLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(n.Attr["latex extra"])

}

func emphOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvKenten"))

	if n.decoOnLeft() {
		w.WriteString(latexOpt("left"))
	}

	switch n.Attr["style"] {

	case "傍点":
		w.WriteString(latexArg("sesame"))

	case "白ゴマ傍点":
		w.WriteString(latexArg("Sesame"))

	case "丸傍点":
		w.WriteString(latexArg("circle"))

	case "白丸傍点":
		w.WriteString(latexArg("Circle"))

	case "黒三角傍点":
		w.WriteString(latexArg("triangle"))

	case "白三角傍点":
		w.WriteString(latexArg("Triangle"))

	case "二重丸傍点":
		w.WriteString(latexArg("bullseye"))

	case "蛇の目傍点":
		w.WriteString(latexArg("fisheye"))

	case "ばつ傍点":
		w.WriteString(latexArg(`$\times$`))

	}

	w.WriteString("{")

}

func emphCloseLaTeX(n *Node, w *strings.Builder) {

	standardCloserLaTeX(n, w)

}

func rubyOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(`}{`)

}

func rubyCloseLaTeX(n *Node, w *strings.Builder) {

	standardCloserLaTeX(n, w)

}

func rubyLikeOpenLaTeX(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {
		if strings.HasPrefix(e.Attr["type"], "ruby") {
			doublesidedRubyOpenLaTeX(n, w)
			return
		}
	}

	w.WriteString(latexCmd(`jruby`))

	opts := "g"

	if n.decoOnLeft() {
		opts = opts + "S"
	}

	w.WriteString(latexOpt(opts))

	w.WriteString(`{`)

}

func rubyLikeCloseLaTeX(n *Node, w *strings.Builder) {

	if n.Attr["latex dblruby"] == "true" {

		doubleSidedRubyCloseLaTeX(n, w)

		return
	}

	w.WriteString(`}{`)

	//w.WriteString(getRefStrings(n.rawCloserString())[0])
	w.WriteString(n.tok.getRefStrings()[0])

	w.WriteString(`}`)

	if !n.hasGaijiWithin {
		return
	}

	w.WriteString(n.firstDescendantOfType("gaiji note").gaijiNoteString())
}

func doublesidedRubyOpenLaTeX(n *Node, w *strings.Builder) {

	n.SetAttr("latex dblruby", "true")

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "ruby group" {

			e.SetAttr("latex dblruby", "true")

			break

		}
	}

	w.WriteString(latexCmd(`truby`))

	opts := "g"

	w.WriteString(latexOpt(opts))

	return
}

func doubleSidedRubyCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(`{`)

	//w.WriteString(getRefStrings(n.rawCloserString())[0])
	w.WriteString(n.tok.getRefStrings()[0])

	w.WriteString(`}`)

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "gaiji note" {

			w.WriteString(e.gaijiNoteString())

		}
	}
}

func fontShapeOpenLaTeX(n *Node, w *strings.Builder) {

	switch n.Attr["font shape"] {
	case "太字":
		w.WriteString(latexCmd("azconvBold"))
	case "斜体":
		w.WriteString(latexCmd("azconvItalic"))
	}

	w.WriteString("{")

}

func fontSizeOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString("{")

	switch n.Attr["size"] {

	case "大きな文字":
		switch n.Attr["step"] {
		case "1":
			w.WriteString(latexCmd("relsize{1}"))
		case "2":
			w.WriteString(latexCmd("relsize{2}"))
		case "3":
			w.WriteString(latexCmd("relsize{3}"))
		case "4":
			w.WriteString(latexCmd("relsize{4}"))
		default:
			w.WriteString(latexCmd("relsize{5}"))
		}

	case "小さな文字":
		switch n.Attr["step"] {
		case "1":
			w.WriteString(latexCmd("relsize{-2}"))
		case "2":
			w.WriteString(latexCmd("relsize{-3}"))
		case "3":
			w.WriteString(latexCmd("relsize{-4}"))
		default:
			w.WriteString(latexCmd("relsize{-2}"))
		}
	}

	w.WriteString(" ")

}

func sectionTitleLaTeX(n *Node, w *strings.Builder) {

	switch n.sectionLevel() {

	case 1:
		w.WriteString(latexCmd("azconvSection"))

	case 2:
		w.WriteString(latexCmd("azconvSubsection"))

	case 3:
		w.WriteString(latexCmd("azconvSubsubsection"))

	default:
		w.WriteString(latexCmd("azconvSsubsubsection"))

	}

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "paragraph" {

			e.Attr["latex lbr"] = "true"

		}
	}

	w.WriteString("{")
}

func sectionTitleEndLaTeX(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func imageLaTeX(n *Node, w *strings.Builder) {

	if n.Attr["style"] == "inline" {
		inlineImageLaTeX(n, w)
		return
	}

	w.WriteString(latexCmd("azconvIncludegraphics"))

	wd, _ := strconv.Atoi(n.Attr["width"])

	hd, _ := strconv.Atoi(n.Attr["height"])

	wratio := float64(wd) / 195

	hratio := float64(hd) / 325

	switch {

	case !(wratio < 1) && !(wratio < hratio):

		w.WriteString(latexArg(`0.62\textheight`))

		w.WriteString(latexArg("!"))

	case !(hratio < 1) && !(hratio < wratio):

		w.WriteString(latexArg("!"))

		w.WriteString(latexArg(`0.75\textwidth`))

	default:

		wsf := float64(wd) / 226
		w.WriteString(latexArg(strconv.FormatFloat(wsf, 'f', 2, 64) + `\textheight`))

		w.WriteString(latexArg("!"))

	}

	w.WriteString(latexArg(n.Attr["file"]))

	w.WriteString("\n\n")

}

func inlineImageLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvIncludegraphicsInline"))

	w.WriteString(latexArg(n.Attr["file"]))

}

func indentationOpenLaTeX(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "section title" {

			return

		}
	}
	w.WriteString(latexCmd("begin"))

	w.WriteString(latexArg("azconvIndented"))

	w.WriteString(latexArg(n.Attr["top margin"] + `\zw`))

	w.WriteString(latexArg(n.Attr["indent"] + `\zw`))

	w.WriteString("\n\n")

}

func indentationCloseLaTeX(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "section title" {

			return

		}
	}

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvIndented"))

	w.WriteString("\n")

}

func bottomAlignOpenLaTeX(n *Node, w *strings.Builder) {

	if !n.isBlockFormat() {

		n2 := newNode("text")

		n2.SetAttr("raw", `\hfill `)

		n.addFirstChild(n2)

		w.WriteString("\n\n")

	}

	w.WriteString(latexCmd(`begin`))

	w.WriteString(latexArg("azconvAlignBottom"))

	w.WriteString(latexArg(n.Attr["bottom margin"] + `\zw`))

	w.WriteString("\n\n")

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "paragraph" {

			e.SetAttr("latex extra", `\hfill `)

		}
	}
}

func bottomAlignCloseLaTeX(n *Node, w *strings.Builder) {

	if !n.isBlockFormat() {

		n.firstChild.Remove()

		w.WriteString("\n\n")

	}

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvAlignBottom"))

	w.WriteString("\n")

}

func inlineNoteOpenLaTeX(n *Node, w *strings.Builder) {

	wt := new(strings.Builder)

	renderInnerTextOnlyNoRuby(n, wt)

	s := strings.Split(wt.String(), "［＃改行］")
	longest := 0
	length := 0

	if len(s) > 1 {
		for j := range s {
			if len([]rune(s[j])) > len([]rune(s[longest])) {
				longest = j
			}
		}
		length = len([]rune(s[longest]))
	} else {
		length = len([]rune(s[0])) / 2
		if len([]rune(s[0]))%2 == 1 {
			length = length + 1
		}
	}

	w.WriteString(latexCmd("azconvWarichuM"))
	w.WriteString(latexArg(strconv.Itoa(length) + `.05\zw`))

	w.WriteString(`{`)

}

func inlineNoteCloseLaTeX(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		delete(e.Attr, "ignore")

	}

	w.WriteString(`}`)

}

func inlineHeaderOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvInlineHeader"))

	w.WriteString(latexArg(n.Attr["raw"]))

	w.WriteString(`{`)

}

func inlineHeaderCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(`}`)

}

func kuntenLaTeX(n *Node, w *strings.Builder) {

	if ok, _ := n.withinScopeOfType("kanbun"); ok {

		kanbunKuntenLaTeX(n, w)

		return

	}

	if n.prev != nil && n.prev.Attr["type"] == "okurigana" {

		l := strconv.Itoa(len([]rune(n.prev.okuriganaString())))

		w.WriteString(`\scriptsize\hskip-` + l + `\zw\normalsize`)

	}

	w.WriteString(latexCmd("azconvKunten"))

	var str string

	switch n.rawStringLaTeX() {
	case "一レ":
		str = `\azconvIR`
	case "二レ":
		str = `\azconvNR`
	case "上レ":
		str = `\azconvJR`
	case "甲レ":
		str = `azconvKR`
	default:
		str = n.rawStringLaTeX()
	}

	w.WriteString(latexArg(str))

}

func kanbunKuntenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(`[`)

	w.WriteString(n.rawStringLaTeX())

	w.WriteString(`]`)

}

func okuriganaLaTeX(n *Node, w *strings.Builder) {

	if ok, _ := n.withinScopeOfType("kanbun"); ok {

		kanbunOkuriganaLaTeX(n, w)

		return

	}
	if n.prev != nil && n.prev.Attr["type"] == "kunten" {

		l := strconv.Itoa(len([]rune(n.prev.Attr["raw"])))

		w.WriteString(`\scriptsize\hskip-` + l + `\zw\normalsize`)

	}

	w.WriteString(latexCmd("azconvOkurigana"))

	w.WriteString(latexArg(n.okuriganaString()))

}

func kanbunOkuriganaLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(`{`)

	w.WriteString(n.okuriganaString())

	w.WriteString(`}`)

}

func lineDecorationOpenLaTeX(n *Node, w *strings.Builder) {

	if n.Attr["style"] == "罫囲み" {

		boxedTextOpenLaTeX(n, w)

		return
	}

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "ruby group" {

			log.Println("WARNING: line " + strconv.Itoa(e.tok.lineNumber()) + " ruby inside underlining not implemented. Dropping underlining.")

			return
		}
	}

	w.WriteString(latexCmd("azconvUline"))

	switch n.Attr["style"] {

	case "傍線":
		w.WriteString(latexArg("solid"))

	case "二重傍線":
		w.WriteString(latexArg("double"))

	case "鎖線":
		w.WriteString(latexArg("dotted"))

	case "破線":
		w.WriteString(latexArg("dashed"))

	case "波線":
		w.WriteString(latexArg("wavy"))

	}

	if n.Attr["position"] == "left" {

		w.WriteString(latexArg("over"))

	}

	if n.Attr["position"] == "right" {

		w.WriteString(latexArg("under"))

	}

	w.WriteString(`{`)

}

func boxedTextOpenLaTeX(n *Node, w *strings.Builder) {

	if n.isBlockFormat() {

		w.WriteString(latexCmd("begin"))

		w.WriteString(latexArg("azconvFramed"))

		if ok, e := n.withinScopeOfType("indentation"); ok {

			opt := "leftmargin=" + e.Attr["top margin"] + `\zw`

			if ok, e = n.withinScopeOfType("narrow paragraph"); ok {

				opt = opt + ", " + `rightmargin=\dimexpr\linewidth -` + e.Attr["width"] + `\zw -5\zw`

			}

			w.WriteString(latexOpt(opt))

		}

		w.WriteString("\n")

		return

	}

	w.WriteString(latexCmd("azconvFbox"))

	w.WriteString("{")

}

func lineDecorationCloseLaTeX(n *Node, w *strings.Builder) {

	if n.Attr["style"] == "罫囲み" {

		boxedTextCloseLaTeX(n, w)

		return

	}

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "ruby group" {
			return
		}
	}
	standardCloserLaTeX(n, w)
}

func boxedTextCloseLaTeX(n *Node, w *strings.Builder) {

	if n.isBlockFormat() {

		w.WriteString(latexCmd("end"))

		w.WriteString(latexArg("azconvFramed"))

		w.WriteString("\n")

		return

	}

	w.WriteString("}")

}

/*
func lineDecorationCloseLaTeX(n *Node, w *strings.Builder) {

	if n.Attr["style"] != "罫囲み" {

		standardCloserLaTeX(n, w)

	}

	boxedTextCloseLaTeX(n, w)

}
*/

func narrowParOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("begin"))

	w.WriteString(latexArg("azconvNarrow"))

	top := 0

	if ok, e := n.withinScopeOfType("indentation"); ok {

		top, _ = strconv.Atoi(e.Attr["top margin"])

	}

	width, _ := strconv.Atoi(n.Attr["width"])

	w.WriteString(latexArg(strconv.Itoa(width+top) + `\zw`))

	w.WriteString("\n")

}

func narrowParCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvNarrow"))

	w.WriteString("\n")

}

func subsupOpenLaTeX(n *Node, w *strings.Builder) {

	switch {

	case strings.HasSuffix(n.Attr["raw"], "行右小書き"):
		w.WriteString(latexCmd("azconvSup"))

	case strings.HasSuffix(n.Attr["raw"], "行左小書き"):
		w.WriteString(latexCmd("azconvSub"))

	case strings.HasSuffix(n.Attr["raw"], "上付き小文字"):
		w.WriteString(latexCmd("azconvSup"))

	case strings.HasSuffix(n.Attr["raw"], "下付き小文字"):
		w.WriteString(latexCmd("azconvSub"))
	}

	w.WriteString(`{`)

}

func textDirOpenLaTeX(n *Node, w *strings.Builder) {

	switch {

	case n.Attr["raw"] == "縦中横":

		w.WriteString(latexCmd("azconvTatechuyoko"))

	case strings.HasSuffix(n.Attr["raw"], "横組み"):

		if n.isBlockFormat() {

			if n.innerParagraphCount() > 10 {

				log.Println("WARNING: line ", strconv.Itoa(n.tok.lineNumber()), " yokogumi specified but ignoring. Please adjust manually.")

				return
			}

			log.Println("WARNING: line ", strconv.Itoa(n.tok.lineNumber()), " yokogumi block. Please make sure result is correct.")

			w.WriteString(latexCmd(`parbox`))

			w.WriteString(latexArg(`\textwidth`))

			w.WriteString(`{\utod`)

			w.WriteString("\n")

			return
		}

		w.WriteString(latexCmd("azconvYokogumi"))
	}

	w.WriteString(`{`)
}

func textDirCloseLaTeX(n *Node, w *strings.Builder) {

	if n.isBlockFormat() {

		if n.innerParagraphCount() > 10 {
			return
		}
	}

	standardCloserLaTeX(n, w)

}

func paginationLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvNewpage"))

	w.WriteString(latexArg(n.Attr["raw"]))

	w.WriteString("\n")

}

func subsupCloseLaTeX(n *Node, w *strings.Builder) {

	standardCloserLaTeX(n, w)

}

func centeringOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd(`vspace*`))

	w.WriteString(latexArg(latexCmd("fill")))

	w.WriteString("\n")

}

func centeringCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("vfill"))

	w.WriteString("\n")
}

func figureOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("begin"))

	w.WriteString(latexArg("azconvFigure"))

	w.WriteString("\n")

}

func figureCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvFigure"))

	w.WriteString("\n")

}

func bibInfoOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd(`newpage`) + "\n")

}

func bibInfoCloseLaTeX(n *Node, w *strings.Builder) {

	return

}

func gaijiNoteOpenLaTeX(n *Node, w *strings.Builder) {

	if n.withinNoteExclScope() {
		return
	}

	unknownOpenLaTeX(n, w)

}

func rubyParentStartLaTeX(n *Node, w *strings.Builder) {

	if n.Attr["latex dblruby"] == "true" {

		w.WriteString("{")

		return
	}

}

func rubyParentCloseLaTeX(n *Node, w *strings.Builder) {

	e := n.firstDescendantOfType("gaiji note")

	if e == nil {

		return
	}

	w.WriteString(e.gaijiNoteString())

}

func specialCharOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(n.rawStringLaTeX())

}

func gaijiCharOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(n.rawStringLaTeX())

}

func kunojiOpenLaTeX(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, `{`, n.rawStringLaTeX(), `}`)

}

func metadataOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd(`thispagestyle`))

	w.WriteString(latexArg(`empty`))

	w.WriteString("\n")

	centeringOpenLaTeX(n, w)

	w.WriteString(latexCmd(`vskip`))

	w.WriteString(`-3\zw` + "\n")

}

func metadataCloseLaTeX(n *Node, w *strings.Builder) {

	centeringCloseLaTeX(n, w)

	w.WriteString(latexCmd(`newpage`))

	w.WriteString("\n\n")

}

func metaTitleOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvTitle"))
	w.WriteString(`{`)

}

func metaTitleCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func metaSubtitleOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvSubtitle"))

	w.WriteString(`{`)
}

func metaSubtitleCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func metaContributorOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvContributor"))

	w.WriteString(`{`)

}

func metaContributorCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(`}` + "\n")

}

func captionOpenLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("caption"))

	w.WriteString(`{`)

}

func captionCloseLaTeX(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func accentOpenLaTeX(n *Node, w *strings.Builder) {

	accentOpenTxt(n, w)

}

func unknownOpenLaTeX(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, `{`, latexCmd(`small`), ` `)

	addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr, `}`)
	//addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr)

}

func unknownBlockOpenLaTeX(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, `{`, latexCmd(`small`), ` `)

	addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr, `}`, "\n")
}

func unknownBlockCloseLaTeX(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, `{`, latexCmd(`small`), ` `)

	addToStringsBuilder(w, noteStartStr, n.Attr["raw closer"], noteEndStr, `}`, "\n")
}

func documentStartLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("begin"))

	w.WriteString(latexArg("document"))

	w.WriteString("\n\n")

}

func documentEndLaTeX(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("document"))

}

func emptyLineLaTeX(n *Node, w *strings.Builder) {

	if n.Parent() != nil {
		if n.Parent().Attr["type"] == "centering" {
			return
		}
	}

	w.WriteString("\n")

}

func rubyGroupOpenLaTeX(n *Node, w *strings.Builder) {

	defer func() {
		for _, e := range linearizeDescendants(n) {

			e.Attr["ignore"] = "true"

		}
	}()

	monolithicParent := false

	for _, e := range linearizeDescendants(n.firstChild) {

		if n.Attr["latex dblruby"] == "true" {
			monolithicParent = true
			break
		}

		switch e.Attr["type"] {

		case "text", "gaiji char", "kunoji", "special char", "accent string":
			continue

		default:

			monolithicParent = true

			break
		}

		break
	}

	var rps, rs []string

	wt := new(strings.Builder)

	if monolithicParent {

		wt.WriteString(`{`)

		renderInnerLaTeX(n.firstChild, wt)

		wt.WriteString(`}`)

		n.Attr["ruby base"] = wt.String()

		wt.Reset()

		for _, e := range linearizeDescendants(n) {

			if e.Attr["type"] == "ruby" {

				renderInnerLaTeX(e, wt)

				n.Attr["ruby string"] = wt.String()

				break
			}
		}
	}

	if !monolithicParent {

		const oLatex = true

		n.splitRuby(oLatex)

		rps = strings.Split(n.Attr["ruby base"], "\t")

		if len(rps) == 1 {

			monolithicParent = true

		} else {

			log.Println("WARNING: line " + strconv.Itoa(n.tok.lineNumber()) + " ruby parent or ruby too long. Splitting.")
		}
	}

	rps = strings.Split(n.Attr["ruby base"], "\t")

	rs = strings.Split(n.Attr["ruby string"], "\t")

	for c := range len(rps) {

		if n.Attr["latex dblruby"] != "true" {

			w.WriteString(latexCmd(`jruby`))

			w.WriteString(latexOpt(`g`))
		}

		w.WriteString(latexArg(rps[c]))

		w.WriteString(latexArg(rs[c]))

	}

}

func rubyGroupCloseLaTeX(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		delete(e.Attr, "ignore")

		if e.Attr["type"] == "gaiji note" {

			addToStringsBuilder(w, `{`, latexCmd(`small`), " ", e.gaijiNoteString(), `}`)

		}
	}

}

func kanbunStartLaTeX(n *Node, w *strings.Builder) {

	if n.prev != nil && n.prev.Attr["type"] == "kanbun" {
		return
	}

	w.WriteString(latexCmd(`Kanbun`))

	w.WriteString("\n")

}

func kanbunEndLaTeX(n *Node, w *strings.Builder) {

	if n.next != nil && n.next.Attr["type"] == "kanbun" {
		w.WriteString("\n\n")
		return
	}

	w.WriteString("\n")

	w.WriteString(latexCmd(`EndKanbun`))

	w.WriteString("\n")

	w.WriteString(latexCmd(`printkanbun`))

	w.WriteString("\n")

}
