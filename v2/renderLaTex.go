package aozoraConvert

import (
	"log"
	"strconv"
	"strings"
)

func renderLaTeX(n *Node, w *strings.Builder) {

	Serialize(n, w, azrLaTexFormatterOpen, azrLaTexFormatterClose)

}

func renderInnerLaTeX(n *Node, w *strings.Builder) {

	SerializeDescendants(n, w, azrLaTexFormatterOpen, azrLaTexFormatterClose)

}

func azrLaTexFormatterOpen(n *Node, w *strings.Builder) {

	if _, ok := n.Attr["ignore"]; ok {
		return
	}

	switch n.Attr["type"] {

	case "text":
		w.WriteString(n.RawString())

	case "paragraph":
		paragraphStartLaTex(n, w)

	case "kanbun":
		kanbunStartLaTex(n, w)

	case "ruby group":
		rubyGroupOpenLatex(n, w)

	case "ruby parent":
		rubyParentStartLatex(n, w)

	case "ruby":
		rubyOpenLaTex(n, w)

	case "empty line":
		emptyLineLatex(n, w)

	case "section":
		return

	case "gaiji char":
		gaijiCharOpenLaTex(n, w)

	case "kunoji":
		kunojiOpenLaTex(n, w)

	case "indentation":
		indentationOpenLaTex(n, w)

	case "bottom align":
		bottomAlignOpenLaTex(n, w)

	case "emphasis":
		emphOpenLaTex(n, w)

	case "line decoration":
		lineDecorationOpenLaTex(n, w)

	case "narrow paragraph":
		narrowParOpenLaTex(n, w)

	case "inline section":
		inlineHeaderOpenLaTex(n, w)

	case "window section":
		inlineHeaderOpenLaTex(n, w)

	case "inline note":
		inlineNoteOpenLaTex(n, w)

	case "rubylike note":
		rubyLikeOpenLaTex(n, w)

	case "caption":
		captionOpenLaTex(n, w)

	case "font shape":
		fontShapeOpenLaTex(n, w)

	case "font size":
		fontSizeOpenLaTex(n, w)

	case "offset":
		subsupOpenLaTex(n, w)

	case "text direction":
		textDirOpenLaTex(n, w)

	case "section title":
		sectionTitleLaTex(n, w)

	case "figure":
		figureOpenLaTex(n, w)

	case "image":
		imageLaTex(n, w)

	case "pagination":
		paginationLaTex(n, w)

	case "kunten":
		kuntenLaTex(n, w)

	case "okurigana":
		okuriganaLaTex(n, w)

	case "centering":
		centeringOpenLaTex(n, w)

	case "bibliographical info":
		bibInfoOpenLaTex(n, w)

	case "gaiji note":
		gaijiNoteOpenLaTex(n, w)

	case "warichu line break":
		w.WriteString(`\\`)

	case "special char":
		specialCharOpenLaTex(n, w)

	case "accent string":
		accentOpenLaTex(n, w)

	case "metadata":
		metadataOpenLaTex(n, w)

	case "meta title":
		metaTitleOpenLaTex(n, w)

	case "meta subtitle":
		metaSubtitleOpenLaTex(n, w)

	case "meta contributor":
		metaContributorOpenLaTex(n, w)

	case "main text":
		return

	case "document":
		return

	case "unknown":
		unknownOpenLaTex(n, w)

	default:
		log.Println("Renderer: unknown node type: " + n.String())
		unknownOpenLaTex(n, w)
	}

}

