package aozoratext

import (
	"errors"
	"strconv"
	"strings"
)

func isValidStructure(prevNode *Node, tok *token) (err error) {

	msgBuilder := new(strings.Builder)

	switch {

	case prevNode == nil:

		addToStringsBuilder(msgBuilder, "\nline ", strconv.Itoa(tok.lineNumber()), ": ", "ERROR: attempting to close node without opener: ", tok.info())

		return errors.New(msgBuilder.String())

	case prevNode.tok == nil:

		addToStringsBuilder(msgBuilder, "\nline ", strconv.Itoa(tok.lineNumber()), ": ", "ERROR: attempting to close node without opener: ", tok.info())

		return errors.New(msgBuilder.String())

	case !matched(tok, prevNode.tok):

		addToStringsBuilder(msgBuilder, "\nERROR: mismatched annotation start and end:\n")

		addToStringsBuilder(msgBuilder, prevNode.tok.info(), "\n closed by: \n")

		addToStringsBuilder(msgBuilder, tok.info())

		return errors.New(msgBuilder.String())

	default:

		return err

	}
}

func matched(closer, opener *token) bool {

	if closer == nil || opener == nil {
		return false
	}

	switch closer.tokType {

	case paragraphEndToken:
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

	case bibInfoEndToken:
		return opener.tokType == bibInfoToken

	case mainTextEndToken:
		return opener.tokType == mainTextStartToken
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
