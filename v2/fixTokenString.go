package aozoraconvert

import (
	"strconv"
	"strings"
)

func (t *token) fixLines() {

	t.cleanup()

	t.fixRuby()

	t.gatherNotes()

	t.cleanupTokenString()

	if oTolerant {
		for e := t.firstToken(); e != nil; e = e.next {
			e.fixnote()
		}
	}

	for e := t.firstToken(); e != nil; e = e.next {

		//		fmt.Print(e)

		e.fixAbbreviatedBlockFormat()
		e.fixInlineAlignBottom()
		e.fixImpliedOpener()
	}

	t.firstToken().fixParagraphs()

}

func (t *token) fixAbbreviatedBlockFormat() {

	if !t.isAbbreviatedBlockFormat() {
		return
	}

	m := emptyStr

	for _, e := range impliedCloserMarker {

		if strings.HasSuffix(t.innerString(), e) {

			m = e

			break
		}
	}
	//now we know we are dealing with a relevant token

	//first fix t itself.
	sbs := new(strings.Builder)

	addToStringsBuilder(sbs, blockStartStr, t.innerString())

	t.originalContent = t.content

	t.setInnerString(sbs.String())

	t.unicodeContent = t.content

	t.jis0213Content = t.content

	t.modified = true

	//need to check if we have to insert a new closing token.

	sbe := new(strings.Builder)

	addToStringsBuilder(sbe, blockEndStr, m, formatEndStr)

	//we need to insert a new closing token.

	end := t.lastTokenInLine()

	end.insertTokenRight(newNote(sbe.String()))

	if end.next.next != nil && end.next.next.tokType != endOfLineToken {
		end.next.insertTokenRight(newTokenOfType(endOfLineToken))
	}

	//need to insert new eol
	if t.next.tokType != endOfLineToken {
		t.insertTokenRight(newTokenOfType(endOfLineToken))
	}

	return

}

func (t *token) fixInlineAlignBottom() {

	if t.isFirstTokenInLine() {

		return

	}

	if t.tokType != noteToken {

		return

	}

	m := emptyStr

	for _, e := range bottomalignMarker {

		if strings.HasSuffix(t.innerString(), e) {

			m = e

			break
		}
	}

	if m == emptyStr {
		return
	}

	//now we know we are dealing with a relevant token.

	t.lastTokenInLine().insertTokenLeft(newTokenOfType(alignBottomCloserToken))

	return

}

func (t *token) fixkunoji() {

	if t.tokType != kunojiToken {
		return
	}

	if t.content == kunojiStr {

		t.unicodeContent = kunojiStrU

		t.jis0213Content = kunojiStrU

		t.modified = true

		return
	}

	t.unicodeContent = kunojiDakuStrU

	t.jis0213Content = kunojiDakuStrU

	t.modified = true
}

func (t *token) fixruby() {

	if t.tokType != rubyStartToken {
		return
	}

	e := new(token)

	found := false

	for e = t.next; e != nil; e = e.next {

		if e.tokType == rubyEndToken {

			if e == t.next {

				if !oTolerant {
					panic(t.lineNumberStr() + "行：ルビの文字列が指定されていません")
				}

				clog.Println(t.lineNumberStr() + "行：空のルビを削除")

				t.tokType = emptyToken

				e.tokType = emptyToken

				return
			}

			e.insertTokenRight(newTokenOfType(rubyGroupEndToken))

			found = true

			break
		}

		if e.tokType == endOfLineToken {
			break
		}

		if e.tokType == emptyLineToken {
			break
		}

		if e.tokType == rubyStartToken {
			break
		}

		if e.tokType == rubyParentStartToken {
			break
		}

	}
	if !found {

		if !oTolerant {
			panic(t.lineNumberStr() + "行：ルビ文字列の終了位置が指定されていません。" + t.info() + "\n surrounding text: " + t.textContext())
		}

		clog.Println(t.lineNumberStr() + "行：《　=> 外字注記")

		t.tokType = specialCharToken

		return

	}

	if e2 := t.prev; e2.tokType == noteEndToken {

		noteStart := e2.matchingNoteStart()

		note := getnote(noteStart, e2)

		replaceTokens(noteStart, e2, note)

		note.inserted = false

		if note.isImage() {

			msg := t.lineNumberStr() + "行：" + " 画像を文字として使用？"

			note.tokType = gaijiImgToken

			msglog.Println("警告: " + msg)

		}

	}

	t.insertTokenLeft(newTokenOfType(rubyParentEndToken))

	if !t.rubyParentExplicit() {

		t.insertRubyParentStart()

	}

	return

}