func azrLaTexFormatterClose(n *Node, w *strings.Builder) {

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
		kanbunEndLaTex(n, w)

	case "ruby group":
		rubyGroupCloseLatex(n, w)

	case "ruby parent":
		rubyParentCloseLaTex(n, w)

	case "ruby":
		rubyCloseLaTex(n, w)

	case "empty line":
		return

	case "section":
		return

	case "indentation":
		indentationCloseLaTex(n, w)

	case "bottom align":
		bottomAlignCloseLaTex(n, w)

	case "emphasis":
		standardCloserLaTeX(n, w)

	case "line decoration":
		lineDecorationCloseLaTex(n, w)

	case "narrow paragraph":
		narrowParCloseLaTex(n, w)

	case "inline section":
		standardCloserLaTeX(n, w)

	case "window section":
		standardCloserLaTeX(n, w)

	case "inline note":
		inlineNoteCloseLaTex(n, w)

	case "rubylike note":
		rubyLikeCloseLaTex(n, w)

	case "caption":
		captionCloseLaTex(n, w)

	case "font shape":
		standardCloserLaTeX(n, w)

	case "font size":
		standardCloserLaTeX(n, w)

	case "offset":
		subsupCloseLaTex(n, w)

	case "text direction":
		textDirCloseLaTex(n, w)

	case "section title":
		sectionTitleEndLaTex(n, w)

	case "figure":
		figureCloseLaTex(n, w)

	case "image":
		return

	case "pagination":
		return

	case "kunten":
		return

	case "okurigana":
		return

	case "centering":
		centeringCloseLaTex(n, w)

	case "bibliographical info":
		return

	case "metadata":
		metadataCloseLaTex(n, w)

	case "meta title":
		metaTitleCloseLaTex(n, w)

	case "meta subtitle":
		metaSubtitleCloseLaTex(n, w)

	case "meta contributor":
		metaContributorCloseLaTex(n, w)

	case "main text":
		return

	case "document":
		return

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

func paragraphStartLaTex(n *Node, w *strings.Builder) {

	w.WriteString(n.Attr["latex extra"])

}

func emphOpenLaTex(n *Node, w *strings.Builder) {

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
		w.WriteString(latexArg("Ｘ"))

	}

	w.WriteString("{")

}

func emphCloseLaTex(n *Node, w *strings.Builder) {

	standardCloserLaTeX(n, w)

}

func rubyOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(`}{`)

}

func rubyCloseLaTex(n *Node, w *strings.Builder) {

	standardCloserLaTeX(n, w)

}

func rubyLikeOpenLaTex(n *Node, w *strings.Builder) {

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

func rubyLikeCloseLaTex(n *Node, w *strings.Builder) {

	if n.Attr["latex dblruby"] == "true" {

		doubleSidedRubyCloseLaTeX(n, w)

		return
	}

	w.WriteString(`}{`)

	w.WriteString(getRefStrings(n.RawCloserString())[0])

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

	w.WriteString(getRefStrings(n.RawCloserString())[0])

	w.WriteString(`}`)

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "gaiji note" {

			w.WriteString(e.gaijiNoteString())

		}
	}
}

func fontShapeOpenLaTex(n *Node, w *strings.Builder) {

	switch n.Attr["font shape"] {
	case "太字":
		w.WriteString(latexCmd("azconvBold"))
	case "斜体":
		w.WriteString(latexCmd("azconvItalic"))
	}

	w.WriteString("{")

}

func fontSizeOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString("{")

	switch n.Attr["size"] {

	case "大きな文字":
		switch n.Attr["step"] {
		case "1":
			w.WriteString(latexCmd("large"))
		case "2":
			w.WriteString(latexCmd("Large"))
		case "3":
			w.WriteString(latexCmd("LARGE"))
		case "4":
			w.WriteString(latexCmd("huge"))
		default:
			w.WriteString(latexCmd("Huge"))
		}

	case "小さな文字":
		switch n.Attr["step"] {
		case "1":
			w.WriteString(latexCmd("small"))
		case "2":
			w.WriteString(latexCmd("footnotesize"))
		case "3":
			w.WriteString(latexCmd("scriptsize"))
		default:
			w.WriteString(latexCmd("tiny"))
		}
	}

	w.WriteString(" ")

}

