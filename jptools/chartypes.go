package jptools

//CharTypeID represents character types.
type CharTypeID int

//Define charater types.
const (
	Symbol CharTypeID = 1 << iota //Symbol captures everything that isn't captured by the other categories.
	KanNum
	ArabNumF
	ArabNum
	Hiragana
	Katakana
	KatakanaH
	Kanji
	LatinF
	Latin
)

//Name retunrs the character type of c
//as a string. E.g. "Katakana" if c
//is katakana.
func (c CharTypeID) Name() string {
	switch c {
	case Symbol:
		return "Symbol"
	case KanNum:
		return "Chinese Numeral"
	case Kanji:
		return "Kanji"
	case Hiragana:
		return "Hiragana"
	case Katakana:
		return "Katakana"
	case KatakanaH:
		return "Katakana Half Width"
	case LatinF:
		return "Latin Full Width"
	case Latin:
		return "Latin standard"
	case ArabNumF:
		return "Arabic Numeral Full Width"
	case ArabNum:
		return "Arabic Numeral Half Width"
	default:
		return "other"
	}
}

//IsKatakana checkes if r is katakana.
func IsKatakana(r rune) bool {

	return CharType(r) == Katakana

}

//IsHiragana checkes if r is hiragana,
func IsHiragana(r rune) bool {

	return CharType(r) == Hiragana

}

//ToKatakana convers r to katakana (iff. r is hiragana).
func ToKatakana(r rune) rune {

	if !IsHiragana(r) {

		return r
	}

	return r + 0x0060

}

//ToHiragana converts r to hiragana (iff. r is katakana)
func ToHiragana(r rune) rune {

	if !IsKatakana(r) {
		return r
	}

	if r > 0x30f6 {
		return r
	}

	return r - 0x0060

}

//CharType returns the character type of r according to the above classification.
//The Kanji numerals are identified both as Kanji and KanNum using the bitmask
// KanNum | Kanji.
func CharType(r rune) CharTypeID {
	if 0xFF9F <= r {
		return Symbol
	}

	if 0xFF66 <= r {
		return KatakanaH
	}

	if 0xFF5B <= r {
		return Symbol
	}

	if 0xFF41 <= r {
		return LatinF
	}

	if 0xFF3B <= r {
		return Symbol
	}

	if 0xFF21 <= r {
		return LatinF
	}

	if 0xFF1A <= r {
		return Symbol
	}

	if 0xFF10 <= r {
		return ArabNumF
	}

	if 0xFA6E <= r {
		return Symbol
	}

	if 0xF900 <= r {
		return Kanji
	}

	if 0x9FF0 <= r {
		return Symbol
	}

	if 0x4E00 <= r {
		if 0x767F <= r {
			return Kanji
		}

		if r == 0x767E {
			return KanNum | Kanji
		}

		if r == 0x56DB {
			return KanNum | Kanji
		}

		if r == 0x5343 {
			return KanNum | Kanji
		}

		if r == 0x5341 {
			return KanNum | Kanji
		}

		if r == 0x516D {
			return KanNum | Kanji
		}

		if r == 0x516B {
			return KanNum | Kanji
		}

		if r == 0x5146 {
			return KanNum | Kanji
		}

		if r == 0x5104 {
			return KanNum | Kanji
		}

		if r == 0x4E94 {
			return KanNum | Kanji
		}

		if r == 0x4E8C {
			return KanNum | Kanji
		}

		if r == 0x4E5D {
			return KanNum | Kanji
		}

		if r == 0x4E09 {
			return KanNum | Kanji
		}

		if r == 0x4E07 {
			return KanNum | Kanji
		}

		if r == 0x4E03 {
			return KanNum | Kanji
		}

		if r == 0x4E00 {
			return KanNum | Kanji
		}

		return Kanji
	}

	if 0x30FF <= r {
		return Symbol
	}

	if 0x30A1 <= r {
		return Katakana
	}

	if 0x30A0 <= r {
		return Symbol
	}

	if 0x3041 <= r {
		return Hiragana
	}

	if 0x3007 <= r {
		return Symbol
	}

	if 0x3005 <= r {
		return Kanji
	}

	if 0x7B <= r {
		return Symbol
	}

	if 0x61 <= r {
		return Latin
	}

	if 0x5B <= r {
		return Symbol
	}

	if 0x41 <= r {
		return Latin
	}

	if 0x3A <= r {
		return Symbol
	}

	if 0x30 <= r {
		return ArabNum
	}

	return Symbol
}
