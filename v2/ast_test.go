package aozoratext

import (
	_ "embed" //for embedding
	"fmt"
	"strings"
	"testing"
)

var renderer func(*node) string

func init() {

	setOutputOption("raw")
	setOutputOption("full")

	renderer = renderHtml
	// renderer = renderAozoraText
	//
	// setCompatible()
}

//go:embed testdata/aztestingHTML.txt
var aozoraTestdata string

func TestAozora(t *testing.T) {

	//	SetJIS0208(true)

	checkline := -7 //42

	checkpairs := strings.Split(aozoraTestdata, "\n&&&\n")

	for counter, p := range checkpairs {

		if strings.HasPrefix(p, "+++++") {
			break
		}

		s := strings.Split(p, "\n:::\n")

		if len(s) != 2 {
			continue
		}

		fmt.Println(counter+1, ": ")

		b := tokenizeAll(s[0], 0)

		//		fmt.Println(b.listAllTokens())

		//		return

		n := getAST(b)

		if renderer(n) != s[1] {

			fmt.Println(s[0], "\n", "want:\n", s[1], "\n", "got:\n", renderer(n))

		} else {
			fmt.Print(":", "ok\n\n")
		}

		if counter+1 == checkline {

			fmt.Println("*****************************")
			fmt.Println("details of line", counter+1)
			fmt.Println()
			fmt.Println(b.listAllTokens())
			fmt.Println()
			fmt.Println(listTypes(n))

			fmt.Println(renderer(n))

			//			return

			r1 := []rune(s[1])
			r2 := []rune(renderer(n))

			k := len(r1)
			if len(r1) != len(r2) {
				fmt.Println("different lengths!")
				if len(r1) > len(r2) {
					k = len(r1)
				} else {
					k = len(r2)
				}
			}

			//	return

			for i := range k {

				if i >= len(r2) {
					fmt.Println(i, ":", r1[i], "---")
					continue
				}

				if i >= len(r1) {
					fmt.Println(i, ":", "---", r2[i])
					continue
				}

				if r1[i] != r2[i] {
					fmt.Println(i, ":", r1[i], r2[i], string(r1[i]), string(r2[i]), "***")

				} else {
					fmt.Println(i, ":", r1[i], r2[i], string(r1[i]))
				}

			}

			return
		}

	}

	return

}
