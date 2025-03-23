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

func renderHtml(n *node) string {

	return strings.TrimSuffix(stringify(n, azrHtmlFormatterOpen, azrHtmlFormatterClose), `<br>`)

}

func azrHtmlFormatterOpen(n *node) string {

	switch n.attr["type"] {

	case "text":
		return n.attr["raw"]

	case "paragraph":
		return "<p>"

	case "ruby parent":
		return newHtag("ruby").String()

	case "ruby":
		return rubyOpenHtml(n)

	case "empty line":
		return newHtag("br").String() + "\n"

	case "section":
		return "<section>\n"

	case "gaiji char":
		return gaijiCharOpenHtml(n)

	case "kunoji":
		return kunojiOpenHtml(n)

	case "indentation":
		return indentationOpenHtml(n)

	case "bottom align":
		return bottomAlignOpenHtml(n)

	case "emphasis":
		if o_compatible {
			return rubyLikeOpenHtml(n)
		}
		return emphOpenHtml(n)

	case "line decoration":
		return lineDecorationOpenHtml(n)

	case "narrow paragraph":
		return narrowParOpenHtml(n)

	case "inline section":
		return inlineHeaderOpenHtml(n)

	case "window section":
		return inlineHeaderOpenHtml(n)

	case "inline note":
		return inlineNoteOpenHtml(n)

	case "rubylike note":
		return rubyLikeOpenHtml(n)

	case "caption":
		return captionOpenHtml(n)

	//	return newHtag("figcaption").String() + "\n"

	case "font shape":
		return fontShapeOpenHtml(n)

	case "font size":
		return fontSizeOpenHtml(n)

	case "offset":
		return subsupOpenHtml(n)

	case "text direction":
		return textDirOpenHtml(n)

	case "section title":
		return sectionTitleHtml(n)

	case "figure":
		return figureOpenHtml(n)

	case "image":
		return imageHtml(n)

	case "pagination":
		return paginationHtml(n)

	case "kunten":
		return kuntenHtml(n)

	case "okurigana":

		return okuriganaHtml(n)

	case "centering":
		return centeringOpenHtml(n)

	case "bibliographical info":
		return bibInfoOpenHtml(n)

	case "gaiji note":
		return gaijiNoteOpenHtml(n)

	case "warichu line break":
		return "\n"

	case "special char":
		return specialCharOpenHtml(n)

	case "accent string":
		return accentOpenHtml(n)

	case "metadata":
		return metadataOpenHtml(n)

	case "meta title":
		return metaTitleOpenHtml(n)

	case "meta subtitle":
		return metaSubtitleOpenHtml(n)

	case "meta contributor":
		return metaContributorOpenHtml(n)

	case "document":
		return ""

	default:
		return n.attr["raw"]
	}

}

func azrHtmlFormatterClose(n *node) string {

	switch n.attr["type"] {

	case "text":
		return ""

	case "paragraph":
		return "</p>\n"

	case "ruby parent":
		return rubyParentCloseHtml(n)
	//	return newCloseHtag("ruby").String()

	case "ruby":
		return rubyCloseHtml(n)

	case "empty line":
		return ""

	case "section":
		return "</section>\n"

	case "indentation":
		return standardCloserHtml(n)

	case "bottom align":
		return standardCloserHtml(n)

	case "emphasis":
		if o_compatible {
			return rubyLikeCloseHtml(n)
		}
		return emphCloseHtml(n)

	case "line decoration":
		return standardCloserHtml(n)

	case "narrow paragraph":
		return standardCloserHtml(n)

	case "inline section":
		return standardCloserHtml(n)

	case "window section":
		return standardCloserHtml(n)

	case "inline note":
		return standardCloserHtml(n)

	case "rubylike note":
		return rubyLikeCloseHtml(n)

	case "caption":
		return captionCloseHtml(n)

	case "font shape":
		return standardCloserHtml(n)

	case "font size":
		return standardCloserHtml(n)

	case "offset":
		return subsupCloseHtml(n)

	case "text direction":
		return standardCloserHtml(n)

	case "section title":
		return sectionTitleEndHtml(n)

	case "figure":
		return figureCloseHtml(n)

	case "image":
		return ""

	case "pagination":
		return ""

	case "kunten":
		return ""

	case "okurigana":

		return ""

	case "centering":
		return standardCloserHtml(n)

	case "bibliographical info":
		return bibInfoCloseHtml(n)

	case "metadata":
		return metadataCloseHtml(n)

	case "meta title":
		return metaTitleCloseHtml(n)

	case "meta subtitle":
		return metaSubtitleCloseHtml(n)

	case "meta contributor":
		return metaContributorCloseHtml(n)

	case "document":
		return ""

	default:
		return ""
	}

	return ""
}

