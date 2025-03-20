package aozoratext

import (
	_ "embed" //for embedding
	"fmt"
	"log"
	"strings"
	"testing"
)

//go:embed testdata/aztestingTXT.txt
var tokenizerTestdata string

func TestTokenize(t *testing.T) {

	//	setTokenizerOption("raw")

	checkpairs := strings.Split(tokenizerTestdata, "\n&&&\n")

	for counter, p := range checkpairs {

		if strings.HasPrefix(p, "+++++") {
			break
		}

		s := strings.Split(p, "\n:::\n")

		if len(s) != 2 {
			continue
		}

		fmt.Println(counter, ": check:\n", s[0])

		b := tokenizeAll(strings.TrimSpace(s[0]), 0)

		if b == nil {

			log.Println("tokenization falied. Exiting.")

			return
		}

		numberTokens(b)

		fmt.Println(b.listAllTokens())

	}
}

func TestCheckAST(t *testing.T) {

	checkpairs := strings.Split(tokenizerTestdata, "\n&&&\n")

	for counter, p := range checkpairs {

		if strings.HasPrefix(p, "+++++") {
			break
		}

		s := strings.Split(p, "\n:::\n")

		if len(s) != 2 {
			continue
		}

		fmt.Println(counter+1, ": check:\n", s[0])

		ast := AST(strings.TrimSpace(s[0]))

		fmt.Println(RenderJSON(ast))

		fmt.Println(RenderAozoraText(ast))

		fmt.Println("#####################################")
	}
}

//go:embed testdata/s1.txt
var xdata string

func TestX(t *testing.T) {

	tokenizeAll(xdata, 0)

}
