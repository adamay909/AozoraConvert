package aozoratext

import (
	"fmt"
	"log"
	"strings"
)

func CheckStructure(text string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("The document has errors.\n", r)
			fmt.Println("Exiting")
			return
		}
	}()

	offset := aztextMainStart(text)

	fmt.Println("offset", offset)

	t := tokenizeAndFix(text, offset)

	var openers []*token

	open := 0

	for e := t.firstToken(); e != nil; e = e.next {

		if e.isEncloseOpen() {

			openers = append(openers, e)

			open++

			continue

		}

		if e.isEncloseClose() {

			open--

			if open < 0 {

				message := e.info() + " does not close any opener."

				panic(message)

			}

			if !matchedPair(openers[len(openers)-1], e) {

				message := "closing tag " + e.info() + "appeared before closing of " + openers[len(openers)-1].info()

				panic(message)

			}

			openers = openers[:len(openers)-1]

		}

	}

	if len(openers) > 0 {

		message := "the following tags are without closers:\n"

		for _, e := range openers {

			message = message + e.info() + "\n"

		}

		panic(message)

	}

	log.Println("OK.")

}

func (t *token) isEncloseOpen() bool {

	if t.isPairMarkerOpen() {
		return true
	}

	switch t.tokType {

	case rubyParentStartToken, rubyStartToken, paragraphToken, sectionTitleStartToken, sectionToken, subsectionToken, subsubsectionToken, figureStartToken:
		return true

	}

	return false

}

func (t *token) isEncloseClose() bool {

	if t.isPairMarkerClose() {
		return true
	}

	switch t.tokType {

	case rubyParentEndToken, rubyEndToken, lineBreakToken, sectionTitleEndToken, sectionEndToken, subsectionEndToken, subsubsectionEndToken, figureEndToken:
		return true
	}

	return false

}

func matchedPair(e1, e2 *token) bool {

	switch e1.tokType {

	case rubyParentStartToken:
		return e2.tokType == rubyParentEndToken

	case rubyStartToken:
		return e2.tokType == rubyEndToken

	case paragraphToken:
		return e2.tokType == lineBreakToken

	case sectionTitleStartToken:
		return e2.tokType == sectionTitleEndToken

	case sectionToken:
		return e2.tokType == sectionEndToken

	case subsectionToken:
		return e2.tokType == subsectionEndToken

	case subsubsectionToken:
		return e2.tokType == subsubsectionEndToken

	case figureStartToken:
		return e2.tokType == figureEndToken

	}

	for _, m := range pairMarker {

		if strings.HasSuffix(e1.innerString(), m) {

			return strings.HasSuffix(strings.TrimSuffix(e2.innerString(), formatEndStr), m)

		}
	}

	return false
}