func sectionTitleLaTex(n *Node, w *strings.Builder) {

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

func sectionTitleEndLaTex(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func imageLaTex(n *Node, w *strings.Builder) {

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

func indentationOpenLaTex(n *Node, w *strings.Builder) {

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

func indentationCloseLaTex(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "section title" {

			return

		}
	}

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvIndented"))

	w.WriteString("\n")

}

func bottomAlignOpenLaTex(n *Node, w *strings.Builder) {

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

func bottomAlignCloseLaTex(n *Node, w *strings.Builder) {

	if !n.isBlockFormat() {

		n.firstChild.remove()

		w.WriteString("\n\n")

	}

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvAlignBottom"))

	w.WriteString("\n")

}

func inlineNoteOpenLaTex(n *Node, w *strings.Builder) {

	wt := new(strings.Builder)

	renderInnerAozoraText(n, wt)

	s := strings.Split(wt.String(), "［＃改行］")

	l := 0

	if len(s) > 1 {

		if len([]rune(s[0])) > len([]rune(s[1])) {

			l = len([]rune(s[0]))

		} else {

			l = len([]rune(s[1]))

		}
	} else {

		l = len([]rune(s[0]))

		if l%2 == 1 {
			l++
		}

		l = l / 2

	}

	/*
		if len(s) == 2 {

			w.WriteString(latexCmd("azconvWarichuM"))

			w.WriteString(latexArg(strconv.Itoa(l) + `.5\zw`))

			addToStringsBuilder(w, `{(`, s[0], `)`, `\\(`, s[1], `)`)

			for _, e := range linearizeDescendants(n) {

				e.SetAttr("ignore", "true")

			}

			return
		}
	*/

	w.WriteString(latexCmd("azconvWarichu"))

	w.WriteString(latexArg(strconv.Itoa(l) + `.5\zw`))

	w.WriteString(`{`)

}

func inlineNoteCloseLaTex(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		delete(e.Attr, "ignore")

	}

	w.WriteString(`}`)

}

func inlineHeaderOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvInlineHeader"))

	w.WriteString(latexArg(n.Attr["raw"]))

	w.WriteString(`{`)

}

func inlineHeaderCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString(`}`)

}

func kuntenLaTex(n *Node, w *strings.Builder) {

	if ok, _ := n.withinScopeOfType("kanbun"); ok {

		kanbunKuntenLaTex(n, w)

		return

	}

	if n.prev != nil && n.prev.Attr["type"] == "okurigana" {

		l := strconv.Itoa(len([]rune(n.prev.okuriganaString())))

		w.WriteString(`\scriptsize\hskip-` + l + `\zw\normalsize`)

	}

	w.WriteString(latexCmd("azconvKunten"))

	var str string

	switch n.RawString() {
	case "一レ":
		str = `\azconvIR`
	case "二レ":
		str = `\azconvNR`
	case "上レ":
		str = `\azconvJR`
	case "甲レ":
		str = `azconvKR`
	default:
		str = n.RawString()
	}

	w.WriteString(latexArg(str))

}

func kanbunKuntenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(`[`)

	w.WriteString(n.RawString())

	w.WriteString(`]`)

}

func okuriganaLaTex(n *Node, w *strings.Builder) {

	if ok, _ := n.withinScopeOfType("kanbun"); ok {

		kanbunOkuriganaLaTex(n, w)

		return

	}
	if n.prev != nil && n.prev.Attr["type"] == "kunten" {

		l := strconv.Itoa(len([]rune(n.prev.Attr["raw"])))

		w.WriteString(`\scriptsize\hskip-` + l + `\zw\normalsize`)

	}

	w.WriteString(latexCmd("azconvOkurigana"))

	w.WriteString(latexArg(n.okuriganaString()))

}

func kanbunOkuriganaLaTex(n *Node, w *strings.Builder) {

	w.WriteString(`{`)

	w.WriteString(n.okuriganaString())

	w.WriteString(`}`)

}

func lineDecorationOpenLaTex(n *Node, w *strings.Builder) {

	if n.Attr["style"] == "罫囲み" {

		boxedTextOpenLaTex(n, w)

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

func boxedTextOpenLaTex(n *Node, w *strings.Builder) {

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

func lineDecorationCloseLaTex(n *Node, w *strings.Builder) {

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

func lineDecorationCloseLatex(n *Node, w *strings.Builder) {

	if n.Attr["style"] != "罫囲み" {

		standardCloserLaTeX(n, w)

	}

	boxedTextCloseLaTeX(n, w)

}

func narrowParOpenLaTex(n *Node, w *strings.Builder) {

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

func narrowParCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvNarrow"))

	w.WriteString("\n")

}

func subsupOpenLaTex(n *Node, w *strings.Builder) {

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

func textDirOpenLaTex(n *Node, w *strings.Builder) {

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
		} else {

			w.WriteString(latexCmd("azconvYokogumi"))
		}
	}

	w.WriteString(`{`)
}

func textDirCloseLaTex(n *Node, w *strings.Builder) {

	if n.isBlockFormat() {

		if n.innerParagraphCount() > 10 {
			return
		}
	}

	standardCloserLaTeX(n, w)

}

func paginationLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvNewpage"))

	w.WriteString(latexArg(n.Attr["raw"]))

	w.WriteString("\n")

}

func subsupCloseLaTex(n *Node, w *strings.Builder) {

	standardCloserLaTeX(n, w)

}

func centeringOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd(`vspace*`))

	w.WriteString(latexArg(latexCmd("fill")))

	w.WriteString("\n")

}

func centeringCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("vfill"))

	w.WriteString("\n")
}

func figureOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("begin"))

	w.WriteString(latexArg("azconvFigure"))

	w.WriteString("\n")

}

func figureCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("end"))

	w.WriteString(latexArg("azconvFigure"))

	w.WriteString("\n")

}

func bibInfoOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd(`newpage`) + "\n")

}

func bibInfoCloseLaTex(n *Node, w *strings.Builder) {

	return

}

func gaijiNoteOpenLaTex(n *Node, w *strings.Builder) {

	if n.withinNoteExclScope() {
		return
	}

	unknownOpenLaTex(n, w)

}

func rubyParentStartLatex(n *Node, w *strings.Builder) {

	if n.Attr["latex dblruby"] == "true" {

		w.WriteString("{")

		return
	}

}

func rubyParentCloseLaTex(n *Node, w *strings.Builder) {

	e := n.firstDescendantOfType("gaiji note")

	if e == nil {

		return
	}

	w.WriteString(e.gaijiNoteString())

}

func specialCharOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(n.RawString())

}

func gaijiCharOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(n.RawString())

}

func kunojiOpenLaTex(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, `{`, n.RawString(), `}`)

}

func metadataOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd(`thispagestyle`))

	w.WriteString(latexArg(`empty`))

	w.WriteString("\n")

	centeringOpenLaTex(n, w)

	w.WriteString(latexCmd(`vskip`))

	w.WriteString(`-3\zw` + "\n")

}

func metadataCloseLaTex(n *Node, w *strings.Builder) {

	centeringCloseLaTex(n, w)

	w.WriteString(latexCmd(`newpage`))

	w.WriteString("\n\n")

}

func metaTitleOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvTitle"))
	w.WriteString(`{`)

}

func metaTitleCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func metaSubtitleOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvSubtitle"))

	w.WriteString(`{`)
}

func metaSubtitleCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func metaContributorOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("azconvContributor"))

	w.WriteString(`{`)

}

func metaContributorCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString(`}` + "\n")

}

func captionOpenLaTex(n *Node, w *strings.Builder) {

	w.WriteString(latexCmd("caption"))

	w.WriteString(`{`)

}

func captionCloseLaTex(n *Node, w *strings.Builder) {

	w.WriteString("}\n")

}

func accentOpenLaTex(n *Node, w *strings.Builder) {

	accentOpenTxt(n, w)

}

func unknownOpenLaTex(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, `{`, latexCmd(`small`), ` `)

	addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr, `}`)
	//addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr)

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

func emptyLineLatex(n *Node, w *strings.Builder) {

	if n.Parent() != nil {
		if n.Parent().Attr["type"] == "centering" {
			return
		}
	}

	w.WriteString("\n")

}

func rubyGroupOpenLatex(n *Node, w *strings.Builder) {

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

		n.splitRuby()

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

func rubyGroupCloseLatex(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {

		delete(e.Attr, "ignore")

		if e.Attr["type"] == "gaiji note" {

			addToStringsBuilder(w, `{`, latexCmd(`small`), " ", e.gaijiNoteString(), `}`)

		}
	}

}

func kanbunStartLaTex(n *Node, w *strings.Builder) {

	if n.prev != nil && n.prev.Attr["type"] == "kanbun" {
		return
	}

	w.WriteString(latexCmd(`Kanbun`))

	w.WriteString("\n")

}

func kanbunEndLaTex(n *Node, w *strings.Builder) {

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
