package aozoratext

import (
	"strings"
)

func renderAozoraText(n *node) string {

	return stringify(n, azrTxtFormatterOpen, azrTxtFormatterClose)

}

func azrTxtFormatterOpen(n *node) string {

	switch n.attr["type"] {

	case "text":
		return n.attr["raw"]

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
		return noteStartStr + n.attr["raw"] + noteEndStr

	}

}

func azrTxtFormatterClose(n *node) string {

	switch n.attr["type"] {

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

func noteStringOpenTxt(n *node) string {

	output.Reset()

	addToStringsBuilder(output, noteStartStr, n.attr["raw"], noteEndStr)

	if n.isBlockFormat() {

		output.WriteString("\n")

	}

	return output.String()

}

func noteStringCloseTxt(n *node) string {

	output.Reset()

	if n.attr["raw closer"] == "" {
		addToStringsBuilder(output, noteStartStr, "missing closer:", n.attr["type"], " ", n.attr["raw"], noteEndStr)
	} else {

		addToStringsBuilder(output, noteStartStr, n.attr["raw closer"], noteEndStr)
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

func indentationOpenTxt(n *node) string {

	return noteStringOpenTxt(n)

	output.Reset()

	if n.firstChild.attr["type"] == "section title" {

		addToStringsBuilder(output, noteStartStr, strings.TrimPrefix(n.attr["raw"], blockStartStr), noteEndStr)

		return output.String()

	}

	if n.innerParagraphCount() == 1 {

		addToStringsBuilder(output, noteStartStr, strings.TrimPrefix(n.attr["raw"], blockStartStr), noteEndStr)

		return output.String()

	}

	n.setAttr("open", "true")

	return noteStringOpenTxt(n)

}

func indentationCloseTxt(n *node) string {

	return noteStringCloseTxt(n)

	if n.next != nil {

		if n.next.attr["type"] == "indentation" {

			return ""

		}
	}

	if n.firstChild.attr["type"] == "section title" {

		return ""

	}

	if n.innerParagraphCount() == 1 {

		for e := n; e.attr["type"] == "indentation"; e = e.prev {

			if e.attr["open"] == "true" {

				delete(e.attr, "open")

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

func bottomAlignOpenTxt(n *node) string {

	return noteStringOpenTxt(n)

	output.Reset()

	if n.firstChild.attr["scope"] == "block" {

		return noteStringOpenTxt(n)

	}

	if n.innerParagraphCount() > 1 {

		return noteStringOpenTxt(n)

	}

	addToStringsBuilder(output, noteStartStr, strings.TrimPrefix(n.attr["raw"], blockStartStr), noteEndStr)

	return output.String()

}

func bottomAlignCloseTxt(n *node) string {

	return noteStringCloseTxt(n)

	if n.firstChild.attr["scope"] == "block" {

		return noteStringCloseTxt(n)

	}

	return noteStringCloseTxt(n)

	if n.innerParagraphCount() > 1 {

		return noteStringCloseTxt(n)

	}

	return ""
}

func sectionTitleCloseTxt(n *node) string {

	output := new(strings.Builder)

	addToStringsBuilder(output, noteStartStr, n.attr["raw closer"], noteEndStr, "\n")

	return output.String()

}

func bibinfostringTxt(n *node) string {

	if n.attr["raw"] == "底本：" {
		return n.attr["raw"]
	}

	return noteStartStr + n.attr["raw"] + noteEndStr + lineBreakStr

}

func captionOpenTxt(n *node) string {

	//	return "［＃ここからキャプション］\n"

	if n.innerParagraphCount() > 1 {

		return "［＃ここからキャプション］\n"

	}

	return "［＃キャプション］"

}

func captionCloseTxt(n *node) string {

	//	return "［＃ここでキャプション終わり］\n"

	if n.innerParagraphCount() > 1 {

		return "［＃ここでキャプション終わり］\n"

	}

	return "［＃キャプション終わり］\n"
}

func gaijiNoteOpenTxt(n *node) string {

	return noteStartStr + strings.TrimPrefix(n.attr["raw"], referenceMarkStr+"は") + noteEndStr

}

func rubyBaseOpenTxt(n *node) string {

	return rubyBaseStartStr

}

func metadataCloseTxt(n *node) string {

	return "\n"

}

func specialCharOpenTxt(n *node) string {

	return n.attr["raw"]

}

func gaijiCharOpenTxt(n *node) string {

	if o_jis0208 {

		return n.attr["raw"]

	}

	return n.attr["alt raw"]

}

func accentOpenTxt(n *node) string {

	if o_jis0208 {

		return accentStartStr + n.attr["raw"] + accentEndStr
	}

	return n.attr["alt raw"]

}

func kunojiOpenTxt(n *node) string {

	if o_jis0208 {

		return n.attr["raw"]
	}

	return n.attr["alt raw"]

}

func bibInfoOpenTxt(n *node) string {

	if n.attr["raw"] != "" {

		return noteStartStr + "本文終わり" + noteEndStr + "\n"
	}

	return ""

}
