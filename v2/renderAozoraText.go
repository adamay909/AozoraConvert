package aozoraconvert

import (
	"strings"
)

func renderAozoraText(n *Node, w *strings.Builder) {

	Serialize(n, w, azrTxtFormatterOpen, azrTxtFormatterClose)

}

func renderInnerAozoraText(n *Node, w *strings.Builder) {

	SerializeDescendants(n, w, azrTxtFormatterOpen, azrTxtFormatterClose)

}

func renderInnerTextOnly(n *Node, w *strings.Builder) {

	SerializeDescendants(n, w, plaintextWriterOpen, plaintextWriterClose)

}

func plaintextWriterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "text", "special char", "kunoji", "gaiji char":
		azrTxtFormatterOpen(n, w)

	default:
		return
	}
}

func plaintextWriterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "text", "special char", "kunoji", "gaiji char":
		azrTxtFormatterClose(n, w)

	default:
		return
	}
}

func azrTxtFormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "text":
		w.WriteString(n.rawString())

	case "paragraph":
		return

	case "kanbun":
		return

	case "ruby group":
		return

	case "ruby parent":
		w.WriteString(rubyParentStartStr)

	case "ruby":
		w.WriteString(rubyStartStr)

	case "empty line":
		w.WriteString(lineBreakStr)

	case "section":
		return

	case "indentation":
		indentationOpenTxt(n, w)

	case "bottom align":
		bottomAlignOpenTxt(n, w)

	case "emphasis":
		noteStringOpenTxt(n, w)

	case "line decoration":
		noteStringOpenTxt(n, w)

	case "narrow paragraph":
		noteStringOpenTxt(n, w)

	case "inline section":
		noteStringOpenTxt(n, w)

	case "window section":
		noteStringOpenTxt(n, w)

	case "inline note":
		noteStringOpenTxt(n, w)

	case "rubylike note":
		noteStringOpenTxt(n, w)

	case "caption":
		captionOpenTxt(n, w)

	case "font shape":
		noteStringOpenTxt(n, w)

	case "font size":
		noteStringOpenTxt(n, w)

	case "offset":
		noteStringOpenTxt(n, w)

	case "text direction":
		noteStringOpenTxt(n, w)

	case "section title":
		noteStringOpenTxt(n, w)

	case "figure":
		return

	case "image":
		noteStringOpenTxt(n, w)

	case "pagination":
		noteStringOpenTxt(n, w)
		w.WriteString("\n")

	case "accent start":
		accentStartTxt(n, w)

	case "accent end":
		accentEndTxt(n, w)

	case "kunten":
		noteStringOpenTxt(n, w)

	case "okurigana":
		noteStringOpenTxt(n, w)

	case "centering":
		noteStringOpenTxt(n, w)

	case "bibliographical info":
		bibInfoOpenTxt(n, w)

	case "gaiji char":
		gaijiCharOpenTxt(n, w)

	case "kunoji":
		kunojiOpenTxt(n, w)

	case "gaiji note":
		gaijiNoteOpenTxt(n, w)

	case "warichu line break":
		noteStringOpenTxt(n, w)

	case "metadata":
		return

	case "meta title":
		return

	case "meta subtitle":
		return

	case "meta contributor":
		return

	case "special char":
		specialCharOpenTxt(n, w)

	case "main text":
		return

	case "document":
		return

	case "end marker":
		endMarkerOpenTxt(n, w)

	case "unknown":
		switch {
		case oJis0208:
			addToStringsBuilder(w, noteStartStr, n.Attr["raw"], noteEndStr)

		case oJis0213:
			addToStringsBuilder(w, noteStartStr, n.Attr["jis0213 raw"], noteEndStr)

		default:

			addToStringsBuilder(w, noteStartStr, n.Attr["unicode raw"], noteEndStr)
		}

		if n.Attr["force linebreak"] == "true" {
			w.WriteString("\n")
		}

		return

	case "unknown block type":
		noteStringOpenTxt(n, w)

	default:
		msglog.Println("Renderer: unknown node type: " + n.String())

		if oJis0208 {
			addToStringsBuilder(w, noteStartStr, n.rawString(), noteEndStr)
			return
		}
		addToStringsBuilder(w, noteStartStr, n.Attr["unicode raw"], noteEndStr)
		return
	}

}

func azrTxtFormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "text":
		return

	case "paragraph":
		w.WriteString("\n")

	case "kanbun":
		w.WriteString("\n")

	case "ruby group":
		return

	case "ruby parent":
		return

	case "ruby":
		w.WriteString(rubyEndStr)

	case "empty line":
		return

	case "section":
		return

	case "indentation":
		indentationCloseTxt(n, w)

	case "bottom align":
		bottomAlignCloseTxt(n, w)

	case "emphasis":
		noteStringCloseTxt(n, w)

	case "line decoration":
		noteStringCloseTxt(n, w)

	case "narrow paragraph":
		noteStringCloseTxt(n, w)

	case "inline section":
		noteStringCloseTxt(n, w)

	case "window section":
		noteStringCloseTxt(n, w)

	case "inline note":
		noteStringCloseTxt(n, w)

	case "rubylike note":
		noteStringCloseTxt(n, w)

	case "caption":
		captionCloseTxt(n, w)

	case "font shape":
		noteStringCloseTxt(n, w)

	case "font size":
		noteStringCloseTxt(n, w)

	case "offset":
		noteStringCloseTxt(n, w)

	case "text direction":
		noteStringCloseTxt(n, w)

	case "section title":
		sectionTitleCloseTxt(n, w)

	case "figure":
		return

	case "image":
		if n.Attr["style"] == "inline" {
			return
		}
		w.WriteString("\n")

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
		return

	case "bibliographical info":
		return

	case "gaiji char":
		return

	case "kunoji":
		return

	case "gaiji note":
		return

	case "warichu line break":
		return

	case "metadata":
		metadataCloseTxt(n, w)

	case "meta title":
		w.WriteString("\n")

	case "meta subtitle":
		w.WriteString("\n")

	case "meta contributor":
		w.WriteString("\n")

	case "special char":
		return

	case "main text":
		return

	case "end marker":
		return

	case "document":
		return

	case "unknown":
		return

	case "unknown block type":
		noteStringCloseTxt(n, w)

	default:
		return
	}

}

func noteStringOpenTxt(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, noteStartStr, n.rawString(), noteEndStr)

	if n.isBlockFormat() {

		w.WriteString("\n")

	}

}

func noteStringCloseTxt(n *Node, w *strings.Builder) {

	if oJis0208 {

		addToStringsBuilder(w, noteStartStr, n.Attr["raw closer"], noteEndStr)

	} else {
		if n.Attr["raw unicode closer"] != "" {

			addToStringsBuilder(w, noteStartStr, n.Attr["raw unicode closer"], noteEndStr)

		} else {

			addToStringsBuilder(w, noteStartStr, n.Attr["raw closer"], noteEndStr)
		}
	}

	if n.isBlockFormat() {

		w.WriteString("\n")
		return
	}

	if n.next != nil && n.next.Attr["type"] == "empty line" {
		w.WriteString("\n")
	}

}

func indentationOpenTxt(n *Node, w *strings.Builder) {

	if n.isJisage() && (n.prev == nil || !n.prev.isJisage()) && (n.next == nil || !n.next.isJisage()) {

		b := new(strings.Builder)

		renderInnerTextOnly(n, b)

		if n.innerParagraphCount() < 2 && len([]rune(b.String())) < 80 {

			addToStringsBuilder(w, noteStartStr, strings.TrimPrefix(n.rawString(), blockStartStr), noteEndStr)

			n.SetAttr("single line jisage", "true")

			return
		}
	}

	if n.isJisage() && n.firstChild.Attr["type"] == "section title" && n.firstChild.next == nil {

		addToStringsBuilder(w, noteStartStr, strings.TrimPrefix(n.rawString(), blockStartStr), noteEndStr)

		n.SetAttr("single line jisage", "true")

		return
	}

	noteStringOpenTxt(n, w)
}

func indentationCloseTxt(n *Node, w *strings.Builder) {

	if n.next.isJisage() {
		return
	}

	if n.Attr["single line jisage"] == "true" {
		return
	}

	noteStringCloseTxt(n, w)

}

