package aozoratext

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type opt int

var o_full, o_raw, o_jis0208 bool

type tokenizer struct {
	data     string
	position int
}

var tokenizerOption opt

var tokenizerLog *log.Logger

var dummyLog *strings.Builder

func init() {

	o_full = true

	o_raw = false

	o_jis0208 = false

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

func tokenize(s string) *token {

	tknz := newTokenizer(s)

	t0 := tknz.nextToken()

	for t1 := t0; !tknz.empty(); t1 = t1.next {

		t1.insertTokenRight(tknz.nextToken())

	}

	return t0

}

func tokenizeAll(text string, lineOffset int) *token {

	fixline := fixfunc()

	lines := strings.Split(text, "\n")

	if len(lines) == 0 {
		return new(token)
	}

	lineTokens := newToken()

	tokenstring := newToken()

	for counter := lineOffset; counter < len(lines); counter++ {

		//fmt.Println(counter + 1)

		tokenizerLog.SetPrefix("line " + strconv.Itoa(counter+1) + ": ")

		lineTokens = fixline(tokenize(lines[counter]))

		lineTokens.firstToken().lineNo = counter

		tokenstring.joinTokens(lineTokens)

	}

	tokenstring.fixEmptyText()

	if !o_raw {

		tokenstring.fixDocument()

		tokenstring.insertSectionEnds()

	}

	tokenstring.firstToken().remove()

	return tokenstring

}

func newTokenizer(d string) *tokenizer {

	r := new(tokenizer)

	r.data = d

	r.position = 0

	return r
}

func (tknz *tokenizer) empty() bool {

	return tknz.position == len(tknz.data)

}

func (tknz *tokenizer) nextToken() (e *token) {

	e = newToken()

	end := 0

	switch typeOf(tknz.data[tknz.position:]) {

	case emptyLine:

		e.tokType = emptyLineToken
		end = 0

	case rubyStartTag:

		e.tokType = rubyStartToken
		end = len(rubyStartStr)

	case rubyEndTag:
		e.tokType = rubyEndToken
		end = len(rubyEndStr)

	//	end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])

	case rubyBaseStartTag:

		e.tokType = rubyParentStartToken
		end = len(rubyBaseStartStr)

	case gaijiMarker:

		e.tokType = gaijiToken
		end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])

	case noteStartTag:

		e.tokType = noteToken
		end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])

	case bibInfoTag:
		e.tokType = bibInfoToken
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
	/*
		switch e.tokType {

		case textToken:
			end = findContiguousText(tknz.data[tknz.position:])

		case emptyLineToken:
			end = 0

		case bibInfoToken:
			end = len(bibInfoStartStr)

		case rubyParentStartToken:
			end = len(rubyBaseStartStr)

		default:
			end = findMatchingCloser(typeOf(tknz.data[tknz.position:]), tknz.data[tknz.position:])
		}
	*/
	e.content = tknz.data[tknz.position : tknz.position+end]

	tknz.position = tknz.position + end

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
