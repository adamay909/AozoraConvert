package aozoratext

import (
	"log"
	"os"
	"strings"
)

type opt int

var o_full, o_raw, o_jis0208, o_debug bool

type tokenizer struct {
	data        string
	position    int
	lineCounter int
}

var tokenizerOption opt

var tokenizerLog *log.Logger

var dummyLog *strings.Builder

func init() {

	o_full = true

	o_raw = false

	o_jis0208 = false

	o_debug = false

	tokenizerLog = log.New(os.Stdout, "", 0)

	log.SetFlags(0)

}

func setOutputOption(o string) {

	switch o {
	case "jis0208":
		o_full = false
		o_jis0208 = true
		o_raw = false

	case "raw":
		o_full = false
		o_jis0208 = true
		o_raw = true

	default:
		o_full = true
		o_jis0208 = false
		o_raw = false
	}
	return
}

func setDebug() {

	o_debug = true

}

func unsetDebug() {

	o_debug = false

}

func tokenize(s string, offset int) *token {

	tknz := newTokenizer(s, offset)

	t0 := tknz.nextToken()

	for t1 := t0; !tknz.empty(); t1 = t1.next {

		t1.addTokenRight(tknz.nextToken())

	}

	return t0

}

func tokenizeAll(text string, lineOffset int) *token {

	tokenString := tokenize(strings.Join(strings.Split(text, "\n")[lineOffset:], "\n"), lineOffset)

	tokenString.lastToken().insertTokenRight(newTokenOfType(endOfLineToken))

	tokenString.fixLines()

	tokenString.fixEmptyText()

	if !o_raw {

		tokenString.fixDocument()

		tokenString.insertSectionEnds()

	}

	return tokenString

}

func newTokenizer(d string, offset int) *tokenizer {

	r := new(tokenizer)

	r.data = d

	r.position = 0

	r.lineCounter = offset

	return r
}

func (tknz *tokenizer) empty() bool {

	return tknz.position == len(tknz.data)

}

func (tknz *tokenizer) nextToken() (e *token) {

	e = newToken()

	end := 0

	switch typeOf(tknz.data[tknz.position:]) {

	case text:

		e.tokType = textToken
		end = findContiguousText(tknz.data[tknz.position:])

	case rubyStartTag:

		e.tokType = rubyStartToken
		end = len(rubyStartStr)

	case rubyEndTag:
		e.tokType = rubyEndToken
		end = len(rubyEndStr)

	case rubyBaseStartTag:

		e.tokType = rubyParentStartToken
		end = len(rubyBaseStartStr)

	case endOfLine:
		e.tokType = endOfLineToken
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(lineBreakStr)

	case gaijiMarker:

		e.tokType = gaijiToken
		end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])

	case noteStartTag:

		e.tokType = noteToken
		end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])

	case bibInfoTag:
		e.tokType = bibInfoToken
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(bibInfoStartStr)

	case accentStartTag:
		e.tokType = accentToken
		end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])

	case kunojiTag:
		e.tokType = kunojiToken
		end = len(kunojiStr)

	case kunojiDakuTag:
		e.tokType = kunojiToken
		end = len(kunojiDakuStr)

	default:

		e.tokType = textToken
		end = findContiguousText(tknz.data[tknz.position:])

	}
	e.content = tknz.data[tknz.position : tknz.position+end]

	if e.tokType == noteToken {

		if e.innerString() == mainTextEndStr {

			e.tokType = bibInfoToken

		}
	}

	tknz.position = tknz.position + end

	//	fmt.Print(e)

	return e
}

func (t *token) fixEmptyText() {

	for pos := t.firstToken(); pos != nil; pos = pos.next {

		if pos.tokType != textToken {
			continue
		}

		if pos.String() == "" {
			pos = pos.next
			pos.prev.remove()
		}
	}

	return
}