func (t *token) rubyParentExplicit() bool {

	for e := t.prev; e != nil; e = e.prev {

		if e.tokType == textToken {
			continue
		}

		if e.tokType == rubyParentStartToken {
			e.insertTokenLeft(newTokenOfType(rubyGroupStartToken))
			return true
		}

		if e.tokType == rubyParentEndToken {
			continue
		}

		if e.tokType == gaijiCharToken {
			continue
		}

		if e.tokType == kunojiToken {
			continue
		}

		if e.tokType == gaijiNoteToken {
			continue
		}

		if e.tokType == accentStartToken {
			continue
		}

		if e.tokType == accentEndToken {
			continue
		}

		if e.tokType == noteEndToken {
			continue
		}

		if e.tokType == noteStartToken {
			continue
		}

		if e.tokType == noteToken {
			continue
		}

		if e.tokType == emptyToken {
			continue
		}

		if e.tokType == gaijiImgToken {
			continue
		}

		break
	}

	return false
}

func (t *token) insertRubyParentStart() {

	e := new(token)

	e = t.prev

	if e.prev.tokType == gaijiImgToken {

		e.prev.insertTokenLeft(newTokenOfType(rubyParentStartToken))

		e.prev.prev.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

		return

	}

	var refType CharTypeID

	var r []rune

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == textToken {
			break
		}

		if e.tokType == gaijiCharToken {
			break
		}

	}

	r = []rune(e.unicodeString())

	refType = CharType(r[len(r)-1])

	k := 0

	for e = t.prev.prev; e != nil; e = e.prev {

		if e.tokType == noteEndToken {

			noteStart := e.matchingNoteStart()

			if noteStart == nil {

				if !oTolerant {
					panic(e.lineNumberStr() + "行：角括弧は外字注記に！")
				}

				clog.Println(e.lineNumberStr() + "行：角括弧 => 外字注記")

				continue

			}

			note := getnote(noteStart, e)

			replaceTokens(noteStart, e, note)

			e = note

		}

		if e.tokType == rubyGroupEndToken {
			break
		}

		if e.tokType == endOfLineToken {
			break
		}

		if e.tokType == emptyLineToken {
			break
		}

		if e.tokType == kunojiToken {
			break
		}

		if e.tokType == specialCharToken {
			break
		}

		if e.isPairClose() {
			break
		}

		if e.isPairOpen() {
			break
		}

		if e.isBlock() {
			break
		}

		if e.isBlockEnd() {
			break
		}

		if e.tokType == noteToken {
			continue
		}

		if e.tokType != textToken && e.tokType != gaijiCharToken {
			continue
		}

		r = []rune(e.unicodeString())

		if refType == Roman {

			for k = len(r) - 1; CharType(r[k]) != Whitespace; k-- {

				if k == 0 {

					break
				}
			}

		} else {

			for k = len(r) - 1; refType == CharType(r[k]); k-- {

				if k == 0 {

					break
				}
			}
		}

		if (k == len(r)-1 && k > 0) || (len(r) == 1 && CharType(r[k]) != refType) {

			for e = e.next; ; e = e.next {

				if e.tokType != noteToken {
					break
				}

				if e.isAbbreviatedBlockFormat() {
					continue
				}

				if ok, _ := e.isImpliedOpener(); ok {
					continue
				}

				if e.isBlock() {
					continue
				}

				if e.isBlockEnd() {
					continue
				}

				if !e.isPairOpen() {
					continue
				}

				if e.isFormatOfType(bottomalignMarker) {
					continue
				}

				break

			}

			e.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

			e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

			return

		}

		if CharType(r[k]) != refType {

			t2 := newTokenOfType(textToken)

			t2.setString(string(r[:k+1]))

			e.originalContent = e.content

			e.setString(string(r[k+1:]))

			e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

			e.prev.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

			e.prev.prev.insertTokenLeft(t2)

			e.modified = true

			return

		}

	}

	if e == nil {

		e = t.firstToken()

	}

	if e.tokType == endOfLineToken || e.tokType == emptyLineToken || e.tokType == rubyGroupEndToken {

		e = e.next
	}

	for ; ; e = e.next {

		if e.tokType != noteToken {
			break
		}

		if e.isAbbreviatedBlockFormat() {
			continue
		}

		if e.isBlock() {
			continue
		}

		if e.isBlockEnd() {
			continue
		}

		if ok, _ := e.isImpliedOpener(); ok {
			continue
		}

		if !e.isPairOpen() {
			continue
		}

		if e.isFormatOfType(bottomalignMarker) {
			continue
		}

		break

	}

	e.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

	e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

	return
}

