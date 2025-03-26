package aozoratext

import (
	"strconv"
	"strings"
)

func validStructure(prevNode *Node, tok *token) (bool, string) {

	msgBuilder := new(strings.Builder)

	switch {

	case prevNode == nil:

		val := false

		addToStringsBuilder(msgBuilder, "line ", strconv.Itoa(tok.lineNumber()), ": ", "attempting to close node without opener: ", tok.info())

		//textContext())

		panic(msgBuilder.String())

		return val, msgBuilder.String()

	case !matched(tok, prevNode.tok):

		val := false

		addToStringsBuilder(msgBuilder, "mismatched annotation start and end:\n")

		addToStringsBuilder(msgBuilder, prevNode.tok.info(), "\n closed by: \n")

		addToStringsBuilder(msgBuilder, tok.info())

		panic(msgBuilder.String())
		return val, msgBuilder.String()

	default:

		return true, ""

	}
}

func matched(closer, opener *token) bool {

	switch closer.tokType {

	case lineBreakToken:
		return opener.tokType == paragraphToken

	case rubyEndToken:
		return opener.tokType == rubyStartToken

	case rubyParentEndToken:
		return opener.tokType == rubyParentStartToken

	case sectionEndToken:
		return opener.tokType == sectionToken

	case subsectionEndToken:
		return opener.tokType == subsectionToken

	case subsubsectionEndToken:
		return opener.tokType == subsubsectionToken

	case figureEndToken:
		return opener.tokType == figureStartToken

	case centeringEndToken:
		return opener.isFormatOfType(centeringMarker)
	}

	if opener.isBlockStartNote() {

		if !closer.isBlockClosingNote() {

			return false

		}
	}

	if closer.isBlockClosingNote() {

		if !opener.isBlockStartNote() {

			return false

		}
	}

	switch {

	case opener.isFormatOfType(bottomalignMarker):
		if !opener.isBlockStartNote() {
			return closer.tokType == alignBottomCloserToken
		}
		fallthrough

	default:
		for _, m := range pairMarker {
			if strings.HasSuffix(opener.innerString(), m) {
				return strings.HasSuffix(strings.TrimPrefix(strings.TrimSuffix(closer.innerString(), formatEndStr), blockEndStr), m)
			}
		}

	}

	return false
}
