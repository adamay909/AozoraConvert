package aozoratext

import (
	"strings"
)

func renderAozoraText(n *Node) string {

	return Serialize(n, azrTxtFormatterOpen, azrTxtFormatterClose)

}

func azrTxtFormatterOpen(n *Node) string {

	switch n.Attr["type"] {

	case "text":
		return n.Attr["raw"]

	case "paragraph":
		return ""

	case "ruby parent":
		return rubyBaseOpenTxt(n)

	case "ruby":
		return rubyStartStr

	case "empty line":
		return lineBreakStr

	case "section":
		return ""

	case "indentation":
		return indentationOpenTxt(n)

	case "bottom align":
		return bottomAlignOpenTxt(n)

	case "emphasis":
		return noteStringOpenTxt(n)

	case "line decoration":
		return noteStringOpenTxt(n)

	case "narrow paragraph":
		return noteStringOpenTxt(n)

	case "inline section":
		return noteStringOpenTxt(n)

	case "window section":
		return noteStringOpenTxt(n)

	case "inline note":
		return noteStringOpenTxt(n)

	case "rubylike note":
		return noteStringOpenTxt(n)

	case "caption":
		return captionOpenTxt(n)

	case "font shape":
		return noteStringOpenTxt(n)

	case "font size":
		return noteStringOpenTxt(n)

	case "offset":
		return noteStringOpenTxt(n)

	case "text direction":
		return noteStringOpenTxt(n)

	case "section title":
		return noteStringOpenTxt(n)

	case "figure":
		return ""

	case "image":
		return noteStringOpenTxt(n) + "\n"

	case "pagination":
		return noteStringOpenTxt(n) + "\n"

	case "kunten":
		return noteStringOpenTxt(n)

	case "okurigana":

		return noteStringOpenTxt(n)

	case "centering":
		return noteStringOpenTxt(n)

	case "bibliographical info":
		return bibInfoOpenTxt(n)

	case "gaiji char":
		return gaijiCharOpenTxt(n)

	case "kunoji":
		return kunojiOpenTxt(n)

	case "gaiji note":
		return gaijiNoteOpenTxt(n)

	case "warichu line break":
		return noteStringOpenTxt(n)

	case "metadata":
		return ""

	case "meta title":
		return ""

	case "meta subtitle":
		return ""

	case "meta contributor":
		return ""

	case "special char":
		return specialCharOpenTxt(n)

	case "accent string":
		return accentOpenTxt(n)

	case "document":
		return ""

	default:
		return noteStartStr + n.Attr["raw"] + noteEndStr

	}

}

func azrTxtFormatterClose(n *Node) string {

	switch n.Attr["type"] {

	case "text":
		return ""

	case "paragraph":
		return "\n"

	case "ruby parent":
		return ""

	case "ruby":
		return rubyEndStr

	case "empty line":
		return ""

	case "section":
		return ""

	case "indentation":
		return indentationCloseTxt(n)

	case "bottom align":
		return bottomAlignCloseTxt(n)

	case "emphasis":
		return noteStringCloseTxt(n)

	case "line decoration":
		return noteStringCloseTxt(n)

	case "narrow paragraph":
		return noteStringCloseTxt(n)

	case "inline section":
		return noteStringCloseTxt(n)

	case "window section":
		return noteStringCloseTxt(n)

	case "inline note":
		return noteStringCloseTxt(n)

	case "rubylike note":
		return noteStringCloseTxt(n)

	case "caption":
		return captionCloseTxt(n)

	case "font shape":
		return noteStringCloseTxt(n)

	case "font size":
		return noteStringCloseTxt(n)

	case "offset":
		return noteStringCloseTxt(n)

	case "text direction":
		return noteStringCloseTxt(n)

	case "section title":
		return sectionTitleCloseTxt(n)

	case "figure":
		return ""

	case "image":
		return ""

	case "pagination":
		return ""

	case "kunten":
		return ""

	case "okurigana":
		return ""

	case "centering":
		return ""
	case "bibliographical info":
		return ""

	case "gaiji note":

		return ""

	case "metadata":
		return metadataCloseTxt(n)

	case "meta title":
		return "\n"

	case "meta subtitle":
		return "\n"

	case "meta contributor":
		return "\n"

	case "document":
		return ""

	default:
		return ""
	}

}

func noteStringOpenTxt(n *Node) string {

	output.Reset()

	addToStringsBuilder(output, noteStartStr, n.Attr["raw"], noteEndStr)

	if n.isBlockFormat() {

		output.WriteString("\n")

	}

	return output.String()

}

