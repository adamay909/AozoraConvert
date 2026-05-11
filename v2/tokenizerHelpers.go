package aozoraconvert

import (
	"errors"
	"strings"
)

func (t *token) cleanup() {

	var stopProcessing bool

	withinAccent := false

	accentStart := new(token)

	converted := false

	for e := t.firstToken(); e != nil; e = e.next {

		//		fmt.Print(e)

		e.fixkunoji()

		if e.tokType == accentStartToken {
			if withinAccent {
				if !oTolerant {
					panic(errors.New(e.lineNumberStr() + "行：アクセント注記が入れ子"))
				} else {
					clog.Println(e.lineNumberStr() + "行：〔 => 外字注記")
					accentStart.tokType = specialCharToken
				}
			}
			withinAccent = true
			accentStart = e
			converted = false
			continue
		}

		if e.tokType == accentEndToken {

			if accentStart == nil {
				clog.Println(e.lineNumberStr() + "行：〕 => 外字注記")
				e.tokType = specialCharToken
				continue
			}

			if !converted {
				e.tokType = specialCharToken
				accentStart.tokType = specialCharToken
			}
			withinAccent = false
			accentStart = nil
			converted = false
			continue
		}

		if withinAccent && e.tokType == textToken {
			e.convertAccentToken()
			if e.content != e.unicodeContent {
				converted = true
			}
		}

		e.fixgaiji()

		stopProcessing = e.fixBibInfo()

		if stopProcessing {
			break
		}
	}
	if withinAccent {
		accentStart.tokType = specialCharToken
	}
}

// fix ruby notes so that all ruby parent start/end, ruby text
// start end are explicit
func (t *token) fixRuby() {

	for e := t.firstToken(); e != nil; e = e.next {
		//we ignore notes
		if e.tokType == noteStartToken {
			e = e.matchingNoteEnd()
			continue
		}
		//bibliographical info at end has no rubys
		if e.tokType == bibInfoToken {
			break
		}
		e.fixruby()

	}
}

func (t *token) gatherNotes() {

	for e := t.firstToken(); e != nil; e = e.next {

		if e.tokType == noteEndToken {
			if !oTolerant {
				panic(errors.New("終わり角括弧は外字注記にしてください"))
			}
			e.tokType = specialCharToken
			clog.Println(e.lineNumberStr() + "行：終わり角括弧 => 外字注記")
		}

		if e.tokType != noteStartToken {
			continue
		}

		f := e.matchingNoteEnd()

		if f == nil {
			panic(errors.New("note is not closed" + t.info()))
		}

		note := getnote(e, f)

		replaceTokens(e, f, note)

		note.inserted = false

		e = note

		if e.innerString() == mainTextEndStr {
			e.tokType = endMarkerToken
		}
	}

}

// getnote returns a noteToken by gathering information from start through end.
// start must be noteStartToken and end must be noteEndToken.
func getnote(start, end *token) *token {

	note := newTokenOfType(noteToken)

	note.content = noteStartStr

	note.unicodeContent = noteStartStr

	for e := start.next; e != end; e = e.next {

		note.content = note.content + e.String()

		note.unicodeContent = note.unicodeContent + e.unicodeString()

	}

	note.content = note.content + noteEndStr

	note.unicodeContent = note.unicodeContent + noteEndStr

	return note

}

func (t *token) cleanupTokenString() {

	e := t.firstToken().next

	if e == nil {
		return
	}

	for ; e != nil; e = e.next {

		switch {

		case e.prev.tokType == emptyToken:
			e.prev.remove()

		case e.tokType == rubyParentStartToken:
			if e.prev == nil || e.prev.tokType != rubyGroupStartToken {
				if !oTolerant {
					panic(errors.New(e.lineNumberStr() + "行：ルビではないルビ親開始記号"))
				}

				clog.Println(e.lineNumberStr() + "行：｜ => 外字注記")

				e.tokType = specialCharToken
			}

		case e.tokType == noteEndToken:
			e.tokType = specialCharToken

		}
	}
}

// replaceTokens replaces the tokens from start to end with replacement.
func replaceTokens(start, end, replacement *token) {

	start.insertTokenLeft(replacement)

	replacement.next = end.next

	if end.next != nil {
		end.next.prev = replacement
	}

}

func (t *token) isAbbreviatedBlockFormat() bool {

	if !t.isFirstTokenInLine() {
		return false
	}

	if t.tokType != noteToken {
		return false
	}

	if strings.HasPrefix(t.innerString(), blockStartStr) {
		return false
	}

	for _, e := range impliedCloserMarker {

		if strings.HasSuffix(t.innerString(), e) {

			return true
		}
	}

	return false
}

func (t *token) isImpliedOpener() (bool, string) {

	if t.tokType != noteToken {
		return false, ""
	}

	if strings.HasPrefix(t.innerString(), blockStartStr) {
		return false, ""
	}

	for _, e := range impliedOpenerMarker {
		if strings.HasSuffix(t.innerString(), e) {
			return true, e
		}
	}

	return false, ""
}

// just get the boolean
func (t *token) isImpliedOpenerTF() bool {
	ok, _ := t.isImpliedOpener()
	return ok
}
