package aozoratext

import (
	"log"
	"strconv"
	"strings"
)

var o_compatible bool

var output *strings.Builder

func init() {
	output = new(strings.Builder)
	o_compatible = false
}

func setCompatible() {
	o_compatible = true
}

func unsetCompatible() {
	o_compatible = true
}

func renderHtml(n *Node, w *strings.Builder) {

	Serialize(n, w, azrHtmlFormatterOpen, azrHtmlFormatterClose)

}

func azrHtmlFormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "text":
		w.WriteString(n.Attr["raw"])

	case "paragraph":
		w.WriteString("<p>")

	case "ruby parent":
		newHtag("ruby").AddStringTo(w)

	case "ruby":
		rubyOpenHtml(n, w)

	case "empty line":
		newHtag("br").AddStringTo(w)
		w.WriteString("\n")

	case "section":
		w.WriteString("<section>\n")

	case "gaiji char":
		gaijiCharOpenHtml(n, w)

	case "kunoji":
		kunojiOpenHtml(n, w)

	case "indentation":
		indentationOpenHtml(n, w)

	case "bottom align":
		bottomAlignOpenHtml(n, w)

	case "emphasis":
		if o_compatible {
			rubyLikeOpenHtml(n, w)
			return
		}
		emphOpenHtml(n, w)

	case "line decoration":
		lineDecorationOpenHtml(n, w)

	case "narrow paragraph":
		narrowParOpenHtml(n, w)

	case "inline section":
		inlineHeaderOpenHtml(n, w)

	case "window section":
		inlineHeaderOpenHtml(n, w)

	case "inline note":
		inlineNoteOpenHtml(n, w)

	case "rubylike note":
		rubyLikeOpenHtml(n, w)

	case "caption":
		captionOpenHtml(n, w)

	case "font shape":
		fontShapeOpenHtml(n, w)

	case "font size":
		fontSizeOpenHtml(n, w)

	case "offset":
		subsupOpenHtml(n, w)

	case "text direction":
		textDirOpenHtml(n, w)

	case "section title":
		sectionTitleHtml(n, w)

	case "figure":
		figureOpenHtml(n, w)

	case "image":
		imageHtml(n, w)

	case "pagination":
		paginationHtml(n, w)

	case "kunten":
		kuntenHtml(n, w)

	case "okurigana":

		okuriganaHtml(n, w)

	case "centering":
		centeringOpenHtml(n, w)

	case "bibliographical info":
		bibInfoOpenHtml(n, w)

	case "gaiji note":
		gaijiNoteOpenHtml(n, w)

	case "warichu line break":
		w.WriteString("\n")

	case "special char":
		specialCharOpenHtml(n, w)

	case "accent string":
		accentOpenHtml(n, w)

	case "metadata":
		metadataOpenHtml(n, w)

	case "meta title":
		metaTitleOpenHtml(n, w)

	case "meta subtitle":
		metaSubtitleOpenHtml(n, w)

	case "meta contributor":
		metaContributorOpenHtml(n, w)

	case "document":
		return

	default:
		unknownOpenHtml(n, w)
	}

}

func azrHtmlFormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "text":
		return

	case "paragraph":
		w.WriteString("</p>\n")

	case "ruby parent":
		rubyParentCloseHtml(n, w)

	case "ruby":
		rubyCloseHtml(n, w)

	case "empty line":
		return

	case "section":
		w.WriteString("</section>\n")

	case "indentation":
		standardCloserHtml(n, w)

	case "bottom align":
		standardCloserHtml(n, w)

	case "emphasis":
		if o_compatible {
			rubyLikeCloseHtml(n, w)
		}
		emphCloseHtml(n, w)

	case "line decoration":
		standardCloserHtml(n, w)

	case "narrow paragraph":
		standardCloserHtml(n, w)

	case "inline section":
		standardCloserHtml(n, w)

	case "window section":
		standardCloserHtml(n, w)

	case "inline note":
		standardCloserHtml(n, w)

	case "rubylike note":
		rubyLikeCloseHtml(n, w)

	case "caption":
		captionCloseHtml(n, w)

	case "font shape":
		standardCloserHtml(n, w)

	case "font size":
		standardCloserHtml(n, w)

	case "offset":
		subsupCloseHtml(n, w)

	case "text direction":
		standardCloserHtml(n, w)

	case "section title":
		sectionTitleEndHtml(n, w)

	case "figure":
		figureCloseHtml(n, w)

	case "image":
		return

	case "pagination":
		return

	case "kunten":
		return

	case "okurigana":

		return

	case "centering":
		standardCloserHtml(n, w)

	case "bibliographical info":
		bibInfoCloseHtml(n, w)

	case "metadata":
		metadataCloseHtml(n, w)

	case "meta title":
		metaTitleCloseHtml(n, w)

	case "meta subtitle":
		metaSubtitleCloseHtml(n, w)

	case "meta contributor":
		metaContributorCloseHtml(n, w)

	case "document":
		return

	default:
		return
	}

	return
}

