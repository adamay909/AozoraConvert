package aozoraConvert

import (
	"strings"
)

func CheckStructure(text string) {

	return

}

func ListRawTokens(text string) string {

	t := tokenize(text)

	return listAllTokens(t)

}

func ListProcessedTokens(text string) string {

	t, err := tokenizeAndFix(text)

	if err != nil {

		return "ERROR:\n" + err.Error()

	}

	return listAllTokens(t)

}

func listAllTokens(t *token) string {

	output := new(strings.Builder)

	for e := t.firstToken(); e != nil; e = e.next {

		addToStringsBuilder(output, e.info(), "\n")

	}
	return output.String()

}