func (t *token) reformgaiji() {

	if t.tokType != gaijiToken {

		return

	}

	t.setString(noteStartStr + t.innerString() + noteEndStr)

	t.tokType = gaijiNoteToken

	t.modified = true

	n := newTokenOfType(textToken)

	n.setString(referenceMarkStr)

	t.insertTokenLeft(n)

}

func (t *token) fixgaiji() {

	if t.tokType != gaijiToken {
		return
	}

	t.replaceGaiji()

	if t.tokType == gaijiToken {

		msglog.Println("gaiji conversion failed:", t.info())

		t.reformgaiji()

	}
}

func (t *token) replaceGaiji() {

	var uni string

	originalNote := t.innerString()

	j := jisCodeOf(originalNote)

	if j != "" {
		uni, _ = convert(j)

		t.jis0213Content = uni

		t.unicodeContent = uni

	} else {

		u := ucode(originalNote)

		if u == "" {
			return
		}

		uni = unicodeOf(u)

		t.unicodeContent = uni

	}

	switch uni {

	case "《":
		t.tokType = specialCharToken

	case "》":
		t.tokType = specialCharToken

	case "［":
		t.tokType = specialCharToken

	case "］":
		t.tokType = specialCharToken

	case "〔":
		t.tokType = specialCharToken

	case "〕":
		t.tokType = specialCharToken

	case "｜":
		t.tokType = specialCharToken

	case "＃":
		t.tokType = specialCharToken

	case "※":
		t.tokType = specialCharToken

	default:
		t.tokType = gaijiCharToken
	}

	if t.tokType == specialCharToken {

		t.content = t.unicodeContent

	}

	t.modified = true

	return
}

// convert everything to nested formatting notes
// to ease parsing of notes. The returned token is new opener.
func (t *token) fixImpliedOpener() {

	ok, m := t.isImpliedOpener()

	if !ok {
		return
	}

	switch m {

	case "ルビ":
		fixLeftRuby(t)

	case "注記":
		fixChuki(t)

	case "大きな文字":
		fixFontSize(t)

	case "小さな文字":
		fixFontSize(t)

	default:
		fixopener(t, m)

	}

	return
}

func fixopener(note *token, m string) {

	rs := note.getRefStrings()

	if rs[0] == "" {
		return
	}

	left := false

	if strings.HasPrefix(strings.TrimPrefix(note.innerString(), refStartStr+rs[0]+refEndStr), "の左に") {
		left = true
	}

	b := new(strings.Builder)

	if left {
		addToStringsBuilder(b, "左に")
	}

	addToStringsBuilder(b, m, formatEndStr)

	note.originalContent = note.content

	note.setInnerString(b.String())

	note.unicodeContent = note.content

	note.jis0213Content = note.content

	note.modified = true

	nt := new(token)

	if left {

		nt = newNote("左に" + m)

	} else {

		nt = newNote(m)
	}

	note.addTokenBefore(rs[0], nt)

	return
}