func standardCloserHtml(n *Node, w *strings.Builder) {

	if n.isBlockFormat() {

		newCloseHtag("div").AddStringTo(w)
		w.WriteString("\n")
		return
	}

	if !n.hasGaijiWithin {
		newCloseHtag("span").AddStringTo(w)
		return
	}

	newCloseHtag("span").AddStringTo(w)
	w.WriteString(n.firstDescendantOfType("gaiji token").gaijiNoteString())

}

func emphOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("em")

	if n.decoOnLeft() {
		h.addClass("emphLeft")
	} else {
		h.addClass("emphRight")
	}

	ext := ""

	switch n.Attr["style"] {

	case "傍点":
		ext = "goma"

	case "白ゴマ傍点":
		ext = "shirogoma"

	case "白丸傍点":
		ext = "shiromaru"

	case "黒三角傍点":
		ext = "sankaku"

	case "白三角傍点":
		ext = "shirosankaku"

	case "二重丸傍点":
		ext = "nijuumaru"

	case "蛇の目傍点":
		ext = "janome"

	case "ばつ傍点":
		ext = "batsu"

	default:
		ext = "maru"
	}

	h.addClass(ext)

	h.AddStringTo(w)

}

func emphCloseHtml(n *Node, w *strings.Builder) {

	if !n.hasGaijiWithin {

		newCloseHtag("em").AddStringTo(w)
		return

	}

	newCloseHtag("em").AddStringTo(w)
	w.WriteString(n.firstDescendantOfType("gaiji note").gaijiNoteString())

}

func rubyLikeOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("ruby")

	if n.decoOnLeft() {
		h.addClass("left")
	}

	h.AddStringTo(w)

}

func rubyLikeCloseHtml(n *Node, w *strings.Builder) {

	h1 := newHtag("rt")

	output.Reset()

	if n.tok.isFormatOfType(decoMarker) {

		dec := emphString(n.Attr["style"])

		wt := new(strings.Builder)

		renderInnerAozoraText(n, wt)

		l := len([]rune(wt.String()))

		for range l {

			output.WriteString(" " + dec + " ")

		}

	} else {

		output.WriteString(getRefStrings(n.RawCloserString())[0])

	}

	h1.setAfter(output.String())

	if !n.hasGaijiWithin {
		h1.AddStringTo(w)
		newCloseHtag("rt").AddStringTo(w)
		newCloseHtag("ruby").AddStringTo(w)
		return
	}

	h1.AddStringTo(w)
	newCloseHtag("rt").AddStringTo(w)
	newCloseHtag("ruby").AddStringTo(w)
	w.WriteString(n.firstDescendantOfType("gaiji note").gaijiNoteString())
}

func bibinfostringHtml(n *Node, w *strings.Builder) {

	if n.Attr["raw"] == "底本：" {
		w.WriteString(n.Attr["raw"])
		return
	}

	addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr, lineBreakStr)

}

func fontShapeOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("")

	if n.isBlockFormat() {
		h.setElement("div")
		h.setAfter("\n")
	} else {
		h.setElement("span")
	}

	switch n.Attr["font shape"] {
	case "太字":
		h.addClass("bold")
	case "斜体":
		h.addClass("italic")
	}

	h.AddStringTo(w)
}

func fontSizeOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("")

	if n.isBlockFormat() {
		h.setElement("div")
		h.setAfter("\n")
	} else {
		h.setElement("span")
	}

	switch n.Attr["size"] {

	case "大きな文字":
		switch n.Attr["step"] {
		case "1":
			h.addClass("large")
		case "2":
			h.addClass("x-large")
		case "3":
			h.addClass("xx-large")
		default:
			h.addClass("xxx-large")
		}

	case "小さな文字":
		switch n.Attr["step"] {
		case "1":
			h.addClass("small")
		case "2":
			h.addClass("x-small")
		default:
			h.addClass("xx-small")
		}
	}

	h.AddStringTo(w)
}

func sectionTitleHtml(n *Node, w *strings.Builder) {

	h := newHtag("")

	h.setID(n.Attr["id"])

	switch n.sectionLevel() {

	case 1:
		h.setElement("h3")
		h.addClass("section")

	case 2:
		h.setElement("h4")
		h.addClass("subsection")

	case 3:
		h.setElement("h5")
		h.addClass("subsubsection")

	default:
		h.setElement("h6")
		h.addClass("subsubsection")
		log.Println("unsupported header level", n.String())

	}

	h.AddStringTo(w)
}

func sectionTitleEndHtml(n *Node, w *strings.Builder) {

	h := newCloseHtag("")

	switch n.sectionLevel() {

	case 1:
		h.setElement("h3")

	case 2:
		h.setElement("h4")

	case 3:
		h.setElement("h5")

	default:
		h.setElement("h6")

	}

	h.setAfter("\n")

	h.AddStringTo(w)
}

func imageHtml(n *Node, w *strings.Builder) {

	h := newHtag("img")

	h.addClass("illustration")

	h.addExtraKeyVal("width", n.Attr["width"])

	h.addExtraKeyVal("height", n.Attr["height"])

	h.addExtraKeyVal("src", n.Attr["file"])

	h.addExtraKeyVal("alt", n.Attr["alt text"])

	h.setAfter("\n")

	h.AddStringTo(w)

}

func indentationOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("topMargin" + n.Attr["top margin"])

	h.addClass("textIndent" + n.Attr["indent"])

	h.setAfter("\n")

	h.AddStringTo(w)

}

func bottomAlignOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("span")

	if n.isBlockFormat() {
		h.setElement("div")
		h.addClass("alignBottom")
		h.setAfter("\n")

	} else {

		h.addClass("flushBottom")

	}

	h.addClass("bottomMargin" + n.Attr["bottom margin"])

	h.AddStringTo(w)

}

func inlineNoteOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("span")

	h.addClass("warichu")

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

	h.addExtraKeyVal("style", "height: "+strconv.Itoa(l)+".5em;")

	h.AddStringTo(w)

}

func inlineHeaderOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("span")

	switch {

	case strings.HasPrefix(n.Attr["raw"], "同行"):
		h.addClass("inlineheader")

	case strings.HasPrefix(n.Attr["raw"], "窓"):
		h.addClass("windowheader")
	}

	switch {

	case strings.HasSuffix(n.Attr["raw"], "大見出し"):
		h.addClass("sectionLevel")

	case strings.HasSuffix(n.Attr["raw"], "中見出し"):
		h.addClass("subsectionLevel")

	case strings.HasSuffix(n.Attr["raw"], "小見出し"):
		h.addClass("subsubsectionLevel")
	}

	h.setID(n.Attr["id"])

	h.AddStringTo(w)

}

func inlineHeaderCloseHtml(n *Node, w *strings.Builder) {

	h := newHtag("span")

	h.setClose()

	h.AddStringTo(w)

}

func kuntenHtml(n *Node, w *strings.Builder) {

	h1 := newHtag("sub")

	h1.addClass("kunten")

	h1.setAfter(n.Attr["raw"])

	h2 := newHtag("sub")

	h2.setClose()

	h1.AddStringTo(w)
	h2.AddStringTo(w)

}

func okuriganaHtml(n *Node, w *strings.Builder) {

	h1 := newHtag("sup")

	h1.addClass("okurigana")

	h1.setAfter(n.okuriganaString())

	h2 := newHtag("sup")

	h2.setClose()

	h1.AddStringTo(w)
	h2.AddStringTo(w)

}

func lineDecorationOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("span")

	switch n.Attr["style"] {

	case "傍線":
		h.addClass("solid")

	case "二重傍線":
		h.addClass("double")

	case "鎖線":
		h.addClass("dotted")

	case "破線":
		h.addClass("dashed")

	case "波線":
		h.addClass("wavy")

	case "罫囲み":
		h.addClass("framed")
		if n.isBlockFormat() {
			h.setElement("div")
			h.setAfter("\n")
		}
	}

	if n.Attr["position"] == "left" {

		h.addClass("overline")

	}

	if n.Attr["position"] == "right" {

		h.addClass("underline")

	}

	h.AddStringTo(w)

}

func narrowParOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("narrow")

	h.addExtraKeyVal("style", `height: `+n.Attr["width"]+`em;`)

	h.setAfter("\n")

	h.AddStringTo(w)

}

func subsupOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("")

	switch {

	case strings.HasSuffix(n.Attr["raw"], "行右小書き"):
		h.setElement("sup")
		h.addClass("superscript")

	case strings.HasSuffix(n.Attr["raw"], "行左小書き"):
		h.setElement("sub")
		h.addClass("subscript")

	case strings.HasSuffix(n.Attr["raw"], "上付き小文字"):
		h.setElement("sup")
		h.addClass("superscript")

	case strings.HasSuffix(n.Attr["raw"], "下付き小文字"):
		h.setElement("sub")
		h.addClass("subscript")
	}

	h.AddStringTo(w)

}

func textDirOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("span")

	if n.isBlockFormat() {
		h.setElement("div")

	}
	switch n.Attr["raw"] {

	case "縦中横":
		h.addClass("horizontal")

	case "横組み":
		h.addClass("sideways")

	default:
		h.addClass("sideways")
		h.setAfter("\n")

	}

	h.AddStringTo(w)
}

func textDirCloseHtml(n *Node, w *strings.Builder) {

	standardCloserHtml(n, w)

}

func paginationHtml(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("pageBreak")

	h.addExtraKeyVal("data-AmznPageBreak", "always")

	h.AddStringTo(w)
	newCloseHtag("div").AddStringTo(w)
	w.WriteString("\n")

}

func subsupCloseHtml(n *Node, w *strings.Builder) {

	if strings.HasSuffix(n.Attr["raw"], "行右小書き") {
		newCloseHtag("sup").AddStringTo(w)
		return

	}

	if strings.HasSuffix(n.Attr["raw"], "上付き小文字") {
		newCloseHtag("sup").AddStringTo(w)
		return
	}

	newCloseHtag("sub").AddStringTo(w)
}

func centeringOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("centering")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func figureOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("figure")

	h.setID(n.Attr["id"])

	h.setAfter("\n")

	h.AddStringTo(w)

}

func figureCloseHtml(n *Node, w *strings.Builder) {

	h := newCloseHtag("figure")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func bibInfoOpenHtml(n *Node, w *strings.Builder) {

	h1 := newHtag("hr")
	h1.setAfter("\n")

	h2 := newHtag("footer")
	h2.setAfter("\n")

	paginationHtml(n, w)
	h1.AddStringTo(w)
	h2.AddStringTo(w)

}

func bibInfoCloseHtml(n *Node, w *strings.Builder) {

	newCloseHtag("footer").AddStringTo(w)
	w.WriteString("\n")

}

func renderNavHtml(n *Node, w *strings.Builder) {

	h1 := newHtag("ol")

	h1.addClass("toc")

	h1.setAfter("\n")

	h2 := newCloseHtag("ol")

	h1.AddStringTo(w)

	Serialize(n, w, navOpenHtml, navCloseHtml)

	h2.AddStringTo(w)

}

func navOpenHtml(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "section":

		t := newHtag("ol")

		t.setAfter("\n")

		switch n.Attr["sectionLevel"] {
		case "1":
			t.addClass("section")
		case "2":
			t.addClass("subsection")
		case "3":
			t.addClass("subsubsection")
		}

		t0 := newHtag("li")

		t0.setAfter("\n")

		if n.prev == nil {
			t0.AddStringTo(w)
			t.AddStringTo(w)
			return
		}

		if n.prev.Attr["type"] != "section" {
			t0.AddStringTo(w)
			t.AddStringTo(w)
			return
		}

		return

	case "section title":
		w.WriteString(n.tocEntryHtml())

	default:
		return
	}
}