func standardCloserHtml(n *node) string {

	if n.isBlockFormat() {

		return newCloseHtag("div").String() + "\n"

	}

	if !n.hasGaijiWithin {
		return newCloseHtag("span").String()

	}

	return newCloseHtag("span").String() + n.firstDescendantOfType("gaiji token").gaijiNoteString()

}

func emphOpenHtml(n *node) string {

	h := newHtag("em")

	if n.decoOnLeft() {
		h.addClass("emphLeft")
	} else {
		h.addClass("emphRight")
	}

	ext := ""

	switch n.attr["style"] {

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

	return h.String()

}

func emphCloseHtml(n *node) string {

	if !n.hasGaijiWithin {

		return newCloseHtag("em").String()

	}

	return newCloseHtag("em").String() + n.firstDescendantOfType("gaiji note").gaijiNoteString()

}

func rubyLikeOpenHtml(n *node) string {

	h := newHtag("ruby")

	if n.decoOnLeft() {
		h.addClass("left")
	}

	return h.String()

}

func rubyLikeCloseHtml(n *node) string {

	h1 := newHtag("rt")

	output.Reset()

	if n.formatOfType(decoMarker) {

		dec := emphString(n.attr["style"])

		l := len([]rune(n.innerText()))

		for range l {

			output.WriteString(" " + dec + " ")

		}

	} else {

		output.WriteString(getRefStrings(n.attr["raw closer"])[0])

	}

	h1.setAfter(output.String())

	if !n.hasGaijiWithin {
		return h1.String() + newCloseHtag("rt").String() + newCloseHtag("ruby").String()

	}

	return h1.String() + newCloseHtag("rt").String() + newCloseHtag("ruby").String() + n.firstDescendantOfType("gaiji note").gaijiNoteString()
}

func bibinfostringHtml(n *node) string {

	if n.attr["raw"] == "底本：" {
		return n.attr["raw"]
	}

	return noteStartStr + n.attr["raw"] + noteEndStr + lineBreakStr

}

func fontShapeOpenHtml(n *node) string {

	h := newHtag("")

	if n.isBlockFormat() {
		h.setElement("div")
		h.setAfter("\n")
	} else {
		h.setElement("span")
	}

	switch n.attr["font shape"] {
	case "太字":
		h.addClass("bold")
	case "斜体":
		h.addClass("italic")
	}

	return h.String()
}

func fontSizeOpenHtml(n *node) string {

	h := newHtag("")

	if n.isBlockFormat() {
		h.setElement("div")
		h.setAfter("\n")
	} else {
		h.setElement("span")
	}

	switch n.attr["size"] {

	case "大きな文字":
		switch n.attr["step"] {
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
		switch n.attr["step"] {
		case "1":
			h.addClass("small")
		case "2":
			h.addClass("x-small")
		default:
			h.addClass("xx-small")
		}
	}

	return h.String()
}

func sectionTitleHtml(n *node) string {

	h := newHtag("")

	h.setID(n.attr["id"])

	switch n.headerLevel() {

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

	return h.String()
}

func sectionTitleEndHtml(n *node) string {

	h := newCloseHtag("")

	switch n.headerLevel() {

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

	return h.String()
}

func imageHtml(n *node) string {

	h := newHtag("img")

	h.addClass("illustration")

	h.addExtraKeyVal("width", n.attr["width"])

	h.addExtraKeyVal("height", n.attr["height"])

	h.addExtraKeyVal("src", n.attr["file"])

	h.addExtraKeyVal("alt", n.attr["alt text"])

	h.setAfter("\n")

	return h.String()

}

func indentationOpenHtml(n *node) string {

	h := newHtag("div")

	h.addClass("topMargin" + n.attr["top margin"])

	h.addClass("textIndent" + n.attr["indent"])

	h.setAfter("\n")

	return h.String()

}

func bottomAlignOpenHtml(n *node) string {

	h := newHtag("span")

	if n.isBlockFormat() {
		h.setElement("div")
		h.addClass("alignBottom")
		h.setAfter("\n")

	} else {

		h.addClass("flushBottom")

	}

	h.addClass("bottomMargin" + n.attr["bottom margin"])

	return h.String()

}

func inlineNoteOpenHtml(n *node) string {

	h := newHtag("span")

	h.addClass("warichu")

	s := strings.Split(n.innerText(), "\n")

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

	return h.String()

}

func inlineHeaderOpenHtml(n *node) string {

	h := newHtag("span")

	switch {

	case strings.HasPrefix(n.attr["raw"], "同行"):
		h.addClass("inlineheader")

	case strings.HasPrefix(n.attr["raw"], "窓"):
		h.addClass("windowheader")
	}

	switch {

	case strings.HasSuffix(n.attr["raw"], "大見出し"):
		h.addClass("sectionLevel")

	case strings.HasSuffix(n.attr["raw"], "中見出し"):
		h.addClass("subsectionLevel")

	case strings.HasSuffix(n.attr["raw"], "小見出し"):
		h.addClass("subsubsectionLevel")
	}

	h.setID(n.attr["id"])

	return h.String()

}

func inlineHeaderCloseHtml(n *node) string {

	h := newHtag("span")

	h.setClose()

	return h.String()

}

func kuntenHtml(n *node) string {

	h1 := newHtag("sub")

	h1.addClass("kunten")

	h1.setAfter(n.attr["raw"])

	h2 := newHtag("sub")

	h2.setClose()

	return h1.String() + h2.String()

}

func okuriganaHtml(n *node) string {

	h1 := newHtag("sup")

	h1.addClass("okurigana")

	h1.setAfter(n.okuriganaString())

	h2 := newHtag("sup")

	h2.setClose()

	return h1.String() + h2.String()

}

func lineDecorationOpenHtml(n *node) string {

	h := newHtag("span")

	switch n.attr["style"] {

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

	if n.attr["position"] == "left" {

		h.addClass("overline")

	}

	if n.attr["position"] == "right" {

		h.addClass("underline")

	}

	return h.String()

}

func narrowParOpenHtml(n *node) string {

	h := newHtag("div")

	h.addClass("narrow")

	h.addExtraKeyVal("style", `height: `+n.attr["width"]+`em;`)

	h.setAfter("\n")

	return h.String()

}

func subsupOpenHtml(n *node) string {

	h := newHtag("")

	switch {

	case strings.HasSuffix(n.attr["raw"], "行右小書き"):
		h.setElement("sup")
		h.addClass("superscript")

	case strings.HasSuffix(n.attr["raw"], "行左小書き"):
		h.setElement("sub")
		h.addClass("subscript")

	case strings.HasSuffix(n.attr["raw"], "上付き小文字"):
		h.setElement("sup")
		h.addClass("superscript")

	case strings.HasSuffix(n.attr["raw"], "下付き小文字"):
		h.setElement("sub")
		h.addClass("subscript")
	}

	return h.String()

}

func textDirOpenHtml(n *node) string {

	h := newHtag("span")

	if n.isBlockFormat() {
		h.setElement("div")

	}
	switch n.attr["raw"] {

	case "縦中横":
		h.addClass("horizontal")

	case "横組み":
		h.addClass("sideways")

	default:
		h.addClass("sideways")
		h.setAfter("\n")

	}

	return h.String()
}

func textDirCloseHtml(n *node) string {

	return standardCloserHtml(n)

}

func paginationHtml(n *node) string {

	h := newHtag("div")

	h.addClass("pageBreak")

	h.addExtraKeyVal("data-AmznPageBreak", "always")

	return h.String() + newCloseHtag("div").String() + "\n"

}

func subsupCloseHtml(n *node) string {

	if strings.HasSuffix(n.attr["raw"], "行右小書き") {
		return newCloseHtag("sup").String()

	}

	if strings.HasSuffix(n.attr["raw"], "上付き小文字") {
		return newCloseHtag("sup").String()
	}

	return newCloseHtag("sub").String()
}

func centeringOpenHtml(n *node) string {

	h := newHtag("div")

	h.addClass("centering")

	h.setAfter("\n")

	return h.String()

}

func figureOpenHtml(n *node) string {

	h := newHtag("figure")

	h.setID(n.attr["id"])

	h.setAfter("\n")

	return h.String()

}

func figureCloseHtml(n *node) string {

	h := newCloseHtag("figure")

	h.setAfter("\n")

	return h.String()

}

func bibInfoOpenHtml(n *node) string {

	h1 := newHtag("hr")
	h1.setAfter("\n")

	h2 := newHtag("footer")
	h2.setAfter("\n")

	return paginationHtml(n) + h1.String() + h2.String()

}

func bibInfoCloseHtml(n *node) string {

	return newCloseHtag("footer").String() + "\n"

}

func renderNavHtml(n *node) string {

	h1 := newHtag("ol")

	h1.addClass("toc")

	h1.setAfter("\n")

	h2 := newCloseHtag("ol")

	return h1.String() + strings.TrimSuffix(stringify(n, navOpenHtml, navCloseHtml), `<br>`) + h2.String()

}

func navOpenHtml(n *node) string {

	switch n.attr["type"] {

	case "section":

		t := newHtag("ol")

		t.setAfter("\n")

		switch n.attr["sectionLevel"] {
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
			return t0.String() + t.String()
		}

		if n.prev.attr["type"] != "section" {
			return t0.String() + t.String()
		}

		return ""

	case "section title":
		return n.tocEntryHtml()

	default:
		return ""
	}
}

func navCloseHtml(n *node) string {

	switch n.attr["type"] {

	case "section":

		t := newCloseHtag("ol")
		t.setAfter("\n")

		t2 := newCloseHtag("li")
		t2.setAfter("\n")

		if n.next == nil {
			return t.String() + t2.String()
		}

		if n.next.attr["type"] != "section" {
			return t.String() + t2.String()
		}

		return ""

	case "section title":
		return newCloseHtag("a").String() + newCloseHtag("li").String() + "\n"

	default:
		return ""

	}
}

func (n *node) tocEntryHtml() string {

	h := newHtag("li")

	//	h.addClass("toc")

	h2 := newHtag("a")

	h2.addExtraKeyVal("href", "#"+n.attr["id"])

	h2.setAfter(renderHtml(n.firstChild))

	return h.String() + h2.String()

}

func gaijiNoteOpenHtml(n *node) string {

	if n.withinNoteExclScope() {
		return ""
	}

	return noteStartStr + n.attr["raw"] + noteEndStr

}

func rubyParentCloseHtml(n *node) string {

	if !n.hasGaijiWithin {

		return newCloseHtag("ruby").String()

	}

	return newCloseHtag("ruby").String() + n.firstDescendantOfType("gaiji note").gaijiNoteString()

}

func (n *node) gaijiNoteString() string {

	if n == nil {
		return ""
	}

	return noteStartStr + n.attr["raw"] + noteEndStr

}

func specialCharOpenHtml(n *node) string {

	if o_jis0208 {

		return n.attr["raw"]
	}

	return n.attr["alt raw"]

}

func gaijiCharOpenHtml(n *node) string {

	if o_jis0208 {

		return n.attr["raw"]
	}

	return n.attr["alt raw"]

}

func kunojiOpenHtml(n *node) string {

	if o_jis0208 {

		return n.attr["raw"]
	}

	return n.attr["alt raw"]

}

func metadataOpenHtml(n *node) string {

	return ""

}

func metadataCloseHtml(n *node) string {

	return "\n\n"

}

func metaTitleOpenHtml(n *node) string {

	h := newHtag("h1")

	h.addClass("title")

	return h.String()

}

func metaTitleCloseHtml(n *node) string {

	h := newCloseHtag("h1")

	h.setAfter("\n")

	return h.String()

}

func metaSubtitleOpenHtml(n *node) string {

	h := newHtag("h2")

	h.addClass("subtitle")

	return h.String()

}

func metaSubtitleCloseHtml(n *node) string {

	h := newCloseHtag("h2")

	h.setAfter("\n")

	return h.String()

}

func metaContributorOpenHtml(n *node) string {

	h := newHtag("h2")

	h.addClass("contributor")

	return h.String()

}

func metaContributorCloseHtml(n *node) string {

	h := newCloseHtag("h2")

	h.setAfter("\n")

	return h.String()

}

func captionOpenHtml(n *node) string {

	h := newHtag("figcaption")

	h.setAfter("\n")

	return h.String()

}

func captionCloseHtml(n *node) string {

	h := newCloseHtag("figcaption")

	if n.attr["scope"] != "block" {

		h.setBefore("\n")

	}

	h.setAfter("\n")

	return h.String()

}

func rubyOpenHtml(n *node) string {

	h1 := newHtag("rp")

	h1.setAfter("（")

	h2 := newCloseHtag("rp")

	h3 := newHtag("rt")

	return h1.String() + h2.String() + h3.String()

}

func rubyCloseHtml(n *node) string {

	h1 := newCloseHtag("rt")

	h2 := newHtag("rp")

	h2.setAfter("）")

	h3 := newCloseHtag("rp")

	return h1.String() + h2.String() + h3.String()

}

func accentOpenHtml(n *node) string {

	return accentOpenTxt(n)

}