func fixFontSize(note *token) {

	rs := note.getRefStrings()

	if rs[0] == "" {
		return
	}

	sb := strings.Split(note.innerString(), refEndStr+"は")[1]

	note.originalContent = note.content

	if strings.HasSuffix(sb, "大きな文字") {

		note.setInnerString("大きな文字終わり")

		note.unicodeContent = note.content

		note.jis0213Content = note.content

	} else {

		note.setInnerString("小さな文字終わり")

		note.unicodeContent = note.content

		note.jis0213Content = note.content
	}

	note.modified = true

	note.addTokenBefore(rs[0], newNote(sb))

	return
}

func fixLeftRuby(note *token) *token {

	rs := note.getRefStrings()

	if len(rs) != 2 {
		return note
	}

	//	pt := note.prevTextToken()

	b := new(strings.Builder)

	addToStringsBuilder(b, "［＃左に", refStartStr, rs[1], refEndStr, "のルビ付き終わり］")

	note.originalContent = note.content

	note.setString(b.String())

	note.unicodeContent = note.content

	note.jis0213Content = note.content

	note.modified = true

	nt := newNote("左にルビ付き")

	note.addTokenBefore(rs[0], nt)

	return note
}

func fixChuki(note *token) *token {

	var left bool

	rs := note.getRefStrings()

	if len(rs) != 2 {
		return note
	}

	if strings.HasPrefix(strings.TrimPrefix(note.innerString(), refStartStr+rs[0]+refEndStr), "の左に") {
		left = true
	}

	//	pt := note.prevTextToken()

	b := new(strings.Builder)

	if left {
		addToStringsBuilder(b, "左に")
	}

	addToStringsBuilder(b, refStartStr, rs[1], refEndStr, "の注記付き終わり")

	note.originalContent = note.content

	note.setInnerString(b.String())

	note.unicodeContent = note.content

	note.jis0213Content = note.content

	note.modified = true

	nt := new(token)

	if left {

		nt = newNote("左に注記付き")

	} else {

		nt = newNote("注記付き")

	}

	note.addTokenBefore(rs[0], nt)

	return note
}

func (t *token) fixParagraphs() {

	start := new(token)
	end := new(token)

	for start = t.firstToken(); start.next != nil; start = end.next {

		if start.tokType == emptyLineToken || start.tokType == endOfLineToken {
			end = start
			if end.next == nil {
				break
			}
			continue
		}

		end = start.nextTokenOfType(endOfLineToken)

		if end == nil {
			break
		}

		fixParagraph(start, end)

		for pos := end.prev; pos.isEmptyText() || pos.tokType == emptyToken; pos = pos.prev {
			pos.tokType = emptyToken
		}

		if end.next == nil {
			break
		}
	}
}

// start should be non-eol and end should be eol
func fixParagraph(start, end *token) {

	for start != nil && start.tokType == emptyToken {
		start = start.next
	}

	if start == nil {
		return
	}

	if start.tokType == noteToken && strings.HasPrefix(start.innerString(), blockStartStr) {
		return
	}

	if start.tokType == noteToken && strings.HasPrefix(start.innerString(), blockEndStr) && strings.HasSuffix(start.innerString(), formatEndStr) {
		return
	}

	oldstart := start

	for start.isEmptyText() || start.tokType == emptyToken {
		start = start.next
	}

	if start == end {

		end.tokType = emptyLineToken

		for pos := oldstart; pos != end; pos = pos.next {
			pos.tokType = emptyToken
		}

		return
	}

	switch start.tokType {

	case textToken, rubyGroupStartToken, gaijiCharToken, specialCharToken, accentStartToken, gaijiImgToken:

		if oldstart.prev != nil && oldstart.prev.tokType == endOfLineToken {
			oldstart.prev.tokType = paragraphToken
		} else {
			start.insertTokenLeft(newParagraphToken())
		}

		end.tokType = paragraphEndToken

		end.modified = true

		return

	case bibInfoToken:

		if start.String() != "" {
			return
		}

		start.insertTokenRight(newParagraphToken())

		end.tokType = paragraphEndToken

		end.modified = true

		return

	}

	if start.isSectionTitleStart() {
		//return
		if end.prev.isSectionTitleEnd() {
			return
		}

		c := new(token)

		for c = end.prev; !c.isSectionTitleEnd(); c = c.prev {
			if c.prev == nil {
				break
			}
		}

		c.insertTokenRight(newTokenOfType(paragraphToken))

		end.tokType = paragraphEndToken

		end.modified = true

		return

	}

	if start.isFormatOfType(captionMarker) {
		return
	}

	pos := start

	//	if pos.isBlockStartNote() {

	if strings.HasPrefix(pos.innerString(), blockStartStr) {

		for e := pos; e != end; e = e.next {

			if e.isSectionTitleStart() {
				return
			}

			if strings.HasPrefix(pos.innerString(), blockStartStr) {
				//			if e.isBlockStartNote() {
				continue
			}

			if e.tokType != noteToken {

				pos.insertTokenLeft(newParagraphToken())

				end.tokType = paragraphEndToken

				end.modified = true

				return
			}

			if e.next == end {
				return
			}
		}
		return

	}

	if pos.tokType == noteToken {

		for e := pos.next; e != end; e = e.next {

			if e.tokType != noteToken {

				pos.insertTokenLeft(newParagraphToken())

				end.tokType = paragraphEndToken

				end.modified = true

				return

			}
		}
	}
	return
}