func navCloseHtml(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "section":

		t := newCloseHtag("ol")
		t.setAfter("\n")

		t2 := newCloseHtag("li")
		t2.setAfter("\n")

		if n.next == nil {
			t.AddStringTo(w)
			t2.AddStringTo(w)
			return
		}

		if n.next.Attr["type"] != "section" {
			t.AddStringTo(w)
			t2.AddStringTo(w)
			return
		}

		return

	case "section title":
		newCloseHtag("a").AddStringTo(w)
		newCloseHtag("li").AddStringTo(w)
		w.WriteString("\n")
		return

	default:
		return

	}
}

func (n *Node) tocEntryHtml() string {

	h := newHtag("li")

	//	h.addClass("toc")

	h2 := newHtag("a")

	h2.addExtraKeyVal("href", "#"+n.Attr["id"])

	wt := new(strings.Builder)

	renderHtml(n.firstChild, wt)

	h2.setAfter(wt.String())

	return h.String() + h2.String()

}

func gaijiNoteOpenHtml(n *Node, w *strings.Builder) {

	if n.withinNoteExclScope() {
		return
	}

	addToStringsBuilder(w, noteStartStr, n.RawString(), noteEndStr)

}

func rubyParentCloseHtml(n *Node, w *strings.Builder) {

	if !n.hasGaijiWithin {

		newCloseHtag("ruby").AddStringTo(w)

		return
	}

	newCloseHtag("ruby").AddStringTo(w)
	w.WriteString(n.firstDescendantOfType("gaiji note").gaijiNoteString())

}

func (n *Node) gaijiNoteString() string {

	if n == nil {
		return ""
	}

	return noteStartStr + n.Attr["raw"] + noteEndStr

}

func specialCharOpenHtml(n *Node, w *strings.Builder) {

	w.WriteString(n.RawString())

}

func gaijiCharOpenHtml(n *Node, w *strings.Builder) {

	w.WriteString(n.RawString())

}

func kunojiOpenHtml(n *Node, w *strings.Builder) {

	w.WriteString(n.RawString())

}

func metadataOpenHtml(n *Node, w *strings.Builder) {

	return

}

func metadataCloseHtml(n *Node, w *strings.Builder) {

	w.WriteString("\n\n")

}

func metaTitleOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("h1")

	h.addClass("title")

	h.AddStringTo(w)

}

func metaTitleCloseHtml(n *Node, w *strings.Builder) {

	h := newCloseHtag("h1")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func metaSubtitleOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("h2")

	h.addClass("subtitle")

	h.AddStringTo(w)

}

func metaSubtitleCloseHtml(n *Node, w *strings.Builder) {

	h := newCloseHtag("h2")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func metaContributorOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("h2")

	h.addClass("contributor")

	h.AddStringTo(w)

}

func metaContributorCloseHtml(n *Node, w *strings.Builder) {

	h := newCloseHtag("h2")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func captionOpenHtml(n *Node, w *strings.Builder) {

	h := newHtag("figcaption")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func captionCloseHtml(n *Node, w *strings.Builder) {

	h := newCloseHtag("figcaption")

	if n.Attr["scope"] != "block" {

		h.setBefore("\n")

	}

	h.setAfter("\n")

	h.AddStringTo(w)

}

func rubyOpenHtml(n *Node, w *strings.Builder) {

	h1 := newHtag("rp")

	h1.setAfter("（")

	h2 := newCloseHtag("rp")

	h3 := newHtag("rt")

	h1.AddStringTo(w)
	h2.AddStringTo(w)
	h3.AddStringTo(w)

}

func rubyCloseHtml(n *Node, w *strings.Builder) {

	h1 := newCloseHtag("rt")

	h2 := newHtag("rp")

	h2.setAfter("）")

	h3 := newCloseHtag("rp")

	h1.AddStringTo(w)
	h2.AddStringTo(w)
	h3.AddStringTo(w)

}

func accentOpenHtml(n *Node, w *strings.Builder) {

	accentOpenTxt(n, w)

}

func unknownOpenHtml(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr)

}
