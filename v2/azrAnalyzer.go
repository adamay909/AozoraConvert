package aozoratext

import (
	"strings"
)

type atomType int

const (
	emptyAtom atomType = iota
	startAtom
	markupNoteStart
	markupNoteEnd
	markupNoteDelimiter
	rubyStartTag
	rubyEndTag
	rubyBaseStartTag
	gaijiMarker
	noteStartTag
	noteEndTag
	bibInfoTag
	text
	endOfLine
	emptyLine
	accentStartTag
	accentEndTag
	kunojiTag
	kunojiDakuTag
)

var (
	markupNoteStartStr = `-------------------------------------------------------
【テキスト中に現れる記号について】`

	markupNoteEndStr = `-------------------------------------------------------`

	markupNoteDelimiterStr = `-------------------------------------------------------`

	rubyStartStr = `《`

	rubyEndStr = `》`

	rubyBaseStartStr = `｜`

	gaijiMarkerStr = `※［＃`

	noteStartStr = `［＃`

	noteEndStr = `］`

	bibInfoStartStr = `底本：`

	emptyStr = ""

	lineBreakStr = "\n"

	refStartStr = "「"

	refEndStr = "」"

	kunojiStr = "／＼"

	kunojiStrU = "〳〵"

	kunojiDakuStr = "／″＼"

	kunojiDakuStrU = "〴〵"

	blockStartStr = "ここから"

	blockEndStr = "ここで"

	formatEndStr = "終わり"

	referenceMarkStr = "※"

	accentStartStr = "〔"

	accentEndStr = "〕"
)

func typeOf(s string) atomType {

	switch {

	case len(s) == 0:
		return emptyLine

	case strings.HasPrefix(s, rubyStartStr):
		return rubyStartTag

	case strings.HasPrefix(s, rubyEndStr):
		return rubyEndTag

	case strings.HasPrefix(s, rubyBaseStartStr):
		return rubyBaseStartTag

	case strings.HasPrefix(s, gaijiMarkerStr):
		return gaijiMarker

	case strings.HasPrefix(s, noteStartStr):
		return noteStartTag

	case strings.HasPrefix(s, noteEndStr):
		return noteEndTag

	case strings.HasPrefix(s, bibInfoStartStr):
		return bibInfoTag

	case strings.HasPrefix(s, accentStartStr):
		return accentStartTag

	case strings.HasPrefix(s, accentEndStr):
		return accentEndTag

	case strings.HasPrefix(s, kunojiStr):
		return kunojiTag

	case strings.HasPrefix(s, kunojiDakuStr):
		return kunojiDakuTag

	default:
		return text

	}
}

func findContiguousText(s string) (i int) {

	for i = range s {

		tp := typeOf(s[i:])

		//		if tp == bibInfoTag {
		//			continue
		//		}

		if tp == noteEndTag {
			tokenizerLog.Println("Unexpected note end string. Treating as text.")
			continue

		}

		if tp != text {
			return i
		}
	}

	return len(s)

}

func findMatchingCloser(o atomType, s string) (i int) {

	counter := 0

	i = 0

	for i = range s {

		if typeOf(s[i:]) == o {
			counter++
			continue
		}

		if o == noteStartTag {

			if typeOf(s[i:]) == gaijiMarker {
				counter++
				continue
			}
		}

		if typeOf(s[i:]) != closingTagOf(o) {

			continue
		}

		counter--

		if counter == 0 {
			break
		}
	}

	if o == rubyBaseStartTag {
		return i
	}

	if counter != 0 {

		tokenizerLog.Println("Found unclosed opening tag: " + o.rawString() + ". Assuming scope ends at end of line.")

		return len(s)
	}

	return i + len(closingTagOf(o).rawString())
}

func (o atomType) rawString() string {

	switch o {

	case rubyStartTag:
		return rubyStartStr

	case rubyEndTag:
		return rubyEndStr

	case rubyBaseStartTag:
		return rubyBaseStartStr

	case gaijiMarker:
		return gaijiMarkerStr

	case noteStartTag:
		return noteStartStr

	case noteEndTag:
		return noteEndStr

	case bibInfoTag:
		return bibInfoStartStr

	case accentStartTag:
		return accentStartStr

	case accentEndTag:
		return accentEndStr
	}

	return ""

}