func (t *token) fixblockformatting() {

	pos := t

	pos = pos.matchingIndentationCloser()

	if !pos.isIndentationEnd() {
		pos.insertTokenLeft(newIndentationCloser())
		pos.insertTokenLeft(newEolToken())

	}

	return
}

func (t *token) fixCentering() {

	pos := t

	for ; !pos.next.isPagination() && pos.next.tokType != mainTextEndToken; pos = pos.next {
	}

	pos.next.insertTokenLeft(newTokenOfType(centeringEndToken))

	return
}

func (t *token) fixSectionTitles() {

	seccount := 0

	for e := t.firstToken(); e != nil; e = e.next {

		if e.isSectionTitleStart() {

			seccount++

			e.fixSectionFormatting(seccount)

		}
	}
}

func (t *token) fixSectionFormatting(seccount int) {

	pos := t

	e1 := new(token)

	if seccount == 1 {

		for e1 = pos; e1 != nil && e1.tokType != mainTextStartToken; e1 = e1.prev {
		}

		if e1 == nil {
			e1 = t.firstToken()
		}
	} else {

		e1 = sectionTitleFormattingStart(pos)
	}

	switch strings.TrimPrefix(pos.innerString(), blockStartStr) {
	case "大見出し":
		e1.insertTokenRight(newTokenOfType(sectionToken))

	case "中見出し":
		e1.insertTokenRight(newTokenOfType(subsectionToken))

	case "小見出し":
		e1.insertTokenRight(newTokenOfType(subsubsectionToken))
	}

	return

}

