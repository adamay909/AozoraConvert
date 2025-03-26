package aozoratext

import (
	"strings"
)

func CheckStructure(text string) {

	return

}

func ListRawTokens(text string) string {

	offset := aztextMainStart(text)

	d := strings.Join(strings.Split(text, "\n")[offset:], "\n")
	t := tokenize(d, offset)

	return listAllTokens(t)

}

func listAllTokens(t *token) string {

	output := new(strings.Builder)

	for e := t.firstToken(); e != nil; e = e.next {

		addToStringsBuilder(output, e.info(), "\n")

	}
	return output.String()

}

func ListProcessedTokens(text string) string {

	offset := aztextMainStart(text)

	t := tokenizeAndFix(text, offset)

	return listAllTokens(t)

}
