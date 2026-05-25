package aozoraconvert

import (
	"log"
	"strconv"
	"strings"
)

var oCompatible bool

var output *strings.Builder

func init() {
	output = new(strings.Builder)
	oCompatible = false
}

func setCompatible() {
	oCompatible = true
}

func unsetCompatible() {
	oCompatible = true
}

func renderHTML(n *Node, w *strings.Builder) {

	Serialize(n, w, azrHTMLFormatterOpen, azrHTMLFormatterClose)

}

func azrHTMLFormatterOpen(n *Node, w *strings.Builder) {

	if _, ok := n.Attr["ignore"]; ok {
		return
	}

	switch n.Attr["type"] {

	case "text":
		w.WriteString(n.rawString())

	case "paragraph":
		w.WriteString("<p>")

	case "kanbun":
		w.WriteString("<p>")

	case "ruby group":
		rubyGroupOpenHTML(n, w)

	case "ruby parent":
		newHtag("ruby").AddStringTo(w)

	case "ruby":
		rubyOpenHTML(n, w)

	case "empty line":
		h := newHtag("br")
		h.AddStringTo(w)
		w.WriteString("\n")

	case "section":
		w.WriteString("<section>\n")

	case "indentation":
		indentationOpenHTML(n, w)

	case "bottom align":
		bottomAlignOpenHTML(n, w)

	case "emphasis":
		if oCompatible {
			rubyLikeOpenHTML(n, w)
			return
		}
		emphOpenHTML(n, w)

	case "line decoration":
		lineDecorationOpenHTML(n, w)

	case "narrow paragraph":
		narrowParOpenHTML(n, w)

	case "inline section":
		inlineHeaderOpenHTML(n, w)

	case "window section":
		inlineHeaderOpenHTML(n, w)

	case "inline note":
		inlineNoteOpenHTML(n, w)

	case "rubylike note":
		rubyLikeOpenHTML(n, w)

	case "caption":
		captionOpenHTML(n, w)

	case "font shape":
		fontShapeOpenHTML(n, w)

	case "font size":
		fontSizeOpenHTML(n, w)

	case "offset":
		subsupOpenHTML(n, w)

	case "text direction":
		textDirOpenHTML(n, w)

	case "section title":
		sectionTitleHTML(n, w)

	case "figure":
		figureOpenHTML(n, w)

	case "image":
		imageHTML(n, w)

	case "pagination":
		return

	case "accent start":
		return

	case "accent end":
		return

	case "kunten":
		kuntenHTML(n, w)

	case "okurigana":
		okuriganaHTML(n, w)

	case "centering":
		centeringOpenHTML(n, w)

	case "bibliographical info":
		bibInfoOpenHTML(n, w)

	case "kunoji":
		kunojiOpenHTML(n, w)

	case "gaiji note":
		gaijiNoteOpenHTML(n, w)

	case "warichu line break":
		h := newHtag("br")
		h.AddStringTo(w)

	case "gaiji char":
		gaijiCharOpenHTML(n, w)

	case "special char":
		specialCharOpenHTML(n, w)

	case "metadata":
		metadataOpenHTML(n, w)

	case "meta title":
		metaTitleOpenHTML(n, w)

	case "meta subtitle":
		metaSubtitleOpenHTML(n, w)

	case "meta contributor":
		metaContributorOpenHTML(n, w)

	case "main text":
		return

	case "end marker":
		return

	case "document":
		return

	case "unknown":
		unknownOpenHTML(n, w)

	case "unknown block type":
		unknownBlockOpenHTML(n, w)

	default:
		msglog.Println("Renderer: unknown node type: " + n.String())
		unknownOpenHTML(n, w)
	}

}

