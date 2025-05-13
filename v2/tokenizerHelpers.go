package aozoraconvert

func (t *token) cleanup() {

	var stopProcessing bool

	withinAccent := false

	accentStart := new(token)

	converted := false

	for e := t.firstToken(); e != nil; e = e.next {

		//	fmt.Print(e)

		e.fixkunoji()

		if e.tokType == accentStartToken {
			if withinAccent {
				panic(e.lineNumberStr() + "行：アクセント注記が入れ子")
			}
			withinAccent = true
			accentStart = e
			converted = false
			continue
		}

		if e.tokType == accentEndToken {

			if accentStart == nil {
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

		if withinAccent && t.tokType == textToken {
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
}

func (t *token) fixRuby() {

	for e := t.firstToken(); e != nil; e = e.next {

		//		fmt.Print(e)

		if e.tokType == noteStartToken {

			e = e.matchingNoteEnd()

			continue

		}
		//		fmt.Print(e)

		if e.tokType == bibInfoToken {
			break
		}

		e.fixruby()

	}
}

func (t *token) gatherNotes() {

	for e := t.firstToken(); e != nil; e = e.next {

		if e.tokType != noteStartToken {
			continue
		}

		f := e.matchingNoteEnd()

		if f == nil {
			panic("note is not closed" + t.info())
		}

		note := getnote(e, f)

		replaceTokens(e, f, note)

		e = note

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

func (t *token) removeEmptyToken() {

	e := t.firstToken().next

	if e == nil {
		return
	}

	for ; e != nil; e = e.next {

		if e.prev.tokType == emptyToken {
			e.prev.remove()
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
