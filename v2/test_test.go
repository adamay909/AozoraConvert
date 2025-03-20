package aozoratext

import (
	_ "embed" //for embedding
	"fmt"
	"testing"
	"time"
)

//go:embed o.txt
var cdata string

func TestComplexity(t *testing.T) {
	/*
		doc := newNode("document")

		para := newNode("paragraph")

		doc.addChild(para)

		prevNode := para

		for _ = range 4000 {

			n := newNode("emphasis")

			n.setAttr("raw", "傍点")

			n.setAttr("format type", "傍点")

			n.setAttr("raw closer", "傍点終わり")

			prevNode.addChild(n)

			prevNode = n

			//		n = newNode("text")

			//		n.setAttr("raw", "hello world")

			//		prevNode.addChild(n)

			//		prevNode = prevNode.parent()

		}

		n := newNode("text")

		n.setAttr("raw", "hello world")

		prevNode.addChild(n)

		cdata = renderAozoraText(doc)

		fmt.Println("done constructing test data")
	*/

	str := `　　姪の浜の大火
　　　　　名刹《めいさつ》如月寺《にょげつじ》に延焼［＃「名刹如月寺に延焼」は６段階大きな文字］
　　　　　　　　　　　　放火女無残の焼死を遂《と》ぐ［＃「放火女無残の焼死を遂ぐ」は３段階大きな文字］`

	//	str = `名刹如月寺に延焼［＃「名刹如月寺に延焼」は６段階大きな文字］`

	str = `名刹《めいさつ》如月寺《にょげつじ》に延焼［＃「如月寺に延焼」は６段階大きな文字］

放火女無残の焼死を遂《と》ぐ［＃「放火女無残の焼死を遂ぐ」は３段階大きな文字］`

	output.Reset()

	for _ = range 1000 {

		output.WriteString(str)

	}

	cdata = output.String()

	fmt.Println(len(cdata))

	startTime := time.Now()

	tok := tokenizeAll(cdata, 0)

	fmt.Println("time needed for tokenization", time.Since(startTime).Milliseconds())

	startTime = time.Now()

	tree := getAST(tok)

	fmt.Println("time needed for constructing ast", time.Since(startTime).Milliseconds())

	return

	fmt.Println(renderHtml(tree))

	return

}

func TestFurigana(t *testing.T) {

	s := "仝々〆〇ヶ〻"

	for _, c := range s {

		fmt.Println(string(c), charType(c))

	}
}