func azrHTMLFormatterClose(n *Node, w *strings.Builder) {

	if _, ok := n.Attr["ignore"]; ok {
		return
	}

	switch n.Attr["type"] {

	case "text":
		return

	case "paragraph":
		w.WriteString("</p>\n")

	case "kanbun":
		w.WriteString("</p>\n")

	case "ruby group":
		rubyGroupCloseHTML(n, w)

	case "ruby parent":
		rubyParentCloseHTML(n, w)

	case "ruby":
		rubyCloseHTML(n, w)

	case "empty line":
		return

	case "section":
		w.WriteString("</section>\n")

	case "indentation":
		standardCloserHTML(n, w)

	case "bottom align":
		standardCloserHTML(n, w)

	case "emphasis":
		if oCompatible {
			rubyLikeCloseHTML(n, w)
			return
		}
		emphCloseHTML(n, w)

	case "line decoration":
		standardCloserHTML(n, w)

	case "narrow paragraph":
		standardCloserHTML(n, w)

	case "inline section":
		standardCloserHTML(n, w)

	case "window section":
		standardCloserHTML(n, w)

	case "inline note":
		standardCloserHTML(n, w)

	case "rubylike note":
		rubyLikeCloseHTML(n, w)

	case "caption":
		captionCloseHTML(n, w)

	case "font shape":
		standardCloserHTML(n, w)

	case "font size":
		standardCloserHTML(n, w)

	case "offset":
		subsupCloseHTML(n, w)

	case "text direction":
		standardCloserHTML(n, w)

	case "section title":
		sectionTitleEndHTML(n, w)

	case "figure":
		figureCloseHTML(n, w)

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
		standardCloserHTML(n, w)

	case "bibliographical info":
		bibInfoCloseHTML(n, w)

	case "gaiji char":
		return

	case "kunoji":
		return

	case "gaiji note":
		return

	case "warichu line break":
		return

	case "special char":
		return

	case "metadata":
		metadataCloseHTML(n, w)

	case "meta title":
		metaTitleCloseHTML(n, w)

	case "meta subtitle":
		metaSubtitleCloseHTML(n, w)

	case "meta contributor":
		metaContributorCloseHTML(n, w)

	case "main text":
		return

	case "end marker":
		return

	case "document":
		return

	case "unknown":
		return

	case "unknown block type":
		return

	default:
		return
	}

	return
}