func bottomAlignOpenTxt(n *Node, w *strings.Builder) {

	noteStringOpenTxt(n, w)

}

func bottomAlignCloseTxt(n *Node, w *strings.Builder) {

	if n.Attr["scope"] == "block" {

		w2 := new(strings.Builder)

		renderInnerAozoraText(n, w2)

		if !strings.HasSuffix(w2.String(), "\n") {
			w.WriteString("\n")
		}

		noteStringCloseTxt(n, w)
		return
	}

	if ok, _ := n.withinScopeOfType("paragraph"); !ok {
		w.WriteString("\n")
	}

	return

}

func sectionTitleCloseTxt(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, noteStartStr, n.Attr["raw closer"], noteEndStr)

	w.WriteString("\n")

}

func bibinfostringTxt(n *Node) string {

	if n.rawString() == "底本：" {
		return n.rawString()
	}

	return noteStartStr + n.rawString() + noteEndStr + lineBreakStr

}

func captionOpenTxt(n *Node, w *strings.Builder) {

	if n.isBlockFormat() || n.innerParagraphCount() > 1 {

		w.WriteString("［＃ここからキャプション］\n")
		return
	}

	w.WriteString("［＃キャプション］")

}

func captionCloseTxt(n *Node, w *strings.Builder) {

	if n.isBlockFormat() || n.innerParagraphCount() > 1 {

		w.WriteString("［＃ここでキャプション終わり］\n")
		return
	}

	w.WriteString("［＃キャプション終わり］")

	if ok, _ := n.withinScopeOfType("paragraph"); ok {
		return
	}

	w.WriteString("\n")
}

func gaijiNoteOpenTxt(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, noteStartStr, strings.TrimPrefix(n.rawString(), referenceMarkStr+"は"), noteEndStr)

}

func rubyBaseOpenTxt(n *Node) string {

	return rubyParentStartStr

}

func metadataCloseTxt(n *Node, w *strings.Builder) {

	w.WriteString("\n")

}

func specialCharOpenTxt(n *Node, w *strings.Builder) {

	if !oParsable {
		w.WriteString(n.rawString())
		return
	}

	w.WriteString(n.Attr["escaped raw"])

}

func gaijiCharOpenTxt(n *Node, w *strings.Builder) {

	w.WriteString(n.rawString())

}

func accentOpenTxt(n *Node, w *strings.Builder) {

	if oJis0208 {

		addToStringsBuilder(w, accentStartStr)

		return
	}
	return
}

func accentCloseTxt(n *Node, w *strings.Builder) {

	if oJis0208 {

		addToStringsBuilder(w, accentEndStr)

		return
	}
	return
}

func kunojiOpenTxt(n *Node, w *strings.Builder) {

	w.WriteString(n.rawString())

}

func bibInfoOpenTxt(n *Node, w *strings.Builder) {

	if n.rawString() != "" {

		addToStringsBuilder(w, noteStartStr, "本文終わり", noteEndStr, "\n")
	}

	return

}

func rubyGroupOpenTxt(n *Node, w *strings.Builder) {

	const noLatex = false

	n.splitRuby(noLatex)

	w.WriteString(rubyParentStartStr)

	w.WriteString(strings.Join(strings.Split(n.Attr["ruby base"], "\t"), ""))

	w.WriteString(rubyStartStr)

	w.WriteString(strings.Join(strings.Split(n.Attr["ruby string"], "\t"), ""))

	w.WriteString(rubyEndStr)

	for _, e := range linearizeDescendants(n) {
		e.Attr["ignore"] = "true"
	}

}

func rubyGroupCloseTxt(n *Node, w *strings.Builder) {

	for _, e := range linearizeDescendants(n) {
		delete(e.Attr, "ignore")
	}
}

func accentStartTxt(n *Node, w *strings.Builder) {

	if oJis0208 {

		w.WriteString(accentStartStr)

	}

	return

}

func accentEndTxt(n *Node, w *strings.Builder) {

	if oJis0208 {

		w.WriteString(accentEndStr)

	}

	return

}

func endMarkerOpenTxt(n *Node, w *strings.Builder) {

	addToStringsBuilder(w, n.Attr["raw"], "\n")

}