func noteStringCloseTxt(n *Node) string {

	output.Reset()

	if n.Attr["raw closer"] == "" {
		addToStringsBuilder(output, noteStartStr, "missing closer:", n.Attr["type"], " ", n.Attr["raw"], noteEndStr)
	} else {

		addToStringsBuilder(output, noteStartStr, n.Attr["raw closer"], noteEndStr)
	}

	if n.isJiage() {

		return ""

	}

	if n.isJisage() {

		if n.next != nil {

			if n.next.isJisage() {

				return ""
			}
		}

	}

	if n.isBlockFormat() {

		output.WriteString("\n")

	}

	return output.String()

}

func indentationOpenTxt(n *Node) string {

	return noteStringOpenTxt(n)

	output.Reset()

	if n.firstChild.Attr["type"] == "section title" {

		addToStringsBuilder(output, noteStartStr, strings.TrimPrefix(n.Attr["raw"], blockStartStr), noteEndStr)

		return output.String()

	}

	if n.innerParagraphCount() == 1 {

		addToStringsBuilder(output, noteStartStr, strings.TrimPrefix(n.Attr["raw"], blockStartStr), noteEndStr)

		return output.String()

	}

	n.SetAttr("open", "true")

	return noteStringOpenTxt(n)

}

func indentationCloseTxt(n *Node) string {

	return noteStringCloseTxt(n)

	if n.next != nil {

		if n.next.Attr["type"] == "indentation" {

			return ""

		}
	}

	if n.firstChild.Attr["type"] == "section title" {

		return ""

	}

	if n.innerParagraphCount() == 1 {

		for e := n; e.Attr["type"] == "indentation"; e = e.prev {

			if e.Attr["open"] == "true" {

				delete(e.Attr, "open")

				return noteStringCloseTxt(n)

			}

			if e == nil {
				break
			}
		}

		return ""

	}

	return noteStringCloseTxt(n)

}

func bottomAlignOpenTxt(n *Node) string {

	return noteStringOpenTxt(n)

	output.Reset()

	if n.firstChild.Attr["scope"] == "block" {

		return noteStringOpenTxt(n)

	}

	if n.innerParagraphCount() > 1 {

		return noteStringOpenTxt(n)

	}

	addToStringsBuilder(output, noteStartStr, strings.TrimPrefix(n.Attr["raw"], blockStartStr), noteEndStr)

	return output.String()

}

func bottomAlignCloseTxt(n *Node) string {

	return noteStringCloseTxt(n)

	if n.firstChild.Attr["scope"] == "block" {

		return noteStringCloseTxt(n)

	}

	return noteStringCloseTxt(n)

	if n.innerParagraphCount() > 1 {

		return noteStringCloseTxt(n)

	}

	return ""
}

func sectionTitleCloseTxt(n *Node) string {

	output := new(strings.Builder)

	addToStringsBuilder(output, noteStartStr, n.Attr["raw closer"], noteEndStr, "\n")

	return output.String()

}

func bibinfostringTxt(n *Node) string {

	if n.Attr["raw"] == "底本：" {
		return n.Attr["raw"]
	}

	return noteStartStr + n.Attr["raw"] + noteEndStr + lineBreakStr

}

func captionOpenTxt(n *Node) string {

	//	return "［＃ここからキャプション］\n"

	if n.innerParagraphCount() > 1 {

		return "［＃ここからキャプション］\n"

	}

	return "［＃キャプション］"

}

func captionCloseTxt(n *Node) string {

	//	return "［＃ここでキャプション終わり］\n"

	if n.innerParagraphCount() > 1 {

		return "［＃ここでキャプション終わり］\n"

	}

	return "［＃キャプション終わり］\n"
}

func gaijiNoteOpenTxt(n *Node) string {

	return noteStartStr + strings.TrimPrefix(n.Attr["raw"], referenceMarkStr+"は") + noteEndStr

}

func rubyBaseOpenTxt(n *Node) string {

	return rubyParentStartStr

}

func metadataCloseTxt(n *Node) string {

	return "\n"

}

func specialCharOpenTxt(n *Node) string {

	return n.Attr["raw"]

}

func gaijiCharOpenTxt(n *Node) string {

	if o_jis0208 {

		return n.Attr["raw"]

	}

	return n.Attr["alt raw"]

}

func accentOpenTxt(n *Node) string {

	if o_jis0208 {

		return accentStartStr + n.Attr["raw"] + accentEndStr
	}

	return n.Attr["alt raw"]

}

func kunojiOpenTxt(n *Node) string {

	if o_jis0208 {

		return n.Attr["raw"]
	}

	return n.Attr["alt raw"]

}

func bibInfoOpenTxt(n *Node) string {

	if n.Attr["raw"] != "" {

		return noteStartStr + "本文終わり" + noteEndStr + "\n"
	}

	return ""

}