func standardCloserHTML(n *Node, w *strings.Builder) {

	if n.Attr["type"] == "metadata" {

		newCloseHtag("div").AddStringTo(w)
		w.WriteString("\n")
		return
	}

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

func emphOpenHTML(n *Node, w *strings.Builder) {

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

func emphCloseHTML(n *Node, w *strings.Builder) {

	if !n.hasGaijiWithin {

		newCloseHtag("em").AddStringTo(w)
		return

	}

	newCloseHtag("em").AddStringTo(w)
	w.WriteString(n.firstDescendantOfType("gaiji note").gaijiNoteString())

}

func rubyLikeOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("ruby")

	if n.decoOnLeft() {
		h.addClass("left")
	}

	h.AddStringTo(w)

}

func rubyLikeCloseHTML(n *Node, w *strings.Builder) {

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

		//output.WriteString(getRefStrings(n.rawCloserString())[0])
		output.WriteString(n.tok.getRefStrings()[0])

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

func bibinfostringHTML(n *Node, w *strings.Builder) {

	if n.Attr["raw"] == "底本：" {
		w.WriteString(n.Attr["raw"])
		return
	}

	addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr, lineBreakStr)

}

func fontShapeOpenHTML(n *Node, w *strings.Builder) {

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

func fontSizeOpenHTML(n *Node, w *strings.Builder) {

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

func sectionTitleHTML(n *Node, w *strings.Builder) {

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

func sectionTitleEndHTML(n *Node, w *strings.Builder) {

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

func imageHTML(n *Node, w *strings.Builder) {

	h := genImageTag(n)

	h.AddStringTo(w)

}

func genImageTag(n *Node) *htmlTagSpec {

	h := newHtag("img")

	if val, ok := n.Attr["style"]; ok {
		if val == "inline" {
			h.addClass("inlineImage")
		}
	} else {
		h.addClass("illustration")
	}

	h.addExtraKeyVal("max-width", n.Attr["width"]+"px")

	h.addExtraKeyVal("max-height", n.Attr["height"]+"px")

	//h.addExtraKeyVal("height", "100%")

	h.addExtraKeyVal("aspect-ratio", n.Attr["width"]+"/"+n.Attr["height"])

	h.addExtraKeyVal("src", n.Attr["file"])

	h.addExtraKeyVal("alt", n.Attr["alt text"])

	h.setAfter("\n")

	return h

}

func indentationOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("topMargin" + n.Attr["top margin"])

	h.addClass("textIndent" + n.Attr["indent"])

	h.setAfter("\n")

	h.AddStringTo(w)

}

func bottomAlignOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("span")

	if n.isBlockFormat() {
		h.setElement("div")
		h.addClass("alignBottom")
		h.setAfter("\n")

	} else {
		h1 := newHtag("br")
		h1.AddStringTo(w)

		h.addClass("flushBottom")

	}

	h.addClass("bottomMargin" + n.Attr["bottom margin"])

	h.AddStringTo(w)

}

func inlineNoteOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("span")

	h.addClass("warichu")

	wt := new(strings.Builder)

	renderInnerAozoraText(n, wt)

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
		h.addExtraKeyVal("style", "height: "+strconv.Itoa(length)+".5em;")
	}

	h.AddStringTo(w)

}

func inlineHeaderOpenHTML(n *Node, w *strings.Builder) {

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

func inlineHeaderCloseHTML(n *Node, w *strings.Builder) {

	h := newHtag("span")

	h.setClose()

	h.AddStringTo(w)

}

func kuntenHTML(n *Node, w *strings.Builder) {

	h1 := newHtag("sub")

	h1.addClass("kunten")

	h1.setAfter(n.Attr["raw"])

	h2 := newHtag("sub")

	h2.setClose()

	h1.AddStringTo(w)
	h2.AddStringTo(w)

}

func okuriganaHTML(n *Node, w *strings.Builder) {

	h1 := newHtag("sup")

	h1.addClass("okurigana")

	h1.setAfter(n.okuriganaString())

	h2 := newHtag("sup")

	h2.setClose()

	h1.AddStringTo(w)
	h2.AddStringTo(w)

}

func lineDecorationOpenHTML(n *Node, w *strings.Builder) {

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

func narrowParOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("narrow")

	h.addExtraKeyVal("style", `height: `+n.Attr["width"]+`em;`)

	h.setAfter("\n")

	h.AddStringTo(w)

}

func subsupOpenHTML(n *Node, w *strings.Builder) {

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

func textDirOpenHTML(n *Node, w *strings.Builder) {

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

func textDirCloseHTML(n *Node, w *strings.Builder) {

	standardCloserHTML(n, w)

}

func paginationHTML(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("pageBreak")

	h.AddStringTo(w)
	newCloseHtag("div").AddStringTo(w)
	w.WriteString("\n")

}

func subsupCloseHTML(n *Node, w *strings.Builder) {

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

func centeringOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("div")

	h.addClass("centering")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func figureOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("figure")

	h.setID(n.Attr["id"])

	h.setAfter("\n")

	h.AddStringTo(w)

}

func figureCloseHTML(n *Node, w *strings.Builder) {

	h := newCloseHtag("figure")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func bibInfoOpenHTML(n *Node, w *strings.Builder) {

	paginationHTML(n, w)

	h2 := newHtag("footer")
	h2.setAfter("\n")

	h1 := newHtag("hr")
	h1.setSelfClose()
	h1.setAfter("\n")

	h2.AddStringTo(w)
	h1.AddStringTo(w)

}

func bibInfoCloseHTML(n *Node, w *strings.Builder) {

	newCloseHtag("footer").AddStringTo(w)
	w.WriteString("\n")

}

func renderNavHTML(n *Node, w *strings.Builder) {

	Serialize(n.sectionStructure(), w, navOpenHTML, navCloseHTML)

}

func navOpenHTML(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "section":

		if n.prev == nil {
			t := newHtag("ol")
			t.setBefore("\n")
			t.setAfter("\n")
			t.AddStringTo(w)

		}

		w.WriteString(n.tocEntryHTML())

	case "top":
		t := newHtag("ol")
		t.setBefore("\n")
		t.setAfter("\n")
		t.AddStringTo(w)

		w.WriteString(n.tocEntryHTML())

	default:
		return

	}
}

func navCloseHTML(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "section":

		if n.next == nil {

			w.WriteString(newCloseHtag(`li`).String())
			w.WriteString("\n")

			t := newCloseHtag("ol")
			t.setAfter("\n")
			t.AddStringTo(w)
			return

		}

		w.WriteString(newCloseHtag(`li`).String())
		w.WriteString("\n")

	case "top":
		w.WriteString(newCloseHtag(`li`).String())
		w.WriteString("\n")
		t := newCloseHtag("ol")
		t.setAfter("\n")
		t.AddStringTo(w)

	default:
		return
	}
}

func (n *Node) tocEntryHTML() string {

	h := newHtag("li")

	h2 := newHtag("a")

	h2.addExtraKeyVal("href", "#"+n.Attr["id"])

	h2.setAfter(n.Attr["title"])

	return h.String() + h2.String() + newCloseHtag(`a`).String()

}

func gaijiNoteOpenHTML(n *Node, w *strings.Builder) {

	if n.withinNoteExclScope() {
		return
	}

	addToStringsBuilder(w, noteStartStr, n.rawString(), noteEndStr)

}

func rubyParentCloseHTML(n *Node, w *strings.Builder) {

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

func specialCharOpenHTML(n *Node, w *strings.Builder) {

	w.WriteString(n.rawString())

}

func gaijiCharOpenHTML(n *Node, w *strings.Builder) {

	w.WriteString(n.rawString())

}

func kunojiOpenHTML(n *Node, w *strings.Builder) {

	w.WriteString(n.rawString())

}

func metadataOpenHTML(n *Node, w *strings.Builder) {

	centeringOpenHTML(n, w)

	h := newHtag("div")

	h.addClass("metadata")

	if id, ok := n.Attr["data-docid"]; ok {
		h.addExtraKeyVal("data-docid", id)
	}

	h.AddStringTo(w)

}

func metadataCloseHTML(n *Node, w *strings.Builder) {

	standardCloserHTML(n, w)
	standardCloserHTML(n, w) //do twice to account for centering close
	paginationHTML(n, w)     //add page break after title+author(s)

}

func metaTitleOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("h1")

	h.addClass("title")

	h.AddStringTo(w)

}

func metaTitleCloseHTML(n *Node, w *strings.Builder) {

	h := newCloseHtag("h1")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func metaSubtitleOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("h2")

	h.addClass("subtitle")

	h.AddStringTo(w)

}

func metaSubtitleCloseHTML(n *Node, w *strings.Builder) {

	h := newCloseHtag("h2")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func metaContributorOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("h2")

	h.addClass("contributor")

	h.AddStringTo(w)

}

func metaContributorCloseHTML(n *Node, w *strings.Builder) {

	h := newCloseHtag("h2")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func captionOpenHTML(n *Node, w *strings.Builder) {

	h := newHtag("figcaption")

	h.setAfter("\n")

	h.AddStringTo(w)

}

func captionCloseHTML(n *Node, w *strings.Builder) {

	h := newCloseHtag("figcaption")

	if n.Attr["scope"] != "block" {

		h.setBefore("\n")

	}

	h.setAfter("\n")

	h.AddStringTo(w)

}

func rubyOpenHTML(n *Node, w *strings.Builder) {

	h1 := newHtag("rp")

	h1.setAfter("（")

	h2 := newCloseHtag("rp")

	h3 := newHtag("rt")

	h1.AddStringTo(w)
	h2.AddStringTo(w)
	h3.AddStringTo(w)

}

func rubyCloseHTML(n *Node, w *strings.Builder) {

	h1 := newCloseHtag("rt")

	h2 := newHtag("rp")

	h2.setAfter("）")

	h3 := newCloseHtag("rp")

	h1.AddStringTo(w)
	h2.AddStringTo(w)
	h3.AddStringTo(w)

}

func accentOpenHTML(n *Node, w *strings.Builder) {

	accentOpenTxt(n, w)

}

func unknownOpenHTML(n *Node, w *strings.Builder) {

	h1 := newHtag("span")

	h1.addClass("annotation")

	addToStringsBuilder(w, h1.String(), noteStartStr, n.Attr["raw"], noteEndStr, newCloseHtag("span").String())

}

func unknownBlockOpenHTML(n *Node, w *strings.Builder) {

	h1 := newHtag("div")

	h1.addClass("unknown")

	h1.setAfter("\n")

	addToStringsBuilder(w, h1.String(), noteStartStr, n.Attr["raw"], noteEndStr, "\n")

}

func unknownBlockCloseHTML(n *Node, w *strings.Builder) {

	h1 := newCloseHtag("div")

	h1.setAfter("\n")

	h1.setBefore("\n")

	addToStringsBuilder(w, noteStartStr, n.Attr["raw closer"], noteEndStr, h1.String())

}

func rubyGroupOpenHTML(n *Node, w *strings.Builder) {

	defer func() {
		for _, e := range linearizeDescendants(n) {

			e.Attr["ignore"] = "true"

		}
	}()

	const noLatex = false

	n.splitRuby(noLatex)

	rps := strings.Split(n.Attr["ruby base"], "\t")

	if len(rps) > 1 {

		log.Println("WARNING: line " + strconv.Itoa(n.tok.lineNumber()) + " ruby parent or ruby too long. Splitting.")
	}

	rs := strings.Split(n.Attr["ruby string"], "\t")

	h1 := newHtag("ruby")

	h1.addClass("right")

	w.WriteString(h1.String())

	for c := range len(rps) {

		w.WriteString(rps[c])

		w.WriteString(newHtag("rt").String())

		w.WriteString(rs[c])

		w.WriteString(newCloseHtag("rt").String())

	}

}

func rubyGroupCloseHTML(n *Node, w *strings.Builder) {

	w.WriteString(newCloseHtag("ruby").String())

	for _, e := range linearizeDescendants(n) {

		delete(e.Attr, "ignore")

		if e.Attr["type"] == "gaiji note" {

			addToStringsBuilder(w, noteStartStr, e.rawString(), noteEndStr)

		}
	}

}