func (t *token) insertSectionEnds() {

	var open int

	open = 0

	var end *token
	var openers []*token

	for pos := t.firstToken(); pos.next != nil; pos = pos.next {

		if pos.isBlockStartNote() {
			openers = append(openers, pos)
			open++
		}

		if pos.isBlockClosingNote() {
			open--
			if open < 0 {
				panic(pos.info() + " closing tag without matching opener! ")
			}

			openers = openers[:len(openers)-1]
		}

		if !pos.isSectionStart() {
			continue
		}

		if open != 0 {
			//ensure proper nesting
			end = openers[len(openers)-1].matchingCloserToken()
		} else {
			end = nil
		}

		switch pos.tokType {

		case sectionToken:
			for pos2 := pos.next; pos2 != nil; pos2 = pos2.next {

				if pos2.tokType == sectionToken {
					pos2.insertTokenLeft(newTokenOfType(sectionEndToken))
					break
				}
				if pos2.tokType == bibInfoToken || pos2.tokType == mainTextEndToken {
					pos2.insertTokenLeft(newTokenOfType(sectionEndToken))
					break
				}
				if pos2 == end {
					pos2.insertTokenLeft(newTokenOfType(sectionEndToken))
					break
				}
				if pos2.next == nil {
					pos2.insertTokenRight(newTokenOfType(sectionEndToken))
					break
				}
			}

		case subsectionToken:
			for pos2 := pos.next; pos2 != nil; pos2 = pos2.next {

				if pos2.tokType == subsectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2.tokType == sectionEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2.tokType == sectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2.tokType == bibInfoToken || pos2.tokType == mainTextEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2 == end {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}

				if pos2.next == nil {
					pos2.insertTokenRight(newTokenOfType(subsectionEndToken))
					break
				}

			}

		case subsubsectionToken:
			for pos2 := pos.next; pos2 != nil; pos2 = pos2.next {

				if pos2.tokType == subsubsectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == subsectionEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == subsectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == sectionEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == sectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == bibInfoToken || pos2.tokType == mainTextEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2 == end {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.next == nil {
					pos2.insertTokenRight(newTokenOfType(subsubsectionEndToken))
					break
				}

			}

		}
	}
}

func newIndentationCloser() *token {

	t := newNote("ここで字下げ終わり")

	return t

}

func sectionTitleFormattingStart(t *token) (s *token) {

	for s = t.prev; ; s = s.prev {

		if s == nil {
			t.firstToken().insertTokenLeft(newTokenOfType(emptyToken))
			return t.firstToken()
		}

		if s.tokType == mainTextStartToken {
			return s
		}

		if s.isPairOpen() {
			continue
		}

		if s.tokType == endOfLineToken {
			continue
		}

		if s.tokType == emptyLineToken {
			continue
		}

		if s.tokType == emptyToken {
			continue
		}

		if s.tokType == paragraphToken {
			s2 := s.matchingCloserToken()
			s2.tokType = emptyToken
			s.tokType = emptyToken
			continue
		}

		if s.isEmptyText() {
			continue
		}

		if s.isImage() {
			continue
		}

		break

	}

	for s = s.next; s.tokType == emptyLineToken || s.tokType == endOfLineToken || s.tokType == emptyToken; s = s.next {
	}

	return s.prev

}

func (t *token) fixaccent() {

	if t.tokType != accentStartToken {
		return
	}

	e := new(token)

	for e = t.next; e != nil; e = e.next {
		if e.tokType == accentEndToken {
			break
		}
	}

	if e == nil {
		t.tokType = textToken
		t.modified = true
		return
	}

	s := ""

	for f := t; f != e; f = f.next {

		if f.tokType == textToken {
			s = s + f.String()
		}
	}

	converted := convertAccent(s)

	if converted == s {
		t.tokType = specialCharToken
		e.tokType = specialCharToken
		t.modified = true
		e.modified = true
		return
	}

	for f := t.next; f != e; f = f.next {

		if f.tokType != textToken {
			continue
		}

		f.unicodeContent = convertAccent(f.String())

		f.jis0213Content = f.unicodeContent

		f.modified = true
	}

	t.next.content = accentStartStr + t.next.content

	t.tokType = emptyToken

	ignore := false

	for pos := e; pos != t; pos = pos.prev {

		if pos.tokType == rubyEndToken {
			ignore = true
			continue
		}
		if pos.tokType == rubyStartToken {
			ignore = false
			continue
		}

		if ignore {
			continue
		}

		if pos.tokType == textToken {

			pos.content = pos.content + accentEndStr

			pos.modified = true

			break
		}
	}

	e.tokType = emptyToken

	return
}

func (t *token) decorationLeft() bool {

	if t.tokType != noteToken {
		return false
	}

	return strings.HasPrefix(t.innerString(), "左に")

}

func (t *token) fixFigures() {

	e := t

	if t.nextSignificantToken().isCaption() {

		t.insertTokenLeft(newTokenOfType(figureStartToken))

		n := newTokenOfType(figureEndToken)

		for e = e.nextSignificantToken(); e != nil; e = e.next {

			if e.tokType != noteToken {
				continue
			}

			if strings.HasSuffix(e.innerString(), "キャプション終わり") {
				break
			}
		}

		e.insertTokenRight(n)

		return

	}

}

func (t *token) fixBibInfo() (stopProcessing bool) {

	if t.tokType != bibInfoToken {
		return false
	}

	pos := new(token)

	//cleanup any stray bibinfo tokens
	for pos = t.next; pos != nil; pos = pos.next {

		if pos.tokType == bibInfoToken {

			pos.tokType = textToken

			pos.modified = true

		}
	}

	//find end of bibliographical info
	for pos = t.lastToken(); pos.tokType != endOfLineToken; pos = pos.prev {
	}

	pos.insertTokenLeft(newTokenOfType(bibInfoEndToken))

	//if notetoken, we are done
	if t.innerString() == noteStartStr+mainTextEndStr+noteEndStr {
		t.tokType = endMarkerToken
		return true
	}

	//extra clean up
	t.tokType = textToken

	t.modified = true

	t.insertTokenLeft(newTokenOfType(bibInfoToken))

	return true
}

func (t *token) fixDocument() {

	for pos := t.firstToken(); pos != nil; pos = pos.next {

		switch {

		case pos.tokType == textToken:

			continue

		case pos.tokType == rubyStartToken:

			continue

		case pos.tokType == rubyEndToken:

			continue

		case pos.tokType == rubyParentStartToken:

			continue

		case pos.tokType == rubyParentEndToken:

			continue

		case pos.innerString() == "ページの左右中央":

			pos.fixCentering()

		case pos.isIndentationStart():

			pos.fixblockformatting()

		case pos.isImage():

			pos.fixFigures()

			/*		case pos.isSectionTitleStart():

					pos.fixSectionFormatting()*/
		}

	}

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

func (t *token) insertAozoraBookMarker() {

	e := new(token)

	for e = t.firstToken(); e.tokType != emptyLineToken; e = e.next {

		if e.next == nil {
			return
		}
	}

	if e.lineNumber() < 3 {
		return
	}

	e.insertTokenRight(newTokenOfType(mainTextStartToken))

	for ; e.tokType != bibInfoToken; e = e.next {

		if e.next == nil {
			t.lastToken().insertTokenRight(newTokenOfType(mainTextEndToken))
			return
		}
	}

	e.insertTokenLeft(newTokenOfType(mainTextEndToken))

	return

}

func (t *token) fixBlockFormat() {

	if t.tokType != noteToken {
		return
	}

}

func (t *token) fixnote() {

	if t.tokType != noteToken {
		return
	}

	r := []rune(t.innerString())

	k := len(r)
	for c := len(r) - 1; c > -1; c-- {
		if r[c] == ' ' {
			k = c
			continue
		}

		if r[c] == '　' {
			k = c
			continue
		}
		break
	}
	if k != len(r) {

		oldStr := t.String()

		t.setInnerString(string(r[:k]))

		clog.Println(t.lineNumberStr()+"行："+oldStr, "=>", t.String())
	}

	for _, e := range pairMarker {

		if strings.HasSuffix(t.innerString(), e+"おわり") || strings.HasSuffix(t.innerString(), e+"終り") {

			newstr := strings.TrimSuffix(t.innerString(), e+"おわり") + e + "終わり"

			clog.Println(t.lineNumberStr()+"行：", t.String(), "=>", newstr)

			t.setInnerString(newstr)

			break
		}
	}

	newstr := ""

	switch t.innerString() {

	case "字下げ終わり":

		newstr = ("ここで字下げ終わり")

	case "ここから割り注":

		newstr = "割り注"

	case "ここで割り注終わり":

		newstr = "割り注終わり"

	case "中見出終わり":

		newstr = "中見出し終わり"

	case "大見出終わり":

		newstr = "大見出し終わり"

	case "小見出終わり":

		newstr = "小見出し終わり"

	}

	if newstr != "" {
		clog.Println(t.lineNumberStr()+"行：", t.String(), "=>", noteStartStr+newstr+noteEndStr)

		t.setInnerString(newstr)
	}

	if t.isBlock() {

		if t.prev != nil && !t.prev.isLineBreak() {

			clog.Println(t.lineNumberStr()+"行: ", "ブロック開始注記の前に改行")

			t.insertTokenLeft(newTokenOfType(endOfLineToken))

		} else {

			if t.next != nil && !t.next.isLineBreak() {

				clog.Println(t.lineNumberStr()+"行: ", "ブロック開始注記の後に改行")

				t.insertTokenRight(newTokenOfType(endOfLineToken))

			}
		}
	}

	if t.isBlockEnd() {

		if t.prev != nil && !t.prev.isLineBreak() {

			clog.Println(t.lineNumberStr()+"行: ", "ブロック終了注記の前に改行")

			t.insertTokenLeft(newTokenOfType(endOfLineToken))

		} else {

			if t.next != nil && !t.next.isLineBreak() {

				clog.Println(t.lineNumberStr()+"行: ", "ブロック終了注記の後に改行")

				t.insertTokenRight(newTokenOfType(endOfLineToken))

			}
		}
	}

	if t.isPairOpen() && strings.HasSuffix(t.innerString(), "字下げ") {

		//if t.prev != nil && t.prev.tokType != endOfLineToken {
		if t.prev != nil && !t.prev.isLineBreak() {

			clog.Println(t.lineNumberStr()+"行：", "字下げ前に改行")

			t.insertTokenLeft(newTokenOfType(endOfLineToken))

		}
	}
	return
}

func (t *token) fixIndentationRound2() {

	c := 0

	for e := t.firstToken(); e != nil; e = e.next {

		if e.isIndentationStart() {
			c++
			continue
		}

		if !e.isIndentationEnd() {
			continue
		}

		c--

		if c != -1 {
			continue
		}

		c2 := 0

		nested := new(token)

		for f := e.prev; f != nil; f = f.prev {

			if !f.isIndentationStart() {
				continue
			}

			c2++

			if c2 > 2 {
				break
			}

			if c2 < 2 {

				nested = f
				/*
					if f.prev == nil {
						fmt.Println("check2")
						break
					}

					if !f.prev.isLineBreak() {
						fmt.Println("check3")
						break
					}

					if f.prev.prev == nil {
						fmt.Println("check4")
						break
					}

					if !f.prev.prev.isIndentationEnd() {
						fmt.Println("check5")
						break
					}
				*/
				continue
			}

			for i := nested.next; i != nil; i = i.next {
				if i.isIndentationEnd() {

					nested.setInnerString(blockStartStr + fwnum(f.getTopMargin()+nested.getTopMargin()) + "字下げ")

					i.insertTokenRight(newTokenOfType(endOfLineToken))
					i.next.insertTokenRight(copyOf(f))
					i.next.next.insertTokenRight(newTokenOfType(endOfLineToken))
					c++
					break
				}
			}
			break
		}

		if c2 != 2 {
			panic(e.lineNumberStr() + "行のあたり：余分な字下げ終了注記")
		}

		clog.Println(e.lineNumberStr() + "行のあたり：字下げが入れ子になっている模様。修復を試みる")

	}

}

func (t *token) getIndentation() int {

	str := strings.TrimPrefix(t.innerString(), "ここから")

	switch {

	case strings.Contains(str, "天付き") && strings.Contains(str, "折り返して"):

		return 0 - t.getTopMargin()

	case strings.Contains(str, "折り返して"):

		part := strings.Split(str, "折り返して")

		a, _ := strconv.Atoi(getFirstNumberString(part[0]))

		b, _ := strconv.Atoi(getFirstNumberString(part[1]))

		return a - b

	default:
		return 0

	}
}

func (t *token) getTopMargin() int {

	str := strings.TrimPrefix(t.innerString(), "ここから")

	switch {

	case strings.Contains(str, "天付き") && strings.Contains(str, "折り返して"):

		m, _ := strconv.Atoi(getFirstNumberString(str))

		return m

	case strings.Contains(str, "折り返して"):

		part := strings.Split(str, "折り返して")

		b, _ := strconv.Atoi(getNumberString(part[1]))

		return b

	default:

		str = strings.TrimSuffix(str, "字下げ")

		m, _ := strconv.Atoi(getNumberString(str))

		return m

	}
}

func (t *token) cleanUpIndentation() {

	for e := t.firstToken(); e != nil; e = e.next {

		if !e.isIndentationStart() {
			continue
		}

		hasContent := false

		f := new(token)

		for f = e.next; !f.isIndentationEnd(); f = f.next {
			if f.tokType != endOfLineToken {
				hasContent = true
				break
			}
		}

		if hasContent {
			continue
		}

		for f = e; !f.isIndentationEnd(); f = f.next {

			f.tokType = emptyToken

		}

		f.tokType = emptyToken

		e = f
	}
}
